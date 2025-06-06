package assistant

import (
	"context" // Important for context propagation
	"fmt"
	"log"
	"net/http" // For HTTP status codes
	"strings"
	"sync" // For sync.WaitGroup
	"time"

	"gonuxt-context-assistant/internal/tools" // Import our tools
)

// Service defines the core assistant logic.
// This struct would hold dependencies like database clients, external API clients, etc.
type Service struct {
	// Add any dependencies here, e.g., Logger *log.Logger
	Logger *log.Logger // Optional: if you want to log within the service
}

// NewService creates a new instance of the Assistant Service.
func NewService() *Service {
	return &Service{}
}

// ProcessQuery takes a context and a query string, returning the answer and an HTTP status code.
func (s *Service) ProcessQuery(ctx context.Context, query string) (string, int) {
	var answer string
	if contains(query, "time") || contains(query, "date") {
		answer = tools.GetCurrentDateTime()
	} else if contains(query, "weather") {
		city := extractCity(query)
		if city != "" {
			ctx, cancel := context.WithTimeout(ctx, 3*time.Second) // Set a timeout for the request context
			defer cancel()

			log.Printf("Invoking GetWeather tool for city: %s", city)
			weatherReport, _ := tools.GetData(ctx, city, tools.GetWeather)
			answer = weatherReport
		} else {
			answer = "Please specify a city for weather information. E.g., 'What's the weather in London?'"
		}
	} else {
		answer = "Hello! I am a simple assistant. I can tell you the current time or the weather in a major city. Try asking me about 'time' or 'weather in London'."
	}

	return answer, http.StatusOK
}

// GetMultiCityWeather takes a context and a slice of city names, returning a map of reports and an HTTP status code.
func (s *Service) GetMultiCityWeather(ctx context.Context, cities []string) (map[string]string, int) {
	reports := make(map[string]string)
	var wg sync.WaitGroup

	type cityReport struct {
		City   string
		Report string
	}
	resultsChan := make(chan cityReport, len(cities))

	for _, city := range cities {
		wg.Add(1)
		go func(currentCity string) {
			defer wg.Done()
			// Pass the context received by GetMultiCityWeather down to GetData
			result, err := tools.GetData(ctx, currentCity, tools.GetWeather) // Reuse GetWeather via GetData
			if err != nil {
				result = fmt.Sprintf("Weather data for %s could not be found.", currentCity)
			}
			resultsChan <- cityReport{City: currentCity, Report: result}
		}(city)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for res := range resultsChan {
		reports[res.City] = res.Report
	}

	return reports, http.StatusOK // If all individual requests handle their own errors, overall OK
}

func (s *Service) GetWeatherForCitiesFromQuery(ctx context.Context, query string) (map[string]string, int) {

	cities := ExtractCitiesFromQuery(query) // Extract cities from the query using a helper function.

	reports, err := tools.GetWeatherForCities(ctx, cities) // Return the result of the private function.
	if err != nil {
		log.Printf("Error fetching weather reports: %v", err)
		return nil, http.StatusInternalServerError
	}
	return reports, http.StatusOK // Return the reports and HTTP status OK.
}

// --- Helper functions (copy from your old main.go if they were there) ---
// You might put these in a separate internal/util package or keep them private to assistant package.

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func extractCity(query string) string {
	// Simplified extraction, actual LLM parsing would be more complex
	lowerQuery := strings.ToLower(query)
	cities := []string{"lisbon", "london", "new york", "paris", "tokyo", "porto"}
	for _, city := range cities {
		if strings.Contains(lowerQuery, city) {
			return strings.Title(city) // Capitalize for consistency
		}
	}
	return ""
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
