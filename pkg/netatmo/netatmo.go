package netatmo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2"
)

type config struct {
	auth *oauth2.Config

	clientID     string
	clientSecret string
	refreshToken string
}

type ApiResponse struct {
	Body struct {
		Devices []Device `json:"devices"`
	} `json:"body"`
}

type Device struct {
	StationName   string        `json:"station_name"`
	HomeName      string        `json:"home_name"`
	HomeID        string        `json:"home_id"`
	Place         Place         `json:"place"`
	DashboardData DashboardData `json:"dashboard_data"`
	Modules       []Module      `json:"modules"`
}

type Place struct {
	Altitude    int 	`json:"altimeter"`
	City        string      `json:"city"`
	Country     string      `json:"country"`
	Location    []float64   `json:"location"`
}

type DashboardData struct {
	TimeUTC          int64   `json:"time_utc"`
	Temperature      float64 `json:"Temperature"`
	CO2              int     `json:"CO2"`
	Humidity         int     `json:"Humidity"`
	Pressure         float64 `json:"Pressure"`
	AbsolutePressure float64 `json:"AbsolutePressure"`
}

type Module struct {
	ID            string              `json:"_id"`
	Type          string              `json:"type"`
	ModuleName    string              `json:"module_name"`
	DataType      []string            `json:"data_type"`
	DashboardData ModuleDashboardData `json:"dashboard_data"`
}

type ModuleDashboardData struct {
	TimeUTC      int64   `json:"time_utc"`
	Rain         float64 `json:"Rain,omitempty"`
	SumRain1     float64 `json:"sum_rain_1,omitempty"`
	SumRain24    float64 `json:"sum_rain_24,omitempty"`
	Temperature  float64 `json:"Temperature,omitempty"`
	Humidity     int     `json:"Humidity,omitempty"`
	WindStrength int     `json:"WindStrength,omitempty"`
	WindAngle    int     `json:"WindAngle,omitempty"`
	GustStrength int     `json:"GustStrength,omitempty"`
	GustAngle    int     `json:"GustAngle,omitempty"`
}

type ModuleData struct {
	StationName     string
	HomeName        string
	HomeID          string
	Lat 		float64
	Lon 		float64
	ModuleID        string
	ModuleName      string
	ModuleType      string
	DataType        []string
	Timestamp       time.Time
	Altimeter       float64
	Humidity        int
	RainLastHour    float64
	RainLast24Hours float64
	Temp            float64
	WindDir         int
	WindGust        int
	WindSpeed       int
}

type NetatmoClient struct {
	client  *http.Client
	baseURL string
}

func (nc *NetatmoClient) unmarshalModuleData(body []byte) ([]ModuleData, error) {
	var response ApiResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}

	var weatherData []ModuleData
	for _, device := range response.Body.Devices {
		for _, module := range device.Modules {
			wd := ModuleData{
				StationName:     device.StationName,
				HomeName:        device.HomeName,
				HomeID:          device.HomeID,
				Lat: 		 device.Place.Location[1],
				Lon:		 device.Place.Location[0],
				ModuleID:        module.ID,
				ModuleName:      module.ModuleName,
				ModuleType:      module.Type,
				DataType:        module.DataType,
				Timestamp:       time.Unix(module.DashboardData.TimeUTC, 0),
				Altimeter:       device.DashboardData.AbsolutePressure,
				Humidity:        module.DashboardData.Humidity,
				RainLastHour:    module.DashboardData.Rain,
				RainLast24Hours: module.DashboardData.SumRain24,
				Temp:            module.DashboardData.Temperature,
				WindDir:         module.DashboardData.WindAngle,
				WindGust:        module.DashboardData.GustStrength,
				WindSpeed:       module.DashboardData.WindStrength,
			}
			weatherData = append(weatherData, wd)
		}
	}

	return weatherData, nil
}

func (nc *NetatmoClient) GetModuleData() ([]ModuleData, error) {
	req, err := nc.newRequest("GET", "/api/getstationsdata")
	if err != nil {
		return nil, err
	}

	resp, err := nc.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	moduleData, err := nc.unmarshalModuleData(body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("moduledata: %+v", moduleData)

	return moduleData, nil
}

func (nc *NetatmoClient) newRequest(method, relativePath string) (*http.Request, error) {
	// Parse the base URL
	parsedBaseURL, err := url.Parse(nc.baseURL)
	if err != nil {
		return nil, err
	}

	// Parse the relative path
	parsedRelPath, err := url.Parse(relativePath)
	if err != nil {
		return nil, err
	}

	// Resolve the relative path against the base URL
	fullURL := parsedBaseURL.ResolveReference(parsedRelPath).String()

	// Create the request
	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func getEnvOrError(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", errors.New("environment variable not set: " + name)
	}
	return value, nil
}

func GetConfig() *config {
	oauth := &oauth2.Config{
		ClientID:     os.Getenv("NETATMO_CLIENT_ID"),
		ClientSecret: os.Getenv("NETATMO_CLIENT_SECRET"),
		Scopes:       []string{"read_station"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.netatmo.com/oauth2/authorize",
			TokenURL: "https://api.netatmo.com/oauth2/token",
		},
	}

	clientID, err := getEnvOrError("NETATMO_CLIENT_ID")
	if err != nil {
		log.Fatalf("Missing configuration: %v", err)
	}
	clientSecret, err := getEnvOrError("NETATMO_CLIENT_SECRET")
	if err != nil {
		log.Fatalf("Missing configuration: %v", err)
	}
	refreshToken, err := getEnvOrError("NETATMO_REFRESH_TOKEN")
	if err != nil {
		log.Fatalf("Missing configuration: %v", err)
	}

	return &config{
		auth: oauth,

		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
	}
}

func (c *config) GetNetatmoClient(ctx context.Context) *NetatmoClient {

	netatmoClient := &NetatmoClient{
		baseURL: "https://api.netatmo.com/api",
	}

	c.auth = &oauth2.Config{
		ClientID:     os.Getenv("NETATMO_CLIENT_ID"),
		ClientSecret: os.Getenv("NETATMO_CLIENT_SECRET"),
		Scopes:       []string{"read_station"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.netatmo.com/oauth2/authorize",
			TokenURL: "https://api.netatmo.com/oauth2/token",
		},
	}

	userRefreshToken := os.Getenv("NETATMO_REFRESH_TOKEN")
	if userRefreshToken == "" {
		log.Fatalf("No refresh token supplied")
	}

	tokenSource := c.auth.TokenSource(ctx, &oauth2.Token{
		RefreshToken: userRefreshToken,
	})

	tok, err := tokenSource.Token() // This will automatically refresh the token if it's expired
	if err != nil {
		log.Fatalf("Error getting token from token source: %v", err)
	}

	// Create an HTTP client using the token
	netatmoClient.client = c.auth.Client(ctx, tok)

	return netatmoClient
}

func GetAllModules() []ModuleData {
	ctx := context.Background()

	c := GetConfig()
	client := c.GetNetatmoClient(ctx)

	moduleData, err := client.GetModuleData()
	if err != nil {
		log.Fatalf("Error getting stationsdata: %v", err)
	}

	return moduleData
}
