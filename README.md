# GoWeatherAPI - Simple Weather API Server
GoWeatherAPI is a lightweight HTTP server written in Go that retrieves the current temperature for a given city. It leverages multiple weather API providers (such as OpenWeatherMap and WeatherAPI) to gather temperature data, calculate the average, and return it as a JSON response. The server listens on port 8080 and handles requests in the format /weather/{city}.


## Features

- Fetches the current temperature for any city using multiple weather APIs.
- Combines data from two providers to calculate an average temperature.
- Responses are returned as JSON with city, temperature (in Fahrenheit), and request duration.
- Simple and fast HTTP server setup with Go.

## Requirements

- Go 1.16 or higher
- Internet connection to fetch weather data from external APIs
- API keys for both OpenWeatherMap and WeatherAPI (see setup instructions below)

## Setup

### 1. Clone the Repository

Clone the repository using Git and navigate into the project directory.

### 2. Get API Keys

- **OpenWeatherMap:** Sign up at [OpenWeatherMap](https://openweathermap.org/api) and obtain your API key.
- **WeatherAPI:** Sign up at [WeatherAPI](https://www.weatherapi.com/) and obtain your API key.

### 3. Update API Keys in Code

In the `main.go` file, update the following placeholders with your actual API keys:

- Replace the `apiKey` in the `openWeatherMap` struct.
- Replace the `apiKey` in the `weatherApi` struct.

### 4. Run the Application

Compile and run the Go application. By default, it will start a server on port 8080.
```bash
go run main.go
```

## Usage

The API exposes a single endpoint to get the current temperature for a given city.

- **Endpoint:** `/weather/{city}`
- **Method:** GET

Example request: `/weather/London`

The response will include:

- The city name
- The average temperature (rounded to the nearest integer)
- The time taken to fetch the result

### Example Response

```json
{
  "city": "London",
  "temp": 63,
  "took": "500ms"
}

```

## Logging

The application logs the temperature fetched from each provider along with the city name to provide visibility into the data source and the calculated temperature.

## License
This project is licensed under the MIT License.