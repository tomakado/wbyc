package apixu

import (
	"../geocoder"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Weather struct {
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
}

func GetCurrentWeather(coordinates geocoder.Coordinates) (*Weather, error) {
	current, err := fetchCurrentWeather(coordinates)
	if err != nil {
		return &Weather{}, err
	}

	return &Weather{
		Temperature: current.TempC,
		FeelsLike:   current.FeelslikeC,
	}, nil
}

func fetchCurrentWeather(coordinates geocoder.Coordinates) (*Current, error) {
	err := godotenv.Load()
	if err != nil {
		return &Current{}, errors.New("Failed to load environment variables. Make sure .env file exists")
	}

	apiKey := os.Getenv("APIXU_API_KEY")
	urlFormat := "http://api.apixu.com/v1/current.json?key=%s&q=%f,%f"

	// Idk why Apixu waits {lon, lat} pair instead of {lat, lon} ¯\_(ツ)_/¯
	url := fmt.Sprintf(urlFormat, apiKey, coordinates.Longitude, coordinates.Latitude)

	resp, err := http.Get(url)
	if err != nil {
		return &Current{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Current{}, err
	}

	var response HttpApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return &Current{}, errors.New("Failed to parse response as JSON data")
	}

	return &response.Current, nil
}
