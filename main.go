package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alvarowolfx/onion-weather-station/middleware"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/alvarowolfx/onion-weather-station/weather"
	"go.opencensus.io/stats/view"

	"periph.io/x/periph/host"
)

const (
	displayChangeInterval = 2000 * time.Millisecond
)

var (
	weatherStation weather.Station
)

func initOpenCensus() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	defaultMonitoringLabels := &stackdriver.Labels{}
	defaultMonitoringLabels.Set("hostname", hostname, "Sensor hostname")
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID:               "weather-station-iot-170004",
		DefaultMonitoringLabels: defaultMonitoringLabels,
	})

	if err != nil {
		log.Fatalf("Failed to register exporter: %v", err)
	}

	defer sd.Flush()

	// Register it as a metrics exporter
	view.RegisterExporter(sd)
	view.SetReportingPeriod(1 * time.Second)

	viewList := []*view.View{
		middleware.CurrentTemperatureView,
		middleware.CurrentPressureView,
		middleware.CurrentLedView,
		middleware.LedChangesView,
		middleware.MeasuresView,
	}

	if err := view.Register(viewList...); err != nil {
		log.Fatalf("Failed to register the views: %v", err)
	}
}

func main() {
	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}

	initOpenCensus()

	//wpsButton := gpioreg.ByName("")
	led := weather.NewPeriphLed("11")
	display, err := weather.NewTM1637Display("17", "16")
	if err != nil {
		log.Fatalf("failed to initialize tm1637: %v", err)
	}
	bmp, err := weather.NewBMP280EnvironmentSensor("/dev/i2c-0")
	if err != nil {
		log.Fatal(err)
	}
	weatherStation = weather.NewWeatherStation(display, led, bmp, nil)
	defer weatherStation.Close()

	go weatherStation.Start()

	port := "8080"
	httpWeatherStation := middleware.NewHttpWeatherStation(weatherStation, port)
	go httpWeatherStation.Start() // Configure http handlers

	weatherStationReporter := middleware.NewWeatherStationStatsReporter(weatherStation)
	go weatherStationReporter.Start() // Configure open census reporter

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigs // Wait for signal
		log.Println(sig)

		// Shutdown all services

		httpWeatherStation.Stop()
		log.Println("Http server stopped")

		weatherStationReporter.Stop()
		log.Println("Reporter stopped")

		weatherStation.Stop()
		log.Println("Weather Station stopped")

		done <- true

	}()

	log.Println("Press ctrl+c to stop...")
	<-done // Wait
}
