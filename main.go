package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alvarowolfx/onion-weather-station/middleware"
	"github.com/alvarowolfx/onion-weather-station/weather"

	"periph.io/x/periph/host"
)

const (
	displayChangeInterval = 2000 * time.Millisecond
)

var (
	weatherStation weather.Station
)

func main() {
	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}

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

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigs // Wait for signal
		log.Println(sig)

		// Shutdown all services

		httpWeatherStation.Stop()
		log.Println("Http server stopped")

		weatherStation.Stop()
		log.Println("Weather Station stopped")

		done <- true

	}()

	log.Println("Press ctrl+c to stop...")
	<-done // Wait
}
