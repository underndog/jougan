package monitor_impl

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	devices prometheus.Gauge
}
type Device struct {
	ID       int    `json:"id"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
}

func (pC *PrometheusConfig) Test() {

	reg := pC.Registry
	m := &metrics{
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "connected_devices",
			Help:      "Number of currently connected devices.",
		}),
	}
	reg.MustRegister(m.devices)
	dvs := []Device{
		{1, "5F-33-CC-1F-43-82", "2.1.6"},
		{2, "EF-2B-C4-F5-D6-34", "2.1.6"},
	}
	m.devices.Set(float64(len(dvs)))
}
