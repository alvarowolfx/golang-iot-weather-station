package weather

import (
	"fmt"
	"io"
	"time"
)

type Display interface {
	io.Closer
	Display(value float64) error
}

type Led interface {
	io.Closer
	On()
	Off()
	GetCurrentState() bool
}

type EnvironmentSensor interface {
	io.Closer
	ReadTemperature() (float64, error)
	ReadPressure() (float64, error)
}

type Station interface {
	io.Closer
	ReadTemperature() (float64, error)
	ReadPressure() (float64, error)
	ToggleLED()
	Display(value float64) error
	Start()
	Stop()
}

type WeatherStation struct {
	display               Display
	led                   Led
	environmentSensor     EnvironmentSensor
	displayChangeInterval time.Duration
	ticker                *time.Ticker
}

type WeatherStationOpts struct {
	displayChangeInterval time.Duration
}

var DefaultWeatherStationOpts = WeatherStationOpts{
	displayChangeInterval: 2000 * time.Millisecond,
}

func NewWeatherStation(display Display, led Led, environmentSensor EnvironmentSensor, opts *WeatherStationOpts) Station {
	if opts == nil {
		opts = &DefaultWeatherStationOpts
	}
	return &WeatherStation{
		display:               display,
		led:                   led,
		environmentSensor:     environmentSensor,
		displayChangeInterval: opts.displayChangeInterval,
	}
}

func (ws *WeatherStation) Start() {
	var err error
	if ws.ticker != nil {
		ws.ticker.Stop()
	}
	ws.ticker = time.NewTicker(ws.displayChangeInterval)
	for tick := false; ; tick = !tick {
		metric := 0.0
		if tick {
			metric, err = ws.ReadTemperature()
		} else {
			metric, err = ws.ReadPressure()
		}

		if err != nil {
			fmt.Println("Error reading sensor %v", err)
		}

		err = ws.Display(metric)

		if err != nil {
			fmt.Println("Error displaying value %v", err)
		}

		ws.ToggleLED()
		<-ws.ticker.C
	}
}

func (ws WeatherStation) Stop() {
	ws.ticker.Stop()
}

func (ws *WeatherStation) ReadTemperature() (float64, error) {
	return ws.environmentSensor.ReadTemperature()
}

func (ws *WeatherStation) ReadPressure() (float64, error) {
	return ws.environmentSensor.ReadPressure()
}

func (ws *WeatherStation) ToggleLED() {
	state := ws.led.GetCurrentState()
	if state {
		ws.led.Off()
	} else {
		ws.led.On()
	}
}

func (ws *WeatherStation) Display(value float64) error {
	if value < 0 {
		return nil
	}
	return ws.display.Display(value)
}

func (ws *WeatherStation) Close() error {
	err := ws.display.Close()
	err = ws.environmentSensor.Close()
	err = ws.led.Close()
	return err
}
