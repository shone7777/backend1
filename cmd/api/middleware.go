package main

import "net/http"

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allow necessary methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Allow Content-Type header

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK) // Send OK for OPTIONS requests
			return
		}

		next.ServeHTTP(w, r) // Proceed to the next handler for other requests
	})
}
