package aprs

import (
	"fmt"

	"github.com/ebarkie/aprs"
	"github.com/ottolote/naaprs/pkg/weather"
)

func SendWeatherData(wr *weather.WeatherReport) string {
	wx := &aprs.Wx{
		Timestamp: 	 wr.Timestamp,
		Altimeter:       wr.Altimeter,
		Humidity:        wr.Humidity,
		Lat: 		 wr.Lat,
		Lon: 		 wr.Lon,
		RainLastHour:    wr.RainLastHour,
		RainLast24Hours: wr.RainLast24Hours,
		RainToday:       wr.RainToday,
		SolarRad:        wr.SolarRad,
		Temp:            wr.Temp,
		WindDir:         wr.WindDir,
		WindGust:        wr.WindGust,
		WindSpeed:       wr.WindSpeed,
	}

	return fmt.Sprintf("SendWeatherData is not implemented, would send weather packet:\n\t%+v", wx)
}
