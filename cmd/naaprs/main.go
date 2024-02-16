package main

import (
	"fmt"

	"github.com/ottolote/naaprs/pkg/aprs"
	"github.com/ottolote/naaprs/pkg/weather"
)

var Version string

func main() {
	if Version == "" {
		Version = "unset"
	}
	fmt.Printf("naaprs started version: %s\n", Version)

	currentWeather := weather.GetWeatherData("netatmo")
	fmt.Printf("weather: %+v\n", currentWeather)

	fmt.Printf("aprs-connection: %+v\n", aprs.SendWeatherData(currentWeather))
	return
}
