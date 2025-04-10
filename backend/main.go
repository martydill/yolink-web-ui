package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go-svelte-app/backend/mqtt"
	"go-svelte-app/backend/websocket"
	"go-svelte-app/backend/yolink"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize MQTT client
	if err := mqtt.Init(); err != nil {
		log.Fatalf("Failed to initialize MQTT client: %v", err)
	}

	// Start WebSocket broadcaster
	go websocket.StartBroadcaster()

	// Create router
	r := mux.NewRouter()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// API routes
	r.HandleFunc("/api/devices", getDevices).Methods("GET")
	r.HandleFunc("/api/devices/state", getDeviceState).Methods("GET")
	r.HandleFunc("/ws", websocket.HandleConnections)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := yolink.GetDevices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(devices)
}

func getDeviceState(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("deviceId")
	deviceType := r.URL.Query().Get("deviceType")
	if deviceID == "" || deviceType == "" {
		http.Error(w, "deviceId and deviceType are required", http.StatusBadRequest)
		return
	}

	// First try to get state from MQTT cache
	state := mqtt.GetDeviceState(deviceID)
	if state != nil {
		json.NewEncoder(w).Encode(state)
		return
	}

	// If not in cache, fetch from YoLink API
	state, err := yolink.GetDeviceState(deviceID, deviceType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(state)
} 