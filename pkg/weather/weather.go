package weather

import (
	"time"

	"github.com/ottolote/naaprs/pkg/netatmo"
)

type WeatherReport struct {
	Timestamp       time.Time
	Lat 		float64
	Lon 		float64
	Altimeter       float64
	Humidity        int
	RainLastHour    float64
	RainLast24Hours float64
	RainToday       float64
	SolarRad        int
	Temp            int
	WindDir         int
	WindGust        int
	WindSpeed       int
}

func containsString(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func filterModules(modules []netatmo.ModuleData, dataType string) []netatmo.ModuleData {
	var result []netatmo.ModuleData
	for _, module := range modules {
		if containsString(module.DataType, dataType) {
			result = append(result, module)
		}
	}

	return result
}

func GetWeatherData(source string) *WeatherReport {
	netatmoModules := netatmo.GetAllModules()

	rainModules := filterModules(netatmoModules, "Rain")
	humidityModules := filterModules(netatmoModules, "Humidity")
	temperatureModules := filterModules(netatmoModules, "Temperature")
	windModules := filterModules(netatmoModules, "Wind")

	// TODO: allow configuration of which module to use, for now just pick first available
	rain := rainModules[0]
	humidity := humidityModules[0]
	temperature := temperatureModules[0]
	wind := windModules[0]

	return &WeatherReport{
		Timestamp: time.Now(),

		Lat: wind.Lat,
		Lon: wind.Lon,

		Humidity:        humidity.Humidity,
		RainLast24Hours: rain.RainLast24Hours,
		RainLastHour:    rain.RainLastHour,
		Temp:            int(temperature.Temp),

		WindDir:   wind.WindDir,
		WindGust:  wind.WindGust,
		WindSpeed: wind.WindSpeed,
	}
}
