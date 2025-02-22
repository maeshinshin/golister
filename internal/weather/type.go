package weather

import "fmt"

type Weather struct {
	Description string `json:"description"`
}

type Main struct {
	Temp       float64 `json:"temp"`
	Feels_like float64 `json:"feels_like"`
	Temp_min   float64 `json:"temp_min"`
	Temp_max   float64 `json:"temp_max"`
	Pressure   float64 `json:"pressure"`
	Humidity   float64 `json:"humidity"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

type Clouds struct {
	All int `json:"all"`
}

type WeatherData struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
	Wind    Wind      `json:"wind"`
	Clouds  Clouds    `json:"clouds"`
	Name    string    `json:"name"`
}

type CityParameterNotFoundError struct{}

func (e *CityParameterNotFoundError) Error() string {
	return weatherError("City parameter is not found.")
}

func newCityParameterNotFoundError() *CityParameterNotFoundError {
	return &CityParameterNotFoundError{}
}

type ApiKeyNotFoundError struct{}

func (e *ApiKeyNotFoundError) Error() string {
	return weatherError("API key is not set.")
}

func newApiKeyNotFoundError() *ApiKeyNotFoundError {
	return &ApiKeyNotFoundError{}
}

type FetchingWeatherDataFailedError struct {
	err error
}

func (e *FetchingWeatherDataFailedError) Error() string {
	return weatherError(fmt.Sprintf("Failed to fetch weather data.(%s)", e.err))
}

func newFetchingWeatherDataFailedError(err error) *FetchingWeatherDataFailedError {
	return &FetchingWeatherDataFailedError{err}
}

type ApiRequestFailedError struct {
	StatusCode int
}

func (e *ApiRequestFailedError) Error() string {
	return weatherError(fmt.Sprintf("External API returned an error. (Status code: %d)", e.StatusCode))
}

func newApiRequestFailedError(StatusCode int) *ApiRequestFailedError {
	return &ApiRequestFailedError{StatusCode: StatusCode}
}

func weatherError(message string) string {
	return fmt.Sprintf("Error /weather: %s", message)
}
