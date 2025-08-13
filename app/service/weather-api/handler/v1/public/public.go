// Package public maintains the handlers for public access.
package public

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/qcbit/weathersrvc/foundation/web"
)

type Handler struct {
	Log *zap.SugaredLogger
}

// Coordinates handles requests for weather data based on coordinates "latitude,longitude" format
// and returns the weather response.
func (h *Handler) Coordinates(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	coord := web.Param(r, "coordinates")

	switch coord {
	case "":
		return web.Response(ctx, w, "", http.StatusOK)
	default:
		h.Log.Infow("fetching weather data for coordinates", "coordinates", coord)
		weather, err := h.getWeatherData(coord)
		if err != nil {
			h.Log.Errorw("failed to fetch weather data", "coordinates", coord, "error", err)
			return web.Response(ctx, w, fmt.Sprintf("failed to fetch weather data: %v", err), http.StatusInternalServerError)
		}
		weather.TemperatureFeeling = getWeatherFeeling(weather.Temperature, weather.TemperatureUnit)
		h.Log.Infow("weather data response ready", "coordinates", coord, "shortForecast", weather.ShortForecast, "temperature", weather.Temperature, "temperatureUnit", weather.TemperatureUnit, "temperatureFeeling", weather.TemperatureFeeling)

		return web.Response(ctx, w, weather, http.StatusOK)
	}
}

// getWeatherFeeling determines the feeling of the temperature.
func getWeatherFeeling(temperature float64, temperatureUnit string) string {
	temperatureUnit = strings.ToLower(temperatureUnit)
	if temperatureUnit != "f" && temperatureUnit != "c" {
		return "-"
	}

	// Convert Celsius to Fahrenheit
	fTemperature := temperature
	if temperatureUnit == "c" {
		fTemperature = (temperature * 9 / 5) + 32
	}

	if fTemperature >= 85 {
		return "hot"
	} else if fTemperature <= 60 {
		return "cold"
	}
	return "moderate"
}

type PointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

// getWeatherData fetches weather data from the external API based on coordinates.
func (h *Handler) getWeatherData(coord string) (weatherResponse, error) {
	resp, err := http.Get("https://api.weather.gov/points/" + coord)
	if err != nil {
		return weatherResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.Log.Errorw("failed to fetch weather data", "status", resp.Status)
		return weatherResponse{}, fmt.Errorf("failed to fetch weather data: %s", resp.Status)
	}

	var points PointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&points); err != nil {
		return weatherResponse{}, err
	}

	// Get the forecast for the URL from the points response
	resp, err = http.Get(points.Properties.Forecast)
	if err != nil {
		return weatherResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.Log.Errorw("failed to fetch forecast data", "status", resp.Status)
		return weatherResponse{}, fmt.Errorf("failed to fetch forecast data: %s", resp.Status)
	}

	// Unmarshal the forecast data
	var forecast forecastData
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return weatherResponse{}, err
	}

	if len(forecast.Properties.Periods) == 0 {
		return weatherResponse{}, fmt.Errorf("no forecast data available")
	}

	// Get the forecast for today
	weather := forecast.Properties.Periods[0] 
	h.Log.Infow("weather data fetched", "shortForecast", weather.ShortForecast, "temperature", weather.Temperature, "temperatureUnit", weather.TemperatureUnit)

	return weather.weatherResponse, nil
}
