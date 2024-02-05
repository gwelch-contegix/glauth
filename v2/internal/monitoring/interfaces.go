package monitoring

import "github.com/gwelch-contegix/ldaps"

type MonitorInterface interface {
	SetResponseTimeMetric(map[string]string, float64) error
	SetLDAPMetric(map[string]string, float64) error
}

type LDAPServerInterface interface {
	SetStats(bool)
	GetStats() ldaps.Stats
}
