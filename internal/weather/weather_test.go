package weather

import (
	"errors"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFetchWeatherData(t *testing.T) {
	type constraint struct {
		noApiKey       bool
		failtoFetching bool
	}

	tests := []struct {
		name       string
		constraint *constraint
		city       string
		err        error
	}{
		{
			name: "not set apiKey",
			constraint: &constraint{
				noApiKey: true,
			},
			city: "",
			err:  newApiKeyNotFoundError(),
		},
		{
			name: "empty city name",
			city: "",
			err:  newCityParameterNotFoundError(),
		},
		{
			name: "fail to fetch weather data",
			constraint: &constraint{
				failtoFetching: true,
			},
			city: "Tokyo",
			err:  newFetchingWeatherDataFailedError(&url.Error{}),
		},
		{
			name: "success fetch weather data",
			city: "Tokyo",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// var apiKey string
			if tt.constraint != nil {
				if tt.constraint.noApiKey {
					apiKey := os.Getenv("OPENWEATHER_API_KEY")
					os.Setenv("OPENWEATHER_API_KEY", "")
					defer func() {
						os.Setenv("OPENWEATHER_API_KEY", apiKey)
					}()
				}
				if tt.constraint.failtoFetching {
					var tmp string
					apiKey := os.Getenv("OPENWEATHER_API_KEY")
					os.Setenv("OPENWEATHER_API_KEY", "testapi")
					apiURL, tmp = "https://gotest.maesh.dev?q=%s&appid=%s", apiURL
					defer func() {
						apiURL = tmp
						os.Setenv("OPENWEATHER_API_KEY", apiKey)
					}()
				}
			}

			_, err := FetchWeatherData(tt.city)
			if diff := cmp.Diff(tt.err, err, cmp.Comparer(func(x, y error) bool {
				if fwdfe, ok := err.(*FetchingWeatherDataFailedError); ok {
					if _, ok := fwdfe.err.(*url.Error); ok {
						return true
					}
				}
				return errors.Is(err, tt.err)
			})); diff != "" {
				t.Errorf("Test %q failed (-want +got):\n%s", tt.name, diff)
			}
		})
	}
}
