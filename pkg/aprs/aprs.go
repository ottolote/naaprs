package aprs

import (
	"fmt"
	"log"

	"github.com/ottolote/naaprs/pkg/weather"

	"github.com/ebarkie/aprs"
	"github.com/spf13/viper"
)

func kmhToMph(speed float64) float64 {
	return speed / 1.609344
}

func celciusToFahrenheit(celcius float64) float64 {
	return (celcius * 1.8) + 32
}

func millimeterToInchHundredths(mm float64) float64 {
	return mm * 3.93700787402
}

func CreateWx(wr *weather.WeatherReport) *aprs.Wx {
	fmt.Printf("altimeter: %+v\n", wr.Altimeter)
	wx := &aprs.Wx{
		Type:            "natm",
		Timestamp:       wr.Timestamp,
		Altimeter:       wr.Altimeter * 0.0295301, // Convert mbar to inches of mercury
		Humidity:        int(wr.Humidity),
		Lat:             wr.Lat,
		Lon:             wr.Lon,
		RainLastHour:    millimeterToInchHundredths(wr.RainLastHour),
		RainLast24Hours: millimeterToInchHundredths(wr.RainLast24Hours),
		RainToday:       millimeterToInchHundredths(wr.RainToday),
		Temp:            int(celciusToFahrenheit(wr.Temp)),
		WindDir:         int(wr.WindDir),
		WindGust:        int(kmhToMph(float64(wr.WindGust))),
		WindSpeed:       int(kmhToMph(float64(wr.WindSpeed))),
	}
	return wx
}

func SendWeatherData(wr *weather.WeatherReport) {
	wx := CreateWx(wr)

	callsign := viper.GetString("CALLSIGN")
	tocall := viper.GetString("TOCALL")
	if tocall == "" {
		tocall = "APZ001"
	}

	f := aprs.Frame{
		Dst:  aprs.Addr{Call: tocall},
		Src:  aprs.Addr{Call: fmt.Sprintf("%s-13", callsign)},
		Path: aprs.Path{aprs.Addr{Call: "TCPIP", Repeated: true}},
		Text: wx.String(),
	}

	if viper.Get("DRY_RUN") != nil {
		log.Printf("Dry run, would send packet: %s", f)
		return
	}

	err := f.SendIS("tcp://cwop.aprs.net:14580", -1)
	if err != nil {
		log.Printf("Upload error: %s", err)
	}
	log.Printf("Sent weather packet: %s", wx)

}
