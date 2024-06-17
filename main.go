package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Output struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func main() {
	http.HandleFunc("/fetch-temperature", FetchZipCode)
	http.ListenAndServe(":8080", nil)
}

func FetchZipCode(w http.ResponseWriter, r *http.Request) {
	// Get the zipcode from the request parameters
	zipcode := r.URL.Query().Get("zipcode")
	if len(zipcode) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	// Make a GET request to the ViaCEP API
	resp, err := http.Get("https://viacep.com.br/ws/" + zipcode + "/json/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch zipcode information", http.StatusInternalServerError)
		return
	}

	// Decode the response body into a struct
	var address struct {
		Localidade string `json:"localidade"`
	}
	err = json.NewDecoder(resp.Body).Decode(&address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := FetchLocationTemperature(address.Localidade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(output)
}

func FetchLocationTemperature(location string) (*Output, error) {
	// Make a GET request to the weather API
	key := "fc7e61344b624b8c87601845241706"

	url, err := url.Parse(fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", key, url.QueryEscape(location)))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(&http.Request{Method: http.MethodGet, URL: url})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode the response body into a struct
	var weather struct {
		Current struct {
			Celsius    float64 `json:"temp_C"`
			Fahrenheit float64 `json:"temp_F"`
		} `json:"current"`
	}

	err = json.Unmarshal(bodyBytes, &weather)
	if err != nil {
		return nil, err
	}

	return &Output{
		Celsius:    weather.Current.Celsius,
		Fahrenheit: weather.Current.Fahrenheit,
		Kelvin:     weather.Current.Celsius + 273}, nil
}
