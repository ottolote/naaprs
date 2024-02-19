package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ottolote/naaprs/pkg/aprs"
	"github.com/ottolote/naaprs/pkg/weather"

	"github.com/spf13/viper"
)

var Version string

var Logo string = `
   ____  ____ _____ _____  __________
  / __ \/ __ '/ __ '/ __ \/ ___/ ___/
 / / / / /_/ / /_/ / /_/ / /  (__  )
/_/ /_/\__,_/\__,_/ .___/_/  /____/
                 /_/

`

type ConfigOption struct {
	ConfigKey string
	Required  bool
}

func configure() {
	options := []ConfigOption{
		{"CALLSIGN", true},
		{"INTERVAL", false},
		{"NETATMO_CLIENT_ID", true},
		{"NETATMO_CLIENT_SECRET", true},
		{"NETATMO_REFRESH_TOKEN", true},
		{"SOURCE_TEMPERATURE", false},
		{"SOURCE_HUMIDITY", false},
		{"SOURCE_RAIN", false},
		{"SOURCE_WIND", false},
		{"ONESHOT", false},
		{"DRY_RUN", false},
	}

	for _, option := range options {
		key := option.ConfigKey
		viper.BindEnv(key)
		if option.Required {
			if viper.Get(key) == nil {
				panic(fmt.Sprintf("missing configuration: %s", key))
			}
		}
	}
}

func main() {
	configure()

	fmt.Printf(Logo)

	if Version == "" {
		Version = "unset"
	}
	log.Printf("Started version: %s\n", Version)

	for {
		currentWeather := weather.GetWeatherData("netatmo")
		log.Printf("Successfully got weather report from Netatmo")

		aprs.SendWeatherData(currentWeather)

		if viper.Get("ONESHOT") != nil {
			log.Printf("Running in oneshot mode, exiting after one packet")
			return
		}

		interval := viper.GetInt("INTERVAL")
		if interval == 0 {
			interval = 600
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return
}
