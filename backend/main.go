package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"yolink-app/yolink"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

func main() {
	// Create a new HTTP server mux
	mux := http.NewServeMux()

	// Serve static files from the frontend build directory
	fs := http.FileServer(http.Dir(filepath.Join("..", "frontend", "public")))
	mux.Handle("/", fs)

	// API endpoints
	mux.HandleFunc("/api/hello", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello from Go backend!"}`))
	}))

	mux.HandleFunc("/api/devices", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		devices, err := yolink.GetDevices()
		if err != nil {
			log.Printf("Error getting devices: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(devices); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}))

	mux.HandleFunc("/api/devices/state", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		deviceID := r.URL.Query().Get("deviceId")
		deviceType := r.URL.Query().Get("deviceType")
		if deviceID == "" || deviceType == "" {
			http.Error(w, "Device ID and type are required", http.StatusBadRequest)
			return
		}

		state, err := yolink.GetDeviceState(deviceID, deviceType)
		if err != nil {
			log.Printf("Error getting device state: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(state); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	}))

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
} 