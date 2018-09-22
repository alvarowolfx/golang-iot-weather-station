package middleware

import (
	"context"
	"fmt"
	"time"

	"go.opencensus.io/stats/view"

	"go.opencensus.io/stats"

	"github.com/alvarowolfx/onion-weather-station/weather"
)

var (
	MTemperature = stats.Float64("station/temperature", "Current Weather Station Temperature", stats.UnitDimensionless)
	MPressure    = stats.Float64("station/pressure", "Current Weather Station Pressure", stats.UnitDimensionless)

	MLedStatus  = stats.Int64("station/led/status", "Current Weather Station Led Status", stats.UnitDimensionless)
	MLedChanges = stats.Int64("station/led/changes", "Current Changes", stats.UnitDimensionless)

	MMeasure = stats.Int64("station/measure", "Current Weather Station Measurements", stats.UnitDimensionless)
)

var (
	CurrentTemperatureView = &view.View{
		Name:        "station/temperature/current",
		Measure:     MTemperature,
		Aggregation: view.LastValue(),
	}

	CurrentPressureView = &view.View{
		Name:        "station/pressure/current",
		Measure:     MPressure,
		Aggregation: view.LastValue(),
	}

	LedChangesView = &view.View{
		Name:        "station/led/changes",
		Measure:     MLedChanges,
		Aggregation: view.Count(),
	}

	MeasuresView = &view.View{
		Name:        "station/measure",
		Measure:     MMeasure,
		Aggregation: view.Count(),
	}

	CurrentLedView = &view.View{
		Name:        "station/led/current",
		Measure:     MLedStatus,
		Aggregation: view.LastValue(),
	}
)

type WeatherStationStatsReporter struct {
	station weather.Station
	ticker  *time.Ticker
}

func NewWeatherStationStatsReporter(station weather.Station) *WeatherStationStatsReporter {
	return &WeatherStationStatsReporter{
		station: station,
	}
}

func (sr *WeatherStationStatsReporter) Start() {
	if sr.ticker != nil {
		sr.ticker.Stop()
	}
	sr.ticker = time.NewTicker(1 * time.Second)
	ctx := context.Background()

	lastLedState := !sr.station.GetLedState()
	for tick := false; ; tick = !tick {

		t, err := sr.station.ReadTemperature()
		p, err := sr.station.ReadPressure()
		currentLedState := sr.station.GetLedState()
		var ledStatus int64
		ledStatus = 0
		if currentLedState {
			ledStatus = 1
		}

		if err != nil {
			fmt.Println("Error reading sensor %v", err)
		}

		stats.Record(ctx, MTemperature.M(t))
		stats.Record(ctx, MPressure.M(p))
		stats.Record(ctx, MLedStatus.M(ledStatus))

		if currentLedState != lastLedState {
			stats.Record(ctx, MLedChanges.M(1))
			lastLedState = currentLedState
		}

		stats.Record(ctx, MMeasure.M(1))

		<-sr.ticker.C
	}
}

func (sr WeatherStationStatsReporter) Stop() {
	sr.ticker.Stop()
}
