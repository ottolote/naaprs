package main

import (
	"fmt"

	"github.com/ottolote/naaprs/pkg/weather"
	"github.com/ottolote/naaprs/pkg/aprs"
)

var Version string


func main() {
	fmt.Println("naaprs version ", Version)

	currentWeather := weather.GetWeatherData("netatmo")
	fmt.Printf("weather: %+v\n", currentWeather)

	fmt.Printf("aprs-connection: %+v\n", aprs.SendWeatherData(currentWeather))
	return
}
