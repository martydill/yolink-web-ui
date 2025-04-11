# YoLink Web UI

A modern web application for monitoring and managing YoLink IoT devices, built with Go and Svelte, and 99% AI-generated with Cursor! 

## Overview

This project was AI-generated and provides a user-friendly interface for monitoring YoLink IoT devices, including door sensors and vibration sensors. The application features real-time updates via MQTT and WebSocket connections, allowing users to track device states, battery levels, and online status.

## Features

- **Real-time Device Monitoring**
  - Live updates of device states via MQTT
  - WebSocket integration for instant UI updates
  - Battery level monitoring with visual indicators
  - Online/offline status tracking

- **Device Types Supported**
  - Door Sensors
  - Vibration Sensors
  - (Extensible for additional device types)

- **Modern UI/UX**
  - Clean, responsive design using Tailwind CSS
  - Interactive device cards with status indicators
  - Detailed device information panels
  - Real-time status updates with visual feedback

- **Technical Features**
  - Go backend with RESTful API
  - Svelte frontend for reactive UI
  - MQTT integration for real-time device updates
  - WebSocket support for live UI updates
  - Secure token-based authentication

## Project Structure

```
yolink-web-ui/
├── backend/           # Go backend server
│   ├── main.go       # Main application entry
│   ├── mqtt/         # MQTT client implementation
│   └── yolink/       # YoLink API integration
├── frontend/         # Svelte frontend
│   ├── src/
│   │   ├── App.svelte
│   │   └── DeviceDetails.svelte
│   └── package.json
└── README.md
```

## Getting Started

1. **Install Dependencies**
   ```bash
   make install
   ```

2. **Run the Application**
   ```bash
   make run
   ```

3. **Access the Dashboard**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080

## Development

- **Backend Development**
  ```bash
  make run-backend
  ```

- **Frontend Development**
  ```bash
  make run-frontend
  ```

## AI Generation

This project was generated using AI assistance, demonstrating the capabilities of modern AI in creating full-stack applications. The AI was used to:
- Design the application architecture
- Implement the backend API
- Create the frontend UI components
- Set up real-time communication
- Configure the build system

## License

This project is open source and available under the MIT License.

## Acknowledgments

- Built with Go and Svelte
- Uses Tailwind CSS for styling
- Integrates with YoLink IoT platform
- Powered by AI-assisted development 