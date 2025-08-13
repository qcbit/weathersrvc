package public

type weatherResponse struct {
	ShortForecast   string  `json:"shortForecast"`
	Temperature     float64 `json:"temperature"`
	TemperatureUnit string  `json:"temperatureUnit"`
	TemperatureFeeling string `json:"temperatureFeeling"`
}

type forecastData struct {
	Properties struct {
		Periods []struct {
			weatherResponse
		} `json:"periods"`
	} `json:"properties"`
}
