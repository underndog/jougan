package monitor_impl

import (
	"github.com/prometheus/client_golang/prometheus"
)

func (pC *PrometheusConfig) FileSizeMonitor(filename string, fileSize float64) {
	pC.Metrics.FileSize.With(prometheus.Labels{"filename": filename}).Set(fileSize)
}

func (pC *PrometheusConfig) SpeedMonitor(filename string, activity string, speed float64, time float64) {
	switch activity {
	case "download":
		pC.Metrics.DownloadSpeed.With(prometheus.Labels{"filename": filename}).Set(speed)
		pC.Metrics.DownloadTime.With(prometheus.Labels{"filename": filename}).Set(time)
	case "save":
		pC.Metrics.SaveSpeed.With(prometheus.Labels{"filename": filename}).Set(speed)
		pC.Metrics.SaveTime.With(prometheus.Labels{"filename": filename}).Set(time)
	case "delete":
		pC.Metrics.DeleteSpeed.With(prometheus.Labels{"filename": filename}).Set(speed)
		pC.Metrics.DeleteTime.With(prometheus.Labels{"filename": filename}).Set(time)
	default:
		// Handle an unknown activity type if needed.
		return
	}
}
