package weather

import (
	"time"

	"github.com/ottolote/naaprs/pkg/netatmo"
	"github.com/spf13/viper"
)

type WeatherReport struct {
	Timestamp       time.Time
	Lat             float64
	Lon             float64
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

func filterModulesByDataType(modules []netatmo.ModuleData, dataType string) []netatmo.ModuleData {
	var result []netatmo.ModuleData
	for _, module := range modules {
		if containsString(module.DataType, dataType) {
			result = append(result, module)
		}
	}

	return result
}

func filterModulesByName(modules []netatmo.ModuleData, name string) []netatmo.ModuleData {
	var result []netatmo.ModuleData
	for _, module := range modules {
		if module.Name == name {
			result = append(result, module)
		}
	}

	return result
}

func GetWeatherData(source string) *WeatherReport {
	netatmoModules := netatmo.GetAllModules()

	rainModules := filterModulesByDataType(netatmoModules, "Rain")
	humidityModules := filterModulesByDataType(netatmoModules, "Humidity")
	temperatureModules := filterModulesByDataType(netatmoModules, "Temperature")
	windModules := filterModulesByDataType(netatmoModules, "Wind")


	sourceRain := viper.GetString("SOURCE_RAIN")
	var rain netatmo.ModuleData
	if sourceRain == "" {
		rain = rainModules[0]
	} else {
		rain = filterModulesByName(netatmoModules, sourceRain)[0]
	}

	sourceWind := viper.GetString("SOURCE_WIND")
	var wind netatmo.ModuleData
	if sourceWind == "" {
		wind = windModules[0]
	} else {
		wind = filterModulesByName(netatmoModules, sourceWind)[0]
	}

	sourceTemperature := viper.GetString("SOURCE_TEMPERATURE")
	var temperature netatmo.ModuleData
	if sourceTemperature == "" {
		temperature = temperatureModules[0]
	} else {
		temperature = filterModulesByName(netatmoModules, sourceTemperature)[0]
	}

	sourceHumidity := viper.GetString("SOURCE_HUMIDITY")
	var humidity netatmo.ModuleData
	if sourceHumidity == "" {
		humidity = humidityModules[0]
	} else {
		humidity = filterModulesByName(netatmoModules, sourceHumidity)[0]
	}

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
