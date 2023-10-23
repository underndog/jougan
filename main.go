package main

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"jougan/handler"
	"jougan/helper/monitor/monitor_impl"
	"jougan/log"
	"jougan/router"
	"os"
)

func init() {
	os.Setenv("APP_NAME", "jougan-inspects-disk")
	log.InitLogger(false)
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
}

func main() {

	// It asserts that the interface returned by NewMonitoring actually holds a pointer to a PrometheusConfig struct
	pConfig := monitor_impl.NewMonitoring(nil).(*monitor_impl.PrometheusConfig)
	reg := pConfig.Registry
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	handler.DiskHandler(pConfig)

	e := echo.New()
	api := router.API{
		Echo:        e,
		PromHandler: promHandler,
	}
	api.SetupRouter()
	e.Logger.Fatal(e.Start(":1994"))
}
