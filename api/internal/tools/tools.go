package tools // Declares this file as part of the "tools" package.

import (
	"fmt"  // Used for formatted string output, like Sprintf.
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

// GetWeather returns a simulated weather report for a given city.
// This is also a public function.
// In a real application, this would make an API call to a weather service.
func GetWeather(city string) string {
	// Here, we're just providing a mock response.
	// We use fmt.Sprintf to easily insert the 'city' variable into the string.
	// This demonstrates string formatting, similar to f-strings in Python or printf in C.
	switch city {
	case "Lisbon", "lisbon": // Case-insensitive check for common city names
		return fmt.Sprintf("The weather in %s is currently sunny with 28°C. Don't forget your sunglasses!", city)
	case "London", "london":
		return fmt.Sprintf("The weather in %s is currently cloudy with 18°C. A typical British day!", city)
	case "New York", "new york", "NYC":
		return fmt.Sprintf("The weather in %s is currently partly cloudy with 22°C. Expect some humidity.", city)
	default:
		// Default response if the city is not recognized in our mock data.
		return fmt.Sprintf("I don't have weather information for %s at the moment. Please try another major city.", city)
	}
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
