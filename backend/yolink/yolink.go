package yolink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	baseURL     = "https://api.yosmart.com/open/yolink/v2/api"
	authURL     = "https://api.yosmart.com/open/yolink/token"
)

type Device struct {
	DeviceID       string      `json:"deviceId"`
	DeviceUDID     string      `json:"deviceUDID"`
	DeviceName     string      `json:"name"`
	Token          string      `json:"token"`
	Type           string      `json:"type"`
	ParentDeviceID interface{} `json:"parentDeviceId"`
	ModelName      string      `json:"modelName"`
	ServiceZone    string      `json:"serviceZone"`
}

type DeviceListData struct {
	Devices []Device `json:"devices"`
}

type DeviceListResponse struct {
	Code    string        `json:"code"`
	Time    int64        `json:"time"`
	MsgID   int64        `json:"msgid"`
	Method  string       `json:"method"`
	Desc    string       `json:"desc"`
	Data    DeviceListData `json:"data"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Code        int    `json:"code"`
	Message     string `json:"message"`
}

type DeviceStateResponse struct {
	Code    string                 `json:"code"`
	Time    int64                 `json:"time"`
	MsgID   int64                 `json:"msgid"`
	Method  string                `json:"method"`
	Desc    string                `json:"desc"`
	Data    map[string]interface{} `json:"data"`
}

var (
	accessToken string
	tokenExpiry time.Time
)

func getAccessToken() (string, error) {
	// If we have a token and it's not expired, return it
	if accessToken != "" && time.Now().Before(tokenExpiry) {
		return accessToken, nil
	}

	clientID := os.Getenv("YOLINK_CLIENT_ID")
	clientSecret := os.Getenv("YOLINK_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("YOLINK_CLIENT_ID and YOLINK_CLIENT_SECRET environment variables must be set")
	}

	// Create URL with parameters
	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("client_id", clientID)
	params.Add("client_secret", clientSecret)
	authURLWithParams := authURL + "?" + params.Encode()

	// Debug: Print the request details
	log.Printf("Making auth request to: %s", authURLWithParams)

	req, err := http.NewRequest("POST", authURLWithParams, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Debug: Print the response body
	log.Printf("Auth Response: %s", string(body))

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse auth response: %v, body: %s", err, string(body))
	}

	// Check if authentication was successful
	if tokenResp.Code != 0 {
		return "", fmt.Errorf("authentication failed: %s (code: %d)", tokenResp.Message, tokenResp.Code)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("authentication failed: no access token received")
	}

	// Set the token and its expiry time (using 90% of the actual expiry time to be safe)
	accessToken = tokenResp.AccessToken
	tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn*9/10) * time.Second)
	
	log.Printf("New token obtained, expires at: %v", tokenExpiry)
	return accessToken, nil
}

func GetDevices() ([]Device, error) {
	token, err := getAccessToken()
	if err != nil {
		// Clear the access token on auth failure
		accessToken = ""
		return nil, fmt.Errorf("authentication required: %v", err)
	}

	// Create the request payload
	payload := map[string]interface{}{
		"method":    "Home.getDeviceList",
		"timestamp": time.Now().Unix(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Debug: Print request details
	log.Printf("Making devices request to: %s", baseURL)
	log.Printf("With payload: %s", string(payloadBytes))
	log.Printf("With token: %s", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Debug: Print response status and headers
	log.Printf("Devices Response status: %s", resp.Status)
	log.Printf("Devices Response headers: %v", resp.Header)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Debug: Print the response body
	log.Printf("Devices Response: %s", string(body))

	var response DeviceListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse devices response: %v, body: %s", err, string(body))
	}

	if response.Code != "000000" {
		// If we get an unauthorized error, clear the token
		if response.Code == "401" {
			accessToken = ""
		}
		return nil, fmt.Errorf("API error: %s (code: %s)", response.Desc, response.Code)
	}

	return response.Data.Devices, nil
}

func GetDeviceState(deviceID string, deviceType string) (map[string]interface{}, error) {
	// First get the access token for the Authorization header
	token, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %v", err)
	}

	// Get the device list to get the device's token
	devices, err := GetDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %v", err)
	}

	// Find the specific device
	var deviceToken string
	for _, device := range devices {
		if device.DeviceID == deviceID {
			deviceToken = device.Token
			break
		}
	}

	if deviceToken == "" {
		return nil, fmt.Errorf("device not found or device token not available")
	}

	// Create the request payload
	payload := map[string]interface{}{
		"method":       deviceType + ".getState",
		"timestamp":    time.Now().Unix(),
		"targetDevice": deviceID,
		"token":        deviceToken,  // Add the device token to the payload
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}

	// Use the access token in the Authorization header
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Debug: Print request details
	log.Printf("Making device state request to: %s", baseURL)
	log.Printf("With payload: %s", string(payloadBytes))
	log.Printf("With access token: %s", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Debug: Print response status and headers
	log.Printf("Device State Response status: %s", resp.Status)
	log.Printf("Device State Response headers: %v", resp.Header)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Debug: Print the response body
	log.Printf("Device State Response: %s", string(body))

	var response DeviceStateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse device state response: %v, body: %s", err, string(body))
	}

	if response.Code != "000000" {
		return nil, fmt.Errorf("API error: %s (code: %s)", response.Desc, response.Code)
	}

	return response.Data, nil
} 