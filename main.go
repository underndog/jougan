package main

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "jougan/docs"
	"jougan/handler"
	"jougan/helper"
	"jougan/helper/aws_cloud/aws_cloud_impl"
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

// @title Jougan API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {

	// It asserts that the interface returned by NewMonitoring actually holds a pointer to a PrometheusConfig struct
	pConfig := monitor_impl.NewMonitoring(nil).(*monitor_impl.PrometheusConfig)
	// reg is a variable that holds a reference to a Prometheus registry object
	reg := pConfig.Registry
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	awsCloud := &aws_cloud_impl.AWSConfiguration{
		Region: helper.GetEnvOrDefault("AWS_REGION", "us-west-2"),
	}

	inspectDiskHandler := handler.InspectDiskHandler{
		Monitoring: pConfig,
		AWSCloud:   awsCloud,
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

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := router.API{
		Echo:               e,
		PromHandler:        promHandler,
		InspectDiskHandler: inspectDiskHandler,
	}
	api.SetupRouter()
	e.Logger.Fatal(e.Start(":1994"))
}
