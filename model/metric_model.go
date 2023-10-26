package model

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	FileSize                    *prometheus.GaugeVec
	DownloadTime, DownloadSpeed *prometheus.GaugeVec
	SaveTime, SaveSpeed         *prometheus.GaugeVec
	DeleteTime, DeleteSpeed     *prometheus.GaugeVec
}
