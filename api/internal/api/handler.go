package api

import (
	"context"
	"encoding/json"
	"gonuxt-context-assistant/internal/app/assistant"
	"log"
	"net/http"
	"time"
)

// Handler struct to hold dependencies like the assistant service (if needed later)
// This is a common pattern for injecting dependencies into handlers.
type Handler struct {
	Assistant *assistant.Service // The core assistant service
}

// NewHandler creates a new Handler instance.
func NewHandler(svc *assistant.Service) *Handler {
	return &Handler{
		Assistant: svc,
	}
}

// askHandler is the HTTP handler function for our /ask endpoint.
// *http is a pointer to the http.ResponseWriter, which is used to write the response.
func (h *Handler) AskHandler(w http.ResponseWriter, r *http.Request) {
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

	answer, httpStatus := h.Assistant.ProcessQuery(r.Context(), query)

	respBody := ResponseBody{Answer: answer}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus) // Set status code before writing body

	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (h *Handler) AskMultipleCityWeatherAsyncHandler(w http.ResponseWriter, r *http.Request) {
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

	reports, httpStatus := h.Assistant.GetMultiCityWeather(ctx, reqBody.Cities)

	// --- Prepare and Send Response (unchanged) ---
	respBody := MultipleAsyncResponseBody{Reports: reports}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// New handler for multi-city weather queries
func (h *Handler) AskMultiCityWeatherFromQueryHandler(w http.ResponseWriter, r *http.Request) {
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

	// This is where our loop will come in (next step)!
	// For now, we'll just acknowledge the cities.
	if len(reqBody.Query) == 0 {
		http.Error(w, "No cities provided in the query.", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second) // Set a timeout for the request context
	defer cancel()                                                 // Ensure we cancel the context to release resources

	// reports["message"] = fmt.Sprintf("Processing weather for %d cities...", len(reqBody.Cities))
	// Placeholder for loop and weather fetching
	reports, httpStatus := h.Assistant.GetWeatherForCitiesFromQuery(ctx, reqBody.Query)
	if httpStatus != http.StatusOK {
		http.Error(w, "Error fetching weather reports", http.StatusInternalServerError)
		return
	}

	// --- Prepare and Send Response ---
	respBody := MultipleCityResponseBody{Reports: reports}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
