package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"yolink-web-ui/backend/yolink"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	client mqtt.Client
	mu     sync.RWMutex
	// Map to store device state updates
	deviceStates = make(map[string]interface{})
	// Channel to broadcast updates to WebSocket clients
	updateChannel = make(chan DeviceUpdate, 100)
)

type DeviceUpdate struct {
	DeviceID string      `json:"deviceId"`
	State    interface{} `json:"state"`
}

// Initialize MQTT client
func Init() error {
	// Get access token from YoLink API
	accessToken, err := yolink.GetAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %v", err)
	}

	// Get home ID
	homeID, err := yolink.GetHomeID()
	if err != nil {
		return fmt.Errorf("failed to get home ID: %v", err)
	}

	// YoLink MQTT broker configuration
	brokerHost := "api.yosmart.com"
	brokerPort := "8003"
	clientID := fmt.Sprintf("yolink-app-%d", time.Now().Unix())

	log.Printf("Connecting to YoLink MQTT broker at %s:%s", brokerHost, brokerPort)
	log.Printf("Using client ID: %s", clientID)
	log.Printf("Using home ID: %s", homeID)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", brokerHost, brokerPort))
	opts.SetClientID(clientID)
	opts.SetUsername(accessToken)
	opts.SetPassword("") // YoLink MQTT doesn't require password
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Connection to MQTT broker lost: %v", err)
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Printf("Successfully connected to MQTT broker")
	})

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	// Subscribe to device state topics using home ID
	topic := fmt.Sprintf("yl-home/%s/+/report", homeID)
	log.Printf("Subscribing to topic: %s", topic)
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %v", topic, token.Error())
	}

	log.Printf("Successfully subscribed to %s", topic)
	return nil
}

// Message handler for MQTT messages
func messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received MQTT message on topic: %s", msg.Topic())
	log.Printf("Message payload: %s", string(msg.Payload()))

	deviceID := extractDeviceID(msg.Topic())
	if deviceID == "" {
		log.Printf("Could not extract device ID from topic: %s", msg.Topic())
		return
	}

	// Parse the message payload
	state := parseStatePayload(msg.Payload())
	if state == nil {
		log.Printf("Could not parse state payload for device %s", deviceID)
		return
	}

	// Update device state
	mu.Lock()
	deviceStates[deviceID] = state
	mu.Unlock()

	log.Printf("Updated state for device %s", deviceID)

	// Broadcast update to WebSocket clients
	updateChannel <- DeviceUpdate{
		DeviceID: deviceID,
		State:    state,
	}
}

// Get the latest state for a device
func GetDeviceState(deviceID string) interface{} {
	mu.RLock()
	defer mu.RUnlock()
	return deviceStates[deviceID]
}

// Get the update channel for WebSocket clients
func GetUpdateChannel() <-chan DeviceUpdate {
	return updateChannel
}

// Helper function to extract device ID from topic
func extractDeviceID(topic string) string {
	// Topic format: yl-home/{homeId}/{deviceId}/report
	parts := strings.Split(topic, "/")
	if len(parts) != 4 {
		return ""
	}
	return parts[2]
}

// Helper function to parse state payload
func parseStatePayload(payload []byte) interface{} {
	var state interface{}
	if err := json.Unmarshal(payload, &state); err != nil {
		log.Printf("Error parsing state payload: %v", err)
		return nil
	}
	return state
}
