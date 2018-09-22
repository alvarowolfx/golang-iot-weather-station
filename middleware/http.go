package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alvarowolfx/onion-weather-station/weather"
)

type HttpWeatherStation struct {
	station weather.Station
	server  *http.Server
	port    string
}

type HttpResponse struct {
	Message string            `json:"message"`
	State   map[string]string `json:"state"`
}

func NewHttpWeatherStation(station weather.Station, port string) *HttpWeatherStation {
	return &HttpWeatherStation{
		station: station,
		port:    port,
	}
}

func (hw *HttpWeatherStation) Start() {
	weatherURL := fmt.Sprintf("/station/weather")
	ledURL := fmt.Sprintf("/station/led/toggle")

	http.HandleFunc(weatherURL, func(res http.ResponseWriter, req *http.Request) {
		temperature, _ := hw.station.ReadTemperature()
		pressure, _ := hw.station.ReadPressure()

		json.NewEncoder(res).Encode(HttpResponse{
			Message: "ok",
			State: map[string]string{
				"temperature": fmt.Sprintf("%2.2f", temperature),
				"pressure":    fmt.Sprintf("%2.2f", pressure),
			},
		})
	})

	http.HandleFunc(ledURL, func(res http.ResponseWriter, req *http.Request) {
		ledStatus := hw.station.GetLedState()
		hw.station.ToggleLED()

		json.NewEncoder(res).Encode(HttpResponse{
			Message: "ok",
			State: map[string]string{
				"led": fmt.Sprintf("%s", ledStatus),
			},
		})
	})

	hw.server = &http.Server{Addr: ":" + hw.port}
	go hw.server.ListenAndServe() // Http server blocks execution
}

func (hw HttpWeatherStation) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	hw.server.Shutdown(ctx)
}
