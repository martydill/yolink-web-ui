<script>
  import { onDestroy } from 'svelte';

  export let device;
  export let onClose;

  let deviceState = null;
  let loading = false;
  let error = null;
  let ws = null;
  let rawDataExpanded = false;

  function formatDeviceType(type) {
    if (!type || typeof type !== 'string') return 'Unknown Device';
    
    // Handle specific known mappings first
    const typeMap = {
      'smartremoter': 'Smart Remote',
      'leaksensor': 'Leak Sensor',
      'motionsensor': 'Motion Sensor',
      'doorwindowsensor': 'Door/Window Sensor',
      'doorsensor': 'Door Sensor',
      'windowsensor': 'Window Sensor',
      'vibrationsensor': 'Vibration Sensor',
      'temperaturesensor': 'Temperature Sensor',
      'watersensor': 'Water Sensor',
      'smokesensor': 'Smoke Sensor',
      'cosensor': 'CO Sensor',
      'co2sensor': 'CO2 Sensor',
      'humiditysensor': 'Humidity Sensor',
      'lightsensor': 'Light Sensor',
      'presssensor': 'Pressure Sensor',
      'pressuresensor': 'Pressure Sensor',
      'airqualitysensor': 'Air Quality Sensor',
      'smartlock': 'Smart Lock',
      'smartswitch': 'Smart Switch',
      'smartplug': 'Smart Plug',
      'smartbulb': 'Smart Bulb',
      'smartthermostat': 'Smart Thermostat',
      'smartcamera': 'Smart Camera',
      'smartdoorbell': 'Smart Doorbell',
      'smartgarage': 'Smart Garage',
      'smartvalve': 'Smart Valve'
    };
    
    const lowerType = type.toLowerCase();
    if (typeMap[lowerType]) {
      return typeMap[lowerType];
    }
    
    // Generic conversion for camelCase/PascalCase to readable format
    return type
      .replace(/([a-z])([A-Z])/g, '$1 $2') // Insert space before capital letters
      .replace(/([A-Z])([A-Z][a-z])/g, '$1 $2') // Handle consecutive capitals
      .split(' ')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(' ');
  }

  function getDeviceIcon(type) {
    if (!type || typeof type !== 'string') return '📱';
    switch (type.toLowerCase()) {
      case 'doorwindowsensor':
      case 'doorsensor':
      case 'door':
      case 'window':
      case 'windowsensor':
        return '🪟';
      case 'hub':
      case 'smarthub':
      case 'yolinkhub':
        return '🏠';
      case 'vibrationsensor':
      case 'vibration':
        return '📳';
      case 'motionsensor':
      case 'motion':
        return '🏃';
      case 'temperaturesensor':
      case 'temperature':
        return '🌡️';
      case 'leaksensor':
      case 'leak':
      case 'water':
      case 'watersensor':
        return '💧';
      case 'smartremoter':
      case 'smartremote':
      case 'remote':
        return '📱';
      case 'siren':
      case 'sirenspeaker':
      case 'alarm':
      case 'alarmsiren':
      case 'speaker':
        return '🚨';
      default:
        return '📱';
    }
  }

  function getStatusColor(status) {
    if (!status || typeof status !== 'string') return 'bg-gray-400';
    switch (status.toLowerCase()) {
      case 'online':
      case 'open':
      case 'active':
        return 'bg-emerald-400';
      case 'offline':
      case 'closed':
      case 'inactive':
        return 'bg-red-400';
      default:
        return 'bg-amber-400';
    }
  }

  function getStatusTextColor(status) {
    if (!status) return 'text-gray-600';
    
    // Convert to string and lowercase for comparison
    const statusStr = String(status).toLowerCase();
    
    // Safe/good states (green)
    if (statusStr.includes('closed') || 
        statusStr.includes('no leak') || 
        statusStr.includes('no vibration') || 
        statusStr.includes('no motion') ||
        statusStr.includes('online') ||
        statusStr.includes('active')) {
      return 'text-emerald-600';
    }
    
    // Alert/bad states (red)
    if (statusStr.includes('open') || 
        statusStr.includes('leak') || 
        statusStr.includes('vibration') || 
        statusStr.includes('motion detected') ||
        statusStr.includes('offline') ||
        statusStr.includes('inactive')) {
      return 'text-red-600';
    }
    
    // Temperature readings or other values (blue)
    if (statusStr.includes('°') || /^\d+/.test(statusStr)) {
      return 'text-blue-600';
    }
    
    // Unknown/neutral states (amber)
    return 'text-amber-600';
  }

  function getDeviceSpecificState(device, deviceState) {
    if (!device) return 'Unknown';
    
    const deviceType = device.type?.toLowerCase() || '';
    
    // Door/Window Sensors
    if (deviceType.includes('door') || deviceType.includes('window')) {
      // Check deviceState first, then device itself
      if (deviceState?.state === 'open' || deviceState?.open === true || device.state === 'open' || device.open === true) return 'Open';
      if (deviceState?.state === 'closed' || deviceState?.open === false || device.state === 'closed' || device.open === false) return 'Closed';
      return 'Closed'; // Default safe state
    }
    
    // Water Leak Detectors
    if (deviceType.includes('leak') || deviceType.includes('water')) {
      // Check deviceState first, then device itself
      if (deviceState?.leak === true || deviceState?.state === 'leak' || device.leak === true || device.state === 'leak') return 'Leak';
      if (deviceState?.leak === false || deviceState?.state === 'no leak' || device.leak === false || device.state === 'no leak') return 'No Leak';
      return 'No Leak'; // Default safe state
    }
    
    // Vibration Sensors
    if (deviceType.includes('vibration')) {
      if (deviceState?.vibration === true || deviceState?.state === 'vibration' || device.vibration === true || device.state === 'vibration') return 'Vibration';
      if (deviceState?.vibration === false || deviceState?.state === 'no vibration' || device.vibration === false || device.state === 'no vibration') return 'No Vibration';
      return 'No Vibration'; // Default safe state
    }
    
    // Motion Sensors
    if (deviceType.includes('motion')) {
      if (deviceState?.motion === true || deviceState?.state === 'motion' || device.motion === true || device.state === 'motion') return 'Motion Detected';
      if (deviceState?.motion === false || deviceState?.state === 'no motion' || device.motion === false || device.state === 'no motion') return 'No Motion';
      return 'No Motion'; // Default safe state
    }
    
    // Temperature Sensors
    if (deviceType.includes('temperature') || deviceType.includes('temp')) {
      if (deviceState?.temperature !== undefined) return `${deviceState.temperature}°`;
      if (device.temperature !== undefined) return `${device.temperature}°`;
      return 'Normal';
    }
    
    // Default fallback - check for any readable state
    if (deviceState?.state && typeof deviceState.state === 'string') {
      return deviceState.state;
    }
    
    if (device.state && typeof device.state === 'string') {
      return device.state;
    }
    
    // If state is an object, try to extract meaningful info
    if (typeof deviceState?.state === 'object' && deviceState.state !== null) {
      // Try to find common state properties
      const stateObj = deviceState.state;
      if (stateObj.open !== undefined) return stateObj.open ? 'Open' : 'Closed';
      if (stateObj.leak !== undefined) return stateObj.leak ? 'Leak' : 'No Leak';
      if (stateObj.motion !== undefined) return stateObj.motion ? 'Motion' : 'No Motion';
      if (stateObj.vibration !== undefined) return stateObj.vibration ? 'Vibration' : 'No Vibration';
      
      return 'Complex State';
    }
    
    return 'Unknown';
  }

  function formatDateTime(timestamp) {
    if (!timestamp) return 'Unknown';
    return new Date(timestamp).toLocaleString();
  }

  function getBatteryColor(level) {
    if (level > 70) return 'text-emerald-600';
    if (level > 30) return 'text-amber-600';
    return 'text-red-600';
  }

  function getBatteryIcon(level) {
    if (level >= 80) return '🔋'; // Full battery
    if (level >= 60) return '🔋'; // High battery  
    if (level >= 40) return '🔋'; // Medium battery
    if (level >= 20) return '🪫'; // Low battery
    return '🪫'; // Critical battery
  }

  function getBatteryLevelIcon(level) {
    // Returns icon based on 0-4 scale from battery field
    if (level === 4) return '🔋'; // Full
    if (level === 3) return '🔋'; // High
    if (level === 2) return '🔋'; // Medium
    if (level === 1) return '🪫'; // Low
    if (level === 0) return '🪫'; // Critical/Empty
    return '🔋'; // Default
  }

  function getBatteryLevelColor(level) {
    if (level >= 3) return 'text-emerald-600';
    if (level >= 2) return 'text-yellow-600';
    if (level >= 1) return 'text-orange-600';
    return 'text-red-600';
  }

  async function fetchDeviceState() {
    loading = true;
    error = null;
    try {
      const response = await fetch(`/api/devices/state?deviceId=${device.deviceId}&deviceType=${device.type}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      deviceState = await response.json();
    } catch (err) {
      error = err.message;
      console.error('Error:', err);
    }
    loading = false;
  }

  function setupWebSocket() {
    if (ws) {
      ws.close();
    }

    // In Docker setup, we need to connect to the backend through the frontend port proxy
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    let wsUrl;
    
    // Check if we're in Docker (frontend on port 3000) or development
    if (window.location.port === '3000') {
      // Docker setup - use the nginx proxy
      wsUrl = `${protocol}//${window.location.host}/ws`;
    } else {
      // Development setup - connect directly to backend
      wsUrl = `${protocol}//${window.location.hostname}:8080/ws`;
    }
    
    console.log('Connecting to WebSocket at:', wsUrl);
    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      console.log('WebSocket connection established');
    };

    ws.onmessage = (event) => {
      try {
        const update = JSON.parse(event.data);
        if (update.deviceId === device.deviceId) {
          deviceState = { ...deviceState, ...update };
        }
      } catch (err) {
        console.error('WebSocket message parse error:', err);
      }
    };

    ws.onerror = (error) => {
      console.warn('WebSocket connection failed - this is normal if backend is not configured with real credentials:', error);
    };

    ws.onclose = (event) => {
      console.log('WebSocket connection closed:', event.code, event.reason);
      // Don't automatically reconnect to avoid spam
    };
  }

  function handleKeydown(event) {
    if (event.key === 'Escape') {
      event.preventDefault();
      onClose();
    }
  }

  function handleBackdropClick(event) {
    // Only close if clicking on the backdrop itself, not on the modal content
    if (event.target === event.currentTarget) {
      onClose();
    }
  }

  function handleClose() {
    onClose();
  }

  // Fetch device state when component mounts or device changes
  $: if (device) {
    fetchDeviceState();
    setupWebSocket();
  }

  onDestroy(() => {
    if (ws) {
      ws.close();
    }
  });
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal Backdrop -->
<div 
  class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
  on:click={handleBackdropClick}
  on:keydown={handleKeydown}
  role="dialog"
  aria-modal="true"
  aria-labelledby="modal-title"
>
  <!-- Modal Content -->
  <div 
    class="bg-white/98 backdrop-blur-xl rounded-3xl shadow-2xl w-full max-w-3xl max-h-[92vh] overflow-hidden border border-white/20 transform transition-all duration-300 scale-100 ring-1 ring-black/5"
    on:click={(e) => e.stopPropagation()}
  >
    <!-- Header -->
    <div class="bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 px-6 py-5 text-white relative overflow-hidden">
      <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,_var(--tw-gradient-stops))] from-purple-400/30 via-transparent to-transparent"></div>
      <div class="absolute inset-0 bg-grid-white/[0.05] bg-[size:20px_20px]"></div>
      <div class="relative flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <div class="text-4xl drop-shadow-lg">{getDeviceIcon(device.type)}</div>
          <div>
            <h2 id="modal-title" class="text-2xl font-bold">{device.name}</h2>
            <p class="text-white/70 text-base">{formatDeviceType(device.type)}</p>
            <div class="flex items-center space-x-2 mt-1">
              <div class="flex items-center space-x-1">
                <div class="w-2 h-2 rounded-full {getStatusColor(deviceState?.online ? 'online' : 'offline')} ring-1 ring-white/30"></div>
                <span class="text-sm text-white/80">
                  {deviceState?.online ? 'Online' : 'Offline'}
                </span>
              </div>
            </div>
          </div>
        </div>
        <button
          on:click={handleClose}
          class="p-2 hover:bg-white/20 rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-white/50"
          aria-label="Close modal"
          type="button"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="p-4 max-h-[calc(92vh-140px)] overflow-y-auto bg-gradient-to-b from-gray-50/50 to-white">
      {#if loading}
        <div class="flex flex-col items-center justify-center py-8">
          <div class="relative">
            <div class="w-12 h-12 border-4 border-purple-200 border-t-purple-600 rounded-full animate-spin"></div>
            <div class="absolute inset-0 flex items-center justify-center">
              <div class="w-6 h-6 bg-gradient-to-r from-purple-500 to-indigo-600 rounded-full animate-pulse"></div>
            </div>
          </div>
          <p class="mt-4 text-gray-600">Loading device state...</p>
        </div>
      {:else if error}
        <div class="bg-gradient-to-br from-red-50 to-orange-50 border border-red-200/60 rounded-xl p-4 shadow-lg">
          <div class="flex items-start space-x-3">
            <div class="w-8 h-8 bg-gradient-to-r from-red-500 to-red-600 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <div class="flex-1">
              <h3 class="text-lg font-semibold text-red-800 mb-2">Error Loading Device State</h3>
              <p class="text-red-700 mb-3 text-sm">{error}</p>
              <button
                on:click={fetchDeviceState}
                class="bg-gradient-to-r from-red-500 to-red-600 text-white px-4 py-2 rounded-lg hover:from-red-600 hover:to-red-700 transition-all duration-200 text-sm font-medium"
              >
                Try Again
              </button>
            </div>
          </div>
        </div>
      {:else}
        <!-- Device Information Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6">
          <!-- Basic Info Card -->
          <div class="bg-gradient-to-br from-blue-50/80 via-indigo-50/60 to-purple-50/40 rounded-2xl p-4 border border-blue-200/30 shadow-lg backdrop-blur-sm">
            <h3 class="text-lg font-semibold text-gray-800 mb-3 flex items-center">
              <div class="w-6 h-6 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-lg flex items-center justify-center mr-2">
                <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              Device Information
            </h3>
            <div class="space-y-1">
              <div class="flex justify-between items-center py-1 border-b border-gray-200/50">
                <span class="text-gray-600 text-sm">Model</span>
                <span class="font-medium text-gray-900 bg-white/60 px-2 py-0.5 rounded text-sm">{device.modelName}</span>
              </div>
              <div class="flex justify-between items-center py-1 border-b border-gray-200/50">
                <span class="text-gray-600 text-sm">Device ID</span>
                <span class="font-medium text-gray-900 bg-white/60 px-2 py-0.5 rounded text-sm">{device.deviceId}</span>
              </div>
              <div class="flex justify-between items-center py-1">
                <span class="text-gray-600 text-sm">Type</span>
                <span class="font-medium text-gray-900 bg-white/60 px-2 py-0.5 rounded text-sm">{formatDeviceType(device.type)}</span>
              </div>
            </div>
          </div>

          <!-- Status Card -->
          <div class="bg-gradient-to-br from-emerald-50/80 via-green-50/60 to-teal-50/40 rounded-2xl p-4 border border-emerald-200/30 shadow-lg backdrop-blur-sm">
            <h3 class="text-lg font-semibold text-gray-800 mb-3 flex items-center">
              <div class="w-6 h-6 bg-gradient-to-r from-emerald-500 to-green-600 rounded-lg flex items-center justify-center mr-2">
                <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              Status & Power
            </h3>
            <div class="space-y-1">
              <div class="flex justify-between items-center py-1 border-b border-gray-200/50">
                <span class="text-gray-600 text-sm">Connection</span>
                <div class="flex items-center space-x-2 bg-white/60 px-2 py-0.5 rounded">
                  <div class="w-2 h-2 rounded-full {getStatusColor(deviceState?.online ? 'online' : 'offline')} ring-1 ring-white"></div>
                  <span class="font-medium {getStatusTextColor(deviceState?.online ? 'online' : 'offline')} text-sm">
                    {deviceState?.online ? 'Online' : 'Offline'}
                  </span>
                </div>
              </div>
              {#if deviceState}
                <div class="flex justify-between items-center py-1 border-b border-gray-200/50">
                  <span class="text-gray-600 text-sm">State</span>
                  <span class="font-medium {getStatusTextColor(String(getDeviceSpecificState(device, deviceState)).toLowerCase())} bg-white/60 px-2 py-0.5 rounded text-sm">
                    {getDeviceSpecificState(device, deviceState)}
                  </span>
                </div>
              {/if}
              {#if deviceState && (deviceState.state?.battery !== undefined || deviceState.battery !== undefined || deviceState.batteryLevel !== undefined || deviceState.power !== undefined)}
                {@const batteryLevel = deviceState.state?.battery || deviceState.battery || deviceState.batteryLevel || deviceState.power || 0}
                {@const batteryPercentage = Math.round((batteryLevel / 4) * 100)}
                <div class="flex justify-between items-center py-1">
                  <span class="text-gray-600 text-sm">Battery</span>
                  <div class="flex items-center space-x-2 bg-white/60 px-2 py-0.5 rounded">
                    <span class="text-lg">{getBatteryLevelIcon(batteryLevel)}</span>
                    <div class="flex flex-col items-end">
                      <span class="font-bold {getBatteryLevelColor(batteryLevel)} text-sm">
                        {batteryPercentage}%
                      </span>
                    </div>
                  </div>
                </div>
              {/if}
            </div>
          </div>
        </div>

        <!-- Detailed State Information -->
        {#if deviceState}
          <div class="bg-gradient-to-br from-slate-50/80 via-gray-50/60 to-zinc-50/40 rounded-2xl p-4 border border-slate-200/30 shadow-lg backdrop-blur-sm">
            <button 
              on:click={() => rawDataExpanded = !rawDataExpanded}
              class="w-full flex items-center justify-between text-lg font-semibold text-gray-800 mb-3 hover:text-gray-600 transition-colors duration-200"
            >
              <div class="flex items-center">
                <div class="w-6 h-6 bg-gradient-to-r from-slate-600 to-gray-700 rounded-lg flex items-center justify-center mr-2">
                  <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                  </svg>
                </div>
                Raw Device Data
              </div>
              <svg class="w-5 h-5 transform transition-transform duration-200 {rawDataExpanded ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
              </svg>
            </button>
            {#if rawDataExpanded}
              <div class="bg-white/80 backdrop-blur-sm rounded-xl p-4 border border-white/60 shadow-sm">
                <div class="bg-slate-900 text-green-400 rounded-lg p-4 max-w-full overflow-x-auto">
                  <pre class="text-xs font-mono whitespace-pre-wrap">
                    {JSON.stringify(
                      Object.fromEntries(
                        Object.entries(deviceState).filter(([key]) => 
                          !['deviceId', 'online', 'reportAt'].includes(key)
                        )
                      ), 
                      null, 
                      2
                    )}
                  </pre>
                </div>
              </div>
            {/if}
          </div>
        {/if}

        <!-- Real-time Updates Indicator -->
        <div class="mt-4 flex items-center justify-center">
          <div class="flex items-center space-x-2 bg-gradient-to-r from-emerald-50 to-teal-50 px-3 py-1 rounded-full border border-emerald-200/50">
            <div class="w-2 h-2 bg-gradient-to-r from-emerald-400 to-teal-500 rounded-full animate-pulse"></div>
            <span class="text-xs font-medium text-emerald-700">Real-time updates enabled</span>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  /* Smooth entrance animation */
  @keyframes modalEnter {
    from {
      opacity: 0;
      transform: scale(0.95) translateY(20px);
    }
    to {
      opacity: 1;
      transform: scale(1) translateY(0);
    }
  }

  .modal-content {
    animation: modalEnter 0.3s ease-out;
  }

  /* Custom scrollbar for the modal content */
  .overflow-y-auto::-webkit-scrollbar {
    width: 6px;
  }

  .overflow-y-auto::-webkit-scrollbar-track {
    background: rgba(0, 0, 0, 0.1);
    border-radius: 3px;
  }

  .overflow-y-auto::-webkit-scrollbar-thumb {
    background: linear-gradient(45deg, #6366f1, #8b5cf6);
    border-radius: 3px;
  }

  .overflow-y-auto::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(45deg, #4f46e5, #7c3aed);
  }
</style>