package handler

import (
	"jougan/helper/monitor/monitor_impl"
)

func DiskHandler(pConfig *monitor_impl.PrometheusConfig) {
	pConfig.Test()
}
