package tools // Declares this file as part of the "tools" package.

import (
	"context"
	"fmt" // Used for formatted string output, like Sprintf.
	"log"
	"strings"
	"time" // Provides functionality for working with time.
)

// GetCurrentDateTime returns the current date and time as a formatted string.
// This is a public function because its name starts with an uppercase letter.
func GetCurrentDateTime() string {
	// time.Now() gets the current local time.
	// Format() formats the time according to the provided layout string.
	// The layout "2006-01-02 15:04:05" is a magic number in Go's time package,
	// representing a specific reference date (Jan 2, 2006, 3:04:05 PM).
	// You use this specific reference to define the desired output format.
	return time.Now().Format("Current time is Monday, January 2, 2006 at 15:04:05 PM (MST)")
}

func ExtractCitiesFromQuery(query string) []string {
	knownCities := []string{"Lisbon", "London", "New York", "Paris", "Berlin", "Madrid"}
	var foundCities []string

	for _, city := range knownCities {
		if strings.Contains(strings.ToLower(query), strings.ToLower(city)) {
			foundCities = append(foundCities, city)
		}
	}
	return foundCities
}

// GetWeather returns a simulated weather report for a given city.
// This is also a public function.
// In a real application, this would make an API call to a weather service.
func GetWeather(ctx context.Context, city string) (string, bool) {
	select {
	case <-ctx.Done():
		fmt.Printf("GetWeather for %s cancelled! Reason: %v\n", city, ctx.Err())
		return "", false // If the context is cancelled, return an empty string and false.
	default: // Simulate quick I/O delay
		// Original weather logic
		weatherData := map[string]string{
			"lisbon":   "sunny with 28°C.",
			"london":   "cloudy with 18°C.",
			"new york": "partly cloudy with 22°C.",
			"paris":    "a delightful 20°C.",
			"tokyo":    "rainy with 15°C.", // Added for variety
		}
		log.Printf("GetWeather called for %s, arg: %s", city, strings.ToLower(city)) // Log the city being queried.
		report, found := weatherData[strings.ToLower(city)]
		if found {
			return fmt.Sprintf("The weather in %s is currently %s", city, report), true
		}
		return fmt.Sprintf("No weather information found for %s.", city), false
	}

}
func GetWeatherForCities(ctx context.Context, cities []string) (map[string]string, error) {
	// This function takes a slice of city names and returns a map with city names as keys
	// and their corresponding weather reports as values.
	reports := make(map[string]string) // Initialize an empty map to store the reports.

	for _, city := range cities {
		// For each city in the input slice, we call GetWeather to get the weather report.
		weatherReport, found := GetWeather(ctx, city)
		if found {
			reports[city] = weatherReport // Store the report in the map with the city name as the key.
		} else {
			reports[city] = fmt.Sprintf("Weather data for %s could not be found.", city)
		}
	}

	return reports, nil // Return the map containing all weather reports.
}

func GetCapital(country string) string {
	// This function returns the capital city of a given country.
	// It's a public function, so it can be accessed from other packages.
	return assessCapital(country) // Return the result of the private function.
}

func assessCapital(country string) string {
	// This is a private function, as it starts with a lowercase letter.
	// It can only be used within this package.
	switch country {
	case "Portugal":
		return "Lisbon"
	case "United Kingdom":
		return "London"
	case "United States":
		return "Washington, D.C."
	default:
		return fmt.Sprintf("I don't know the capital of %s. Please check your input.", country)
	}
}

// function with two arguments, go context and a method to call wich can be GetWeather or GetCapital
// arg can be either a string or a slice of strings depending on the method called
/*
Cons / Areas for Improvement (Go Idioms & Best Practices):

Use of interface{} for arg and return value:

Loss of Type Safety: This is the biggest issue. interface{} (the empty interface) can hold any type. While powerful, it pushes type checking from compile-time (where Go excels) to runtime.
Runtime Type Assertions (arg.([]string)): You're constantly performing type assertions (if cityNames, ok := arg.([]string); ok). This makes your code more brittle. If arg is not exactly the type you expect, ok will be false, and you'll fall into error paths, but the compiler can't help you.
Verbosity and Error Prone: It makes the code verbose with all the if/else if and ok checks, and it's easy to miss a type case or make a typo.
method as string and Manual Dispatch:

Fragile Dispatch: Using a string ("GetWeatherForCities") to determine which function to call is common in languages like JavaScript, but less idiomatic in Go. It's prone to typos (e.g., "GetWeatherForCitties" won't be caught by the compiler) and doesn't scale well.
No Compile-Time Checks: If you rename a function or change its signature, your GetData function won't know until runtime.
Centralizing All Calls:

While centralization has its benefits, it often leads to a "God function" that knows too much about many different operations.
Go generally favors smaller, single-purpose functions that are composed together.
Simulated Delay:

The time.After(3 * time.Second) applies a 3-second delay to every method call. This might not be desirable. You'd want delays/timeouts to be specific to the underlying operations (e.g., GetWeatherForCities might be slower due to multiple external calls).
Context-aware functions usually manage their own timeouts internally if they involve I/O. The GetData function's role should primarily be to pass context, not enforce uniform delays.

func GetData(ctx context.Context, method string, arg interface{}) (interface{}, error) {
	// This function takes a context and a method name (as a string) to call either GetWeather or GetCapital.
	// It returns the result of the method call or an error if the context is cancelled.

	select {
	case <-ctx.Done():
		return "", ctx.Err() // If the context is cancelled, return an error.
	case <-time.After(3 * time.Second): // Simulate a delay for the operation.
		if method == "GetWeatherForCities" {
			if cityNames, ok := arg.([]string); ok {
				return GetWeatherForCities(cityNames)
			}
			return "", fmt.Errorf("invalid argument type for GetWeatherForCities: %T", arg)
		} else if method == "GetCapital" {
			if countryName, ok := arg.(string); ok {
				return GetCapital(countryName), nil
			}
			return "", fmt.Errorf("invalid argument type for GetCapital: %T", arg)
		} else if method == "GetWeather" {
			if cityName, ok := arg.(string); ok {
				weatherReport, found := GetWeather(cityName)
				if found {
					return weatherReport, nil
				} else {
					return fmt.Sprintf("Weather data for %s could not be found.", cityName), nil
				}
			}
			return "", fmt.Errorf("invalid argument type for GetWeather: %T", arg)
		}
		return "", fmt.Errorf("unknown method: %s", method) // Return an error if the method is not recognized.
	}
}
*/
