package main

import (
	"fmt"

	"github.com/ottolote/naaprs/pkg/weather"
	"github.com/ottolote/naaprs/pkg/aprs"
)


func main() {
	currentWeather := weather.GetWeatherData("netatmo")
	fmt.Printf("weather: %+v\n", currentWeather)

	fmt.Printf("aprs-connection: %+v\n", aprs.SendWeatherData(currentWeather))
	return
}
