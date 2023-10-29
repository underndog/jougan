package model

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	FileSize                    *prometheus.GaugeVec
	DownloadTime, DownloadSpeed *prometheus.GaugeVec
	SaveTime, SaveSpeed         *prometheus.GaugeVec
	DeleteTime, DeleteSpeed     *prometheus.GaugeVec
}

type MetricResponse struct {
	FileSize      int     `json:"fileSize,omitempty"`
	DownloadTime  float64 `json:"downloadTime,omitempty"`
	DownloadSpeed float64 `json:"downloadSpeed,omitempty"`
	SaveTime      float64 `json:"saveTime,omitempty"`
	SaveSpeed     float64 `json:"saveSpeed,omitempty"`
	DeleteTime    float64 `json:"deleteTimestring,omitempty"`
	DeleteSpeed   float64 `json:"deleteSpeed,omitempty"`
}
