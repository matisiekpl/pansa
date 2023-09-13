package repository

import (
	"io"
	"net/http"
	"strings"
)

type WeatherRepository interface {
	GetMETAR(icao string) string
	GetTAF(icao string) string
}

type weatherRepository struct {
}

func newWeatherRepository() WeatherRepository {
	return &weatherRepository{}
}

func (w weatherRepository) GetMETAR(icao string) string {
	resp, err := http.Get("https://beta.aviationweather.gov/cgi-bin/data/metar.php?ids=" + icao)
	if err != nil {
		return ""
	}
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(result))
}

func (w weatherRepository) GetTAF(icao string) string {
	resp, err := http.Get("https://beta.aviationweather.gov/cgi-bin/data/taf.php?ids=" + icao)
	if err != nil {
		return ""
	}
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(result))
}
