package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

func main() {
	// Create an instance of multiWeatherProvider with two weather APIs
	mw := multiWeatherProvider{
		openWeatherMap{
			apiKey: "YOUR_API_KEY", // API key for openWeatherMap
		},
		weatherApi{
			apiKey: "YOUR_API_KEY", // API key for weatherApi
		},
	}

	// Handle HTTP requests to /weather/ endpoint
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now() // Record the start time to calculate request duration
		// Extract the city name from the URL path
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		// Get the temperature from the multiWeatherProvider
		temp, err := mw.temperature(city)
		if err != nil {
			// Return an internal server error if fetching temperature fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set response header and send the JSON response with city, temperature, and duration
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"city": city,
			"temp": temp,
			"took": time.Since(begin).String(), // Calculate how long the request took
		})
	})

	// Start the HTTP server on port 8080
	http.ListenAndServe(":8080", nil)
}

// Struct to represent the openWeatherMap API with its API key
type openWeatherMap struct {
	apiKey string
}

// Fetch the temperature for a city from the openWeatherMap API
func (w openWeatherMap) temperature(city string) (float64, error) {
	// Perform a GET request to the openWeatherMap API with the city and API key
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + w.apiKey + "&q=" + city)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close() // Ensure the response body is closed

	// Struct to parse the JSON response from the API
	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"` // Temperature is provided in Kelvin
		} `json:"main"`
	}

	// Decode the JSON response into the struct
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	// Convert the temperature from Kelvin to Fahrenheit
	fahrenheit := math.Round((d.Main.Kelvin-273.15)*9/5 + 32)

	// Log the city and its temperature in Fahrenheit
	log.Printf("openWeatherMap: %s: %.2f", city, fahrenheit)
	return fahrenheit, nil
}

// Struct to represent the weatherApi with its API key
type weatherApi struct {
	apiKey string
}

// Fetch the temperature for a city from the weatherApi API
func (w weatherApi) temperature(city string) (float64, error) {
	// Perform a GET request to the weatherApi with the city and API key
	resp, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + w.apiKey + "&q=" + city + "&aqi=no")
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close() // Ensure the response body is closed

	// Struct to parse the JSON response from the API
	var d struct {
		Current struct {
			Fahrenheit float64 `json:"temp_f"` // Temperature is provided in Fahrenheit
		} `json:"current"`
	}

	// Decode the JSON response into the struct
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	// Round the temperature to the nearest whole number
	fahrenheit := math.Round(d.Current.Fahrenheit)
	// Log the city and its temperature in Fahrenheit
	log.Printf("weatherApi: %s: %.2f", city, fahrenheit)
	return fahrenheit, nil
}

// Interface for weather providers that can return the temperature for a city
type weatherProvider interface {
	temperature(city string) (float64, error)
}

// multiWeatherProvider is a slice of weatherProvider interfaces
type multiWeatherProvider []weatherProvider

// Fetch the temperature for a city using multiple weather providers and return the average
func (w multiWeatherProvider) temperature(city string) (float64, error) {
	// Make a channel for temperatures, and a channel for errors.
	// Each provider will push a value into only one.
	temps := make(chan float64, len(w))
	errs := make(chan error, len(w))

	// For each provider, spawn a goroutine with an anonymous function.
	// That function will invoke the temperature method, and forward the response.
	for _, provider := range w {
		go func(p weatherProvider) {
			k, err := p.temperature(city)
			if err != nil {
				errs <- err // Push error to the errs channel
				return
			}
			temps <- k // Push temperature to the temps channel
		}(provider)
	}

	sum := 0.0 // Initialize a variable to hold the sum of temperatures

	// Collect a temperature or an error from each provider.
	for i := 0; i < len(w); i++ {
		select {
		case temp := <-temps: // Receive temperature from temps channel
			sum += temp // Add temperature to the sum
		case err := <-errs: // If an error is received, return the error
			return 0, err
		}
	}

	// Return the average of the temperatures rounded to the nearest whole number
	return math.Round(sum / float64(len(w))), nil
}
