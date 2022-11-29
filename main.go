package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type ApiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (ApiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return ApiConfigData{}, err
	}

	var c ApiConfigData
	if err := json.Unmarshal(bytes, &c); err != nil {
		return ApiConfigData{}, err
	}

	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

func query(city string) (WeatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return WeatherData{}, err
	}

	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return WeatherData{}, err
	}

	defer res.Body.Close()
	var d WeatherData
	if err := json.NewDecoder(res.Body).Decode(&d); err != nil {
		return WeatherData{}, err
	}

	return d, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})

	http.ListenAndServe(":8080", nil)
}
