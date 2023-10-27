package main

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"jougan/handler"
	"jougan/helper/monitor/monitor_impl"
	"jougan/log"
	"jougan/router"
	"os"
	"time"
)

func init() {
	os.Setenv("APP_NAME", "jougan-inspects-disk")
	log.InitLogger(false)
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
}

func main() {

	// It asserts that the interface returned by NewMonitoring actually holds a pointer to a PrometheusConfig struct
	pConfig := monitor_impl.NewMonitoring(nil).(*monitor_impl.PrometheusConfig)
	// reg is a variable that holds a reference to a Prometheus registry object
	reg := pConfig.Registry
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	inspectDiskHandler := handler.InspectDiskHandler{
		Monitoring: pConfig,
	}
	// Run the DiskHandler function continuously in a loop
	go func() {
		for {
			startTime := time.Now()

			inspectDiskHandler.DiskHandler()

			elapsedTime := time.Since(startTime)
			sleepTime := 15*time.Second - elapsedTime

			// Sleep only if the elapsedTime is less than 15 seconds
			if sleepTime > 0 {
				time.Sleep(sleepTime)
			}
		}
	}()

	e := echo.New()
	api := router.API{
		Echo:        e,
		PromHandler: promHandler,
	}
	api.SetupRouter()
	e.Logger.Fatal(e.Start(":1994"))
}
