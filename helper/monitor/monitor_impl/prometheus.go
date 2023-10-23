package monitor_impl

import (
	"github.com/prometheus/client_golang/prometheus"
	"jougan/helper/monitor"
)

type PrometheusConfig struct {
	Registry *prometheus.Registry
}

func NewMonitoring(pC *PrometheusConfig) monitor.Monitoring {
	reg := prometheus.NewRegistry()
	pC = &PrometheusConfig{Registry: reg}
	return pC
}
