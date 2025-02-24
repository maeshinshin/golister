package weatherclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var apiURL string = "https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric"

func FetchWeatherData(city string) (*WeatherData, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")

	if apiKey == "" {
		return nil, newApiKeyNotFoundError()
	}

	if city == "" {
		return nil, newCityParameterNotFoundError()
	}

	url := fmt.Sprintf(apiURL, city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, newFetchingWeatherDataFailedError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, newApiRequestFailedError(resp.StatusCode)
	}

	var weatherData WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
}
