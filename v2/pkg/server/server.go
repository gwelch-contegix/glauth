package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"plugin"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"

	"github.com/GeertJohan/yubigo"
	"github.com/gwelch-contegix/glauth/v2/internal/monitoring"
	"github.com/gwelch-contegix/glauth/v2/pkg/config"
	"github.com/gwelch-contegix/glauth/v2/pkg/handler"
	"github.com/gwelch-contegix/ldaps"
	proxyproto "github.com/pires/go-proxyproto"
)

type LdapSvc struct {
	c        *config.Config
	yubiAuth *yubigo.YubiAuth
	l        *ldaps.Server

	monitor monitoring.MonitorInterface
	tracer  trace.Tracer
	log     zerolog.Logger
}

func NewServer(opts ...Option) (*LdapSvc, error) {
	options := newOptions(opts...)

	s := LdapSvc{
		log:     options.Logger,
		c:       options.Config,
		monitor: options.Monitor,
		tracer:  options.Tracer,
	}

	var err error

	if len(s.c.YubikeyClientID) > 0 && len(s.c.YubikeySecret) > 0 {
		s.yubiAuth, err = yubigo.NewYubiAuth(s.c.YubikeyClientID, s.c.YubikeySecret)

		if err != nil {
			return nil, errors.New("Yubikey Auth failed")
		}
	}

	var helper handler.Handler

	loh := handler.NewLDAPOpsHelper(s.tracer)

	// instantiate the helper, if any
	if s.c.Helper.Enabled {
		switch s.c.Helper.Datastore {
		case "config":
			helper = handler.NewConfigHandler(
				handler.Logger(&s.log),
				handler.Config(s.c),
				handler.YubiAuth(s.yubiAuth),
				handler.LDAPHelper(loh),
				handler.Tracer(s.tracer),
			)
		case "plugin":
			plug, err := plugin.Open(s.c.Helper.Plugin)
			if err != nil {
				return nil, fmt.Errorf("Unable to load specified helper plugin: %s", err)
			}
			nph, err := plug.Lookup(s.c.Helper.PluginHandler)
			if err != nil {
				return nil, errors.New("unable to find 'NewPluginHandler' in loaded helper plugin")
			}
			initFunc, ok := nph.(func(...handler.Option) handler.Handler)

			if !ok {
				return nil, errors.New("loaded helper plugin lacks a proper NewPluginHandler function")
			}
			// Normally, here, we would somehow have imported our plugin into our
			// handler namespace. Oops?
			helper = initFunc(
				handler.Logger(&s.log),
				handler.Config(s.c),
				handler.YubiAuth(s.yubiAuth),
				handler.LDAPHelper(loh),
				handler.Tracer(s.tracer),
			)
		default:
			return nil, fmt.Errorf("unsupported helper %s - must be one of 'config', 'plugin'", s.c.Helper.Datastore)
		}
		s.log.Info().Str("datastore", s.c.Helper.Datastore).Msg("Using helper")
	}

	backendCounter := -1
	allHandlers := handler.HandlerWrapper{Handlers: make([]handler.Handler, 10), Count: &backendCounter}

	// configure the backends
	s.l = ldaps.NewServer()
	s.l.EnforceLDAP = true

	if tlsConfig := options.TLSConfig; tlsConfig != nil {
		s.l.TLSConfig = tlsConfig
		s.log.Info().Interface("tls.certificates", tlsConfig.Certificates).Msg("enabling LDAP over TLS")
	}

	for i, backend := range s.c.Backends {
		var h handler.Handler
		switch backend.Datastore {
		case "ldap":
			h = handler.NewLdapHandler(
				handler.Backend(backend),
				handler.Handlers(allHandlers),
				handler.Logger(&s.log),
				handler.Helper(helper),
				handler.Monitor(s.monitor),
				handler.Tracer(s.tracer),
			)
		case "owncloud":
			h = handler.NewOwnCloudHandler(
				handler.Backend(backend),
				handler.Logger(&s.log),
				handler.Monitor(s.monitor),
				handler.Tracer(s.tracer),
			)
		case "config":
			h = handler.NewConfigHandler(
				handler.Backend(backend),
				handler.Logger(&s.log),
				handler.Config(s.c), // TODO only used to access Users and Groups, move that to dedicated options
				handler.YubiAuth(s.yubiAuth),
				handler.LDAPHelper(loh),
				handler.Monitor(s.monitor),
				handler.Tracer(s.tracer),
			)
		case "plugin":
			plug, err := plugin.Open(backend.Plugin)
			if err != nil {
				return nil, fmt.Errorf("Unable to load specified backend plugin: %s", err)
			}
			nph, err := plug.Lookup(backend.PluginHandler)
			if err != nil {
				return nil, errors.New("unable to find 'NewPluginHandler' in loaded backend plugin")
			}
			initFunc, ok := nph.(func(<-chan bool, ...handler.Option) handler.Handler)

			if !ok {
				return nil, errors.New("loaded backend plugin lacks a proper NewPluginHandler function")
			}
			// Normally, here, we would somehow have imported our plugin into our
			// handler namespace. Oops?
			h = initFunc(
				s.l.Quit,
				handler.Backend(backend),
				handler.Logger(&s.log),
				handler.Config(s.c),
				handler.YubiAuth(s.yubiAuth),
				handler.LDAPHelper(loh),
				handler.Monitor(s.monitor),
				handler.Tracer(s.tracer),
			)
		default:
			return nil, fmt.Errorf("unsupported backend %s - must be one of 'config', 'ldap','owncloud' or 'plugin'", backend.Datastore)
		}
		s.log.Info().Str("datastore", backend.Datastore).Int("position", i).Msg("Loading backend")

		// Only our first backend will answer proper LDAP queries.
		// Note that this could evolve towars something nicer where we would maintain
		// multiple binders in addition to the existing multiple LDAP backends
		if i == 0 {
			s.l.BindFunc("", h)
			s.l.SearchFunc("", h)
			s.l.CloseFunc("", h)
		}
		allHandlers.Handlers[i] = h
		backendCounter++
	}

	monitoring.NewLDAPMonitorWatcher(s.l, s.monitor, &s.log)

	return &s, nil
}

// ListenAndServe listens on the TCP network address s.c.LDAP.Listen
func (s *LdapSvc) ListenAndServe() error {
	s.log.Info().Str("address", s.c.LDAP.Listen).Msg("LDAP server listening")

	ln, err := net.Listen("tcp", s.c.LDAP.Listen)
	if err != nil {
		return err
	}
	ln = &proxyproto.Listener{
		Listener:          ln,
		ReadHeaderTimeout: 2 * time.Second,
	}
	defer ln.Close()
	return s.l.Serve(ln)
}

// ListenAndServeTLS listens on the TCP network address s.c.LDAPS.Listen
func (s *LdapSvc) ListenAndServeTLS() error {
	s.log.Info().Str("address", s.c.LDAPS.Listen).Msg("LDAPS server listening")

	cert, err := tls.LoadX509KeyPair(s.c.LDAPS.Cert, s.c.LDAPS.Key)
	if err != nil {
		return err
	}
	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}
	tlsConfig.ServerName = "localhost"
	ln, err := tls.Listen("tcp", s.c.LDAPS.Listen, &tlsConfig)
	if err != nil {
		return err
	}
	ln = &proxyproto.Listener{
		Listener:          ln,
		ReadHeaderTimeout: 2 * time.Second,
	}
	defer ln.Close()
	return s.l.Serve(ln)
}

// Shutdown ends listeners by closing the ldap server quit channel
func (s *LdapSvc) Shutdown() {
	close(s.l.Quit)
}
