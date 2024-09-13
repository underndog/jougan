package monitor_impl

import (
	"github.com/prometheus/client_golang/prometheus"
	"jougan/helper/monitor"
	"jougan/model"
)

type PrometheusConfig struct {
	Registry *prometheus.Registry
	Metrics  *model.Metrics
}

/*
**
Create a helper function to initialize a GaugeVec metric
Avoid "panic: duplicate metrics collector registration attempted" when calling again "reg.MustRegister(metric)"
with Each Metrics, Only one to call reg.MustRegister()
*/
func initMetric(reg *prometheus.Registry, namespace, name, help string, labels []string) *prometheus.GaugeVec {
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      help,
		},
		labels,
	)
	reg.MustRegister(metric)
	return metric
}

func NewMonitoring(pC *PrometheusConfig) monitor.Monitoring {
	reg := prometheus.NewRegistry()

	//Use the helper function to initialize each metric
	m := &model.Metrics{
		FileSize:      initMetric(reg, "jougan", "file_size", "Size of the file in Bytes.", []string{"filename"}),
		DownloadSpeed: initMetric(reg, "jougan", "download_speed", "Download speed (unit: B/s).", []string{"filename"}),
		DownloadTime:  initMetric(reg, "jougan", "download_time_seconds", "Time required to download the file (unit: seconds).", []string{"filename"}),
		SaveSpeed:     initMetric(reg, "jougan", "save_speed", "Save speed (unit: B/s).", []string{"filename"}),
		SaveTime:      initMetric(reg, "jougan", "save_time_seconds", "Time required to save the file (unit: seconds).", []string{"filename"}),
		DeleteSpeed:   initMetric(reg, "jougan", "delete_speed", "Delete speed (unit: B/s).", []string{"filename"}),
		DeleteTime:    initMetric(reg, "jougan", "delete_time_seconds", "Time required to delete the file (unit: seconds).", []string{"filename"}),
		UploadSpeed:   initMetric(reg, "jougan", "upload_speed", "Upload speed (unit: B/s).", []string{"filename"}),
		UploadTime:    initMetric(reg, "jougan", "upload_time_seconds", "Time required to upload the file (unit: seconds).", []string{"filename"}),
	}
	return &PrometheusConfig{
		Registry: reg,
		Metrics:  m,
	}
}
