module github.com/gwelch-contegix/glauth/v2

go 1.21

require (
	github.com/BurntSushi/toml v1.4.0
	github.com/GeertJohan/yubigo v0.0.0-20190917122436-175bc097e60e
	github.com/arl/statsviz v0.6.0
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/fsnotify/fsnotify v1.7.0
	github.com/go-ldap/ldap/v3 v3.4.8
	github.com/gwelch-contegix/ldaps v0.0.0-20240724011018-d6deb1b4df72
	github.com/jinzhu/copier v0.4.0
	github.com/pires/go-proxyproto v0.7.0
	github.com/pquerna/otp v1.4.0
	github.com/prometheus/client_golang v1.19.1
	github.com/prometheus/client_model v0.6.1
	github.com/rs/zerolog v1.33.0
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.2.3
	github.com/yaegashi/msgraph.go v0.1.4
	go.opentelemetry.io/contrib/propagators/jaeger v1.28.0
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.28.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/trace v1.28.0
	go.uber.org/mock v0.4.0
	golang.org/x/crypto v0.25.0
	gopkg.in/amz.v3 v3.0.0-20201001071545-24fc1eceb27b
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.0.2 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.7 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rickb777/date v1.20.7 // indirect
	github.com/rickb777/plural v1.4.2 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240723171418-e6d459c13d2a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240723171418-e6d459c13d2a // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/smartystreets/assertions => github.com/smarty/assertions v1.15.1
