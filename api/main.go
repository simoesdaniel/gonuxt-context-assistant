package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	// Required for unicode.IsSpace and unicode.ToUpper
	"gonuxt-context-assistant/api/internal/tools" // Make sure this path is correct for your module

	"github.com/rs/cors"
)

type RequestBody struct {
	Query string `json:"query"`
}

type ResponseBody struct {
	Answer string `json:"answer"`
}

// askHandler is the HTTP handler function for our /ask endpoint.
func askHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in askHandler: %v", r)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Ensure the request body is closed

	log.Printf("Received query: \"%s\"", reqBody.Query)

	var answer string
	query := reqBody.Query

	if contains(query, "time") || contains(query, "date") {
		// Log the tool invocation
		log.Println("Invoking GetCurrentDateTime tool.")
		answer = tools.GetCurrentDateTime()
	} else if contains(query, "weather") {
		city := extractCity(query)
		if city != "" {
			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second) // Set a timeout for the request context
			defer cancel()

			log.Printf("Invoking GetWeather tool for city: %s", city)
			weatherReport, _ := tools.GetWeather(ctx, city)
			answer = weatherReport
		} else {
			answer = "Please specify a city for weather information. E.g., 'What's the weather in London?'"
		}
	} else {
		answer = "Hello! I am a simple assistant. I can tell you the current time or the weather in a major city. Try asking me about 'time' or 'weather in London'."
	}

	respBody := ResponseBody{Answer: answer}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set status code before writing body

	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

type MultipleCityRequestBody struct {
	Query string `json:"query"`
}

type MultipleCityResponseBody struct {
	Reports map[string]string `json:"reports"`
}

type MultipleAsyncRequestBody struct {
	Cities []string `json:"cities"`
}

type MultipleAsyncResponseBody struct {
	Reports map[string]string `json:"reports"`
}

func askMultipleCityWeatherAsyncHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in askMultiCityWeatherHandler: %v", r)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // Set a timeout for the request context
	defer cancel()

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody MultipleAsyncRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Received multi-city query for cities: %v", reqBody.Cities)

	if len(reqBody.Cities) == 0 {
		respBody := MultipleAsyncResponseBody{Reports: map[string]string{"error": "No cities provided in the query."}}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // Bad Request for empty city list
		json.NewEncoder(w).Encode(respBody)
		return
	}

	// --- Concurrency Setup ---
	var wg sync.WaitGroup // A WaitGroup to wait for all goroutines to finish.
	// Make a buffered channel to receive results from goroutines.
	// The buffer size is the number of cities, so goroutines don't block
	// if the main goroutine isn't ready to receive immediately.
	resultsChan := make(chan struct { // Channel to send structs containing city and report
		City   string
		Report string
	}, len(reqBody.Cities))

	// Iterate over cities and launch a goroutine for each.
	for _, city := range reqBody.Cities {
		wg.Add(1) // Increment the WaitGroup counter for each goroutine.

		// Launch a goroutine
		// `go` keyword followed by a function call.
		// It's a best practice to pass loop variables as arguments to goroutines
		// to avoid issues with variable capture in closures.
		go func(currentCity string) {
			defer wg.Done() // Decrement the WaitGroup counter when this goroutine finishes.

			log.Printf("Goroutine started for city: %s", currentCity)
			weatherReport, found := tools.GetWeather(ctx, currentCity) // Call our tool
			log.Printf("Goroutine finished for city: %s", currentCity)

			if !found {
				weatherReport = fmt.Sprintf("Weather data for %s could not be found.", currentCity)
			}
			// Send the result back on the channel.
			resultsChan <- struct {
				City   string
				Report string
			}{City: currentCity, Report: weatherReport}

		}(city) // Pass 'city' as an argument to the anonymous function.
	}

	// This goroutine waits for all other goroutines to complete, then closes the channel.
	go func() {
		wg.Wait()          // Block until all wg.Done() calls have been made.
		close(resultsChan) // Close the channel when all results have been sent.
	}()

	// Collect results from the channel.
	reports := make(map[string]string)
	for res := range resultsChan { // Loop until the channel is closed.
		reports[res.City] = res.Report
	}

	// --- Prepare and Send Response (unchanged) ---
	respBody := MultipleAsyncResponseBody{Reports: reports}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// New handler for multi-city weather queries
func askMultiCityWeatherHandler(w http.ResponseWriter, r *http.Request) {
	// --- Error Handling Best Practice: Defer for Panic Recovery ---
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in askMultiCityWeatherHandler: %v", r)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	// --- Check HTTP Method ---
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// --- Decode Request Body ---
	var reqBody MultipleCityRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Received multi-city query: %s", reqBody.Query)

	// --- Initialize response map ---
	reports := make(map[string]string)

	// This is where our loop will come in (next step)!
	// For now, we'll just acknowledge the cities.
	if len(reqBody.Query) == 0 {
		reports["error"] = "No cities provided in the query."
	} else {

		ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second) // Set a timeout for the request context
		defer cancel()                                                 // Ensure we cancel the context to release resources

		// reports["message"] = fmt.Sprintf("Processing weather for %d cities...", len(reqBody.Cities))
		// Placeholder for loop and weather fetching
		cities := tools.ExtractCitiesFromQuery(reqBody.Query)
		reports, err = tools.GetWeatherForCities(ctx, cities)
		if err != nil {
			log.Printf("Error fetching weather reports: %v", err)
			http.Error(w, "Error fetching weather reports: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// --- Prepare and Send Response ---
	respBody := MultipleCityResponseBody{Reports: reports}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// main function is the entry point of our server application.
func main() {
	// 1. Create a new HTTP multiplexer (router). This is best practice for custom routing.
	// We could use http.DefaultServeMux, but creating our own gives more control.
	mux := http.NewServeMux()

	// 2. Register our askHandler with the multiplexer.
	// http.HandlerFunc(askHandler) converts the function into an http.Handler.
	mux.Handle("/ask", http.HandlerFunc(askHandler))
	mux.Handle("/ask-multiple-city-weather", http.HandlerFunc(askMultiCityWeatherHandler))
	mux.Handle("/ask-multi-city-weather-async", http.HandlerFunc(askMultipleCityWeatherAsyncHandler))

	// 3. Create the CORS middleware instance.
	// The `cors` package expects an `http.Handler` to wrap.
	// We wrap our `mux` (which is an http.Handler).
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow Nuxt.js dev server
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
		Debug:          true, // Enable CORS logging for debugging
	}).Handler(mux) // <--- Correct usage: wrap the mux (router)

	// 4. Start the HTTP server with the CORS-wrapped handler.
	fmt.Println("Server starting on port 8080...")
	// We pass our `handler` (which is the mux wrapped by CORS) to ListenAndServe.
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// Helper functions (kept outside main for clarity and reusability within this package)

// contains checks if s contains substr, case-insensitively.
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// extractCity attempts to extract a city name from a query string.
func extractCity(query string) string {
	lowerQuery := strings.ToLower(query)
	city := ""

	keywords := []string{"weather in ", "weather for ", "weather of ", "weather "}
	// extract logic if city ends with ?
	// Remove trailing '?' or '.' if present
	if strings.HasSuffix(lowerQuery, "?") || strings.HasSuffix(lowerQuery, ".") {
		lowerQuery = strings.TrimRight(lowerQuery, "?.")
	}

	for _, keyword := range keywords {
		if strings.Contains(lowerQuery, keyword) {
			parts := strings.SplitN(lowerQuery, keyword, 2)
			log.Printf("Extracted parts: %v", parts)
			if len(parts) > 1 {
				city = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if city != "" {
		// strings.Title is deprecated, but simpler for this example.
		// For robust i18n, use golang.org/x/text/cases.Title.
		return strings.Title(city)
	}
	return ""
}
