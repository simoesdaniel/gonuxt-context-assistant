package main

import (
	"fmt"
	"log"
	"net/http"

	// Required for unicode.IsSpace and unicode.ToUpper
	// Make sure this path is correct for your module

	"gonuxt-context-assistant/internal/api"
	"gonuxt-context-assistant/internal/app/assistant"

	"github.com/rs/cors"
)

// main function is the entry point of our server application.
func main() {

	// Initialize the core assistant service
	assistantSvc := assistant.NewService()

	// Initialize the API handlers, injecting the assistant service
	apiHandlers := api.NewHandler(assistantSvc)

	// 1. Create a new HTTP multiplexer (router). This is best practice for custom routing.
	// We could use http.DefaultServeMux, but creating our own gives more control.
	mux := http.NewServeMux()

	// 2. Register our askHandler with the multiplexer.
	// http.HandlerFunc(askHandler) converts the function into an http.Handler.
	mux.Handle("/ask", http.HandlerFunc(apiHandlers.AskHandler))
	mux.Handle("/ask-multiple-city-weather", http.HandlerFunc(apiHandlers.AskMultiCityWeatherFromQueryHandler))
	mux.Handle("/ask-multi-city-weather-async", http.HandlerFunc(apiHandlers.AskMultipleCityWeatherAsyncHandler))

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
