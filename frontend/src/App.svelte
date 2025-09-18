<script>
  import './app.css';
  import DeviceDetails from './DeviceDetails.svelte';
  
  let devices = [];
  let loading = false;
  let error = null;
  let selectedDevice = null;
  let searchTerm = '';
  let selectedDeviceType = 'all';
  let viewMode = 'grid'; // 'grid' or 'list'

  async function fetchDevices() {
    loading = true;
    error = null;
    try {
      const response = await fetch('/api/devices');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      devices = await response.json();
    } catch (err) {
      error = err.message;
      console.error('Error:', err);
    }
    loading = false;
  }

  function handleDeviceClick(device) {
    selectedDevice = device;
  }

  function handleCloseDetails() {
    selectedDevice = null;
  }

  $: filteredDevices = devices.filter(device => {
    const searchLower = searchTerm.toLowerCase();
    const matchesSearch = (
      (device.name || '').toLowerCase().includes(searchLower) ||
      (device.type || '').toLowerCase().includes(searchLower) ||
      (device.modelName || '').toLowerCase().includes(searchLower)
    );
    
    const matchesType = selectedDeviceType === 'all' || device.type === selectedDeviceType;
    
    return matchesSearch && matchesType;
  });

  $: uniqueDeviceTypes = [...new Set(devices.map(d => d.type))].sort();

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

  function getDeviceStatus(device) {
    // Check multiple possible status fields
    if (device.online === true) return 'online';
    if (device.online === false) return 'offline';
    if (device.isOnline === true) return 'online';
    if (device.isOnline === false) return 'offline';
    if (typeof device.status === 'string') {
      const status = device.status.toLowerCase();
      if (status === 'online' || status === 'active' || status === 'connected') return 'online';
      if (status === 'offline' || status === 'inactive' || status === 'disconnected') return 'offline';
    }
    // For debugging - assume all devices are online for now
    return 'online';
  }

  function getDeviceCurrentState(device) {
    if (!device) return 'Unknown';
    
    const deviceType = device.type?.toLowerCase() || '';
    
    // Door/Window Sensors
    if (deviceType.includes('door') || deviceType.includes('window')) {
      if (device.state === 'open' || device.open === true) return 'Open';
      if (device.state === 'closed' || device.open === false) return 'Closed';
      // Fallback to online/offline if no specific state
      return getDeviceStatus(device) === 'online' ? 'Closed' : 'Offline';
    }
    
    // Water Leak Detectors
    if (deviceType.includes('leak') || deviceType.includes('water')) {
      if (device.leak === true || device.state === 'leak') return 'Leak';
      if (device.leak === false || device.state === 'no leak') return 'No Leak';
      return getDeviceStatus(device) === 'online' ? 'No Leak' : 'Offline';
    }
    
    // Vibration Sensors
    if (deviceType.includes('vibration')) {
      if (device.vibration === true || device.state === 'vibration') return 'Vibration';
      if (device.vibration === false || device.state === 'no vibration') return 'No Vibration';
      return getDeviceStatus(device) === 'online' ? 'No Vibration' : 'Offline';
    }
    
    // Motion Sensors
    if (deviceType.includes('motion')) {
      if (device.motion === true || device.state === 'motion') return 'Motion';
      if (device.motion === false || device.state === 'no motion') return 'No Motion';
      return getDeviceStatus(device) === 'online' ? 'No Motion' : 'Offline';
    }
    
    // Temperature Sensors
    if (deviceType.includes('temperature') || deviceType.includes('temp')) {
      if (device.temperature !== undefined) return `${device.temperature}°`;
      return getDeviceStatus(device) === 'online' ? 'Normal' : 'Offline';
    }
    
    // Default fallback to online/offline
    return getDeviceStatus(device) === 'online' ? 'Online' : 'Offline';
  }

  function getStateColor(device) {
    if (!device) return 'text-gray-600';
    
    const state = getDeviceCurrentState(device);
    const deviceType = device.type?.toLowerCase() || '';
    
    // Door/Window sensors - Open is typically an alert state
    if (deviceType.includes('door') || deviceType.includes('window')) {
      return state === 'Open' ? 'text-red-600' : state === 'Closed' ? 'text-emerald-600' : 'text-gray-600';
    }
    
    // Leak sensors - Leak is an alert state
    if (deviceType.includes('leak') || deviceType.includes('water')) {
      return state === 'Leak' ? 'text-red-600' : state === 'No Leak' ? 'text-emerald-600' : 'text-gray-600';
    }
    
    // Vibration sensors - Vibration is an alert state
    if (deviceType.includes('vibration')) {
      return state === 'Vibration' ? 'text-red-600' : state === 'No Vibration' ? 'text-emerald-600' : 'text-gray-600';
    }
    
    // Motion sensors - Motion detected is an alert state
    if (deviceType.includes('motion')) {
      return state === 'Motion' ? 'text-red-600' : state === 'No Motion' ? 'text-emerald-600' : 'text-gray-600';
    }
    
    // Temperature sensors - show in blue
    if (deviceType.includes('temperature') || deviceType.includes('temp')) {
      return 'text-blue-600';
    }
    
    // Default online/offline colors
    return state === 'Online' ? 'text-emerald-600' : state === 'Offline' ? 'text-red-600' : 'text-gray-600';
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

  // Fetch devices when component mounts
  fetchDevices();
</script>

<!-- Background with gradient -->
<div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-cyan-50">
  <!-- Header Section -->
  <div class="bg-white/80 backdrop-blur-md border-b border-gray-200/50 sticky top-0 z-40">
    <div class="container mx-auto px-6 py-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <div class="w-10 h-10 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
            </svg>
          </div>
          <div>
            <h1 class="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
              YoLink Dashboard
            </h1>
          </div>
        </div>
        <div class="flex flex-col sm:flex-row items-start sm:items-center space-y-3 sm:space-y-0 sm:space-x-3 w-full sm:w-auto">
          <div class="flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-3 w-full sm:w-auto">
            <div class="relative w-full sm:w-auto">
              <span class="absolute inset-y-0 left-0 flex items-center pl-3">
                <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                </svg>
              </span>
              <input
                type="text"
                bind:value={searchTerm}
                placeholder="Search devices..."
                class="w-full sm:w-auto pl-10 pr-4 py-2 border border-gray-200 rounded-xl bg-white/50 backdrop-blur-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200"
              />
            </div>
            <div class="relative w-full sm:w-auto">
              <select
                bind:value={selectedDeviceType}
                class="w-full sm:w-auto appearance-none bg-white/50 backdrop-blur-sm border border-gray-200 rounded-xl pl-4 pr-10 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200"
              >
                <option value="all">All Types</option>
                {#each uniqueDeviceTypes as deviceType}
                  <option value={deviceType}>{formatDeviceType(deviceType)}</option>
                {/each}
              </select>
              <div class="absolute inset-y-0 right-0 flex items-center px-3 pointer-events-none">
                <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
                </svg>
              </div>
            </div>
          </div>
          <div class="flex items-center justify-between w-full sm:w-auto space-x-3">
            <div class="flex bg-gray-100 rounded-lg p-1">
              <button
                on:click={() => viewMode = 'grid'}
                class="p-2 rounded-md transition-all duration-200 {viewMode === 'grid' ? 'bg-white shadow-sm text-indigo-600' : 'text-gray-500 hover:text-gray-700'}"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/>
                </svg>
              </button>
              <button
                on:click={() => viewMode = 'list'}
                class="p-2 rounded-md transition-all duration-200 {viewMode === 'list' ? 'bg-white shadow-sm text-indigo-600' : 'text-gray-500 hover:text-gray-700'}"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16"/>
                </svg>
              </button>
            </div>
            <button 
              on:click={fetchDevices}
              class="bg-gradient-to-r from-indigo-500 to-purple-600 text-white px-4 sm:px-6 py-2 rounded-xl hover:from-indigo-600 hover:to-purple-700 transition-all duration-200 flex items-center space-x-2 shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 sm:h-5 sm:w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
              </svg>
              <span class="hidden sm:inline">Refresh</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Main Content -->
  <main class="container mx-auto px-6 py-8">
    {#if loading}
      <div class="flex flex-col justify-center items-center h-64">
        <div class="relative">
          <div class="w-16 h-16 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
          <div class="absolute inset-0 flex items-center justify-center">
            <div class="w-8 h-8 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-full animate-pulse"></div>
          </div>
        </div>
        <p class="mt-4 text-gray-600 animate-pulse">Loading your devices...</p>
      </div>
    {:else if error}
      <div class="max-w-md mx-auto">
        <div class="bg-gradient-to-r from-red-50 to-red-100 border border-red-200 rounded-2xl p-6 shadow-lg">
          <div class="flex items-start space-x-4">
            <div class="flex-shrink-0">
              <div class="w-10 h-10 bg-red-500 rounded-full flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-red-800 mb-2">Connection Error</h3>
              <p class="text-red-700">{error}</p>
              <button
                on:click={fetchDevices}
                class="mt-4 bg-red-500 text-white px-4 py-2 rounded-lg hover:bg-red-600 transition-colors duration-200"
              >
                Try Again
              </button>
            </div>
          </div>
        </div>
      </div>
    {:else if filteredDevices.length === 0 && devices.length > 0}
      <div class="text-center py-12">
        <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-700 mb-2">No devices found</h3>
        <p class="text-gray-500">Try adjusting your search terms</p>
      </div>
    {:else if devices.length === 0}
      <div class="text-center py-12">
        <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-700 mb-2">No devices connected</h3>
        <p class="text-gray-500">Connect your YoLink devices to get started</p>
      </div>
    {:else}
      <!-- Device Stats -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white/60 backdrop-blur-sm rounded-2xl p-6 border border-gray-200/50 shadow-lg">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Total Devices</p>
              <p class="text-2xl font-bold text-gray-900">{devices.length}</p>
            </div>
            <div class="w-12 h-12 bg-gradient-to-r from-blue-400 to-blue-600 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
              </svg>
            </div>
          </div>
        </div>
        <div class="bg-white/60 backdrop-blur-sm rounded-2xl p-6 border border-gray-200/50 shadow-lg">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Online</p>
              <p class="text-2xl font-bold text-emerald-600">{devices.filter(d => getDeviceStatus(d) === 'online').length}</p>
            </div>
            <div class="w-12 h-12 bg-gradient-to-r from-emerald-400 to-emerald-600 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
            </div>
          </div>
        </div>
        <div class="bg-white/60 backdrop-blur-sm rounded-2xl p-6 border border-gray-200/50 shadow-lg">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Offline</p>
              <p class="text-2xl font-bold text-red-600">{devices.filter(d => getDeviceStatus(d) === 'offline').length}</p>
            </div>
            <div class="w-12 h-12 bg-gradient-to-r from-red-400 to-red-600 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </div>
          </div>
        </div>
        <div class="bg-white/60 backdrop-blur-sm rounded-2xl p-6 border border-gray-200/50 shadow-lg">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Device Types</p>
              <p class="text-2xl font-bold text-indigo-600">{new Set(devices.map(d => d.type)).size}</p>
            </div>
            <div class="w-12 h-12 bg-gradient-to-r from-indigo-400 to-indigo-600 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"/>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Device Grid/List -->
      {#if viewMode === 'grid'}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {#each filteredDevices as device, index}
            <div 
              role="button"
              tabindex="0"
              on:click={() => handleDeviceClick(device)}
              on:keydown={(e) => e.key === 'Enter' && handleDeviceClick(device)}
              class="group bg-white/60 backdrop-blur-sm rounded-2xl p-6 border border-gray-200/50 shadow-lg hover:shadow-xl transition-all duration-300 cursor-pointer transform hover:scale-105 hover:bg-white/80"
              style="animation-delay: {index * 100}ms"
            >
              <div class="flex items-start justify-between mb-4">
                <div class="flex items-center space-x-2 min-w-0 flex-1">
                  <div class="text-2xl sm:text-3xl flex-shrink-0">{getDeviceIcon(device.type)}</div>
                  <div class="flex-1 min-w-0">
                    <h3 class="font-semibold text-gray-900 truncate group-hover:text-indigo-600 transition-colors duration-200 text-sm sm:text-base">
                      {device.name}
                    </h3>
                    <p class="text-xs sm:text-sm text-gray-500 truncate">{formatDeviceType(device.type)}</p>
                  </div>
                </div>
                <div class="flex items-center space-x-2 flex-shrink-0">
                  <div class="w-2 h-2 sm:w-3 sm:h-3 rounded-full {getStatusColor(getDeviceStatus(device))} ring-1 sm:ring-2 ring-white shadow-lg"></div>
                </div>
              </div>
              
              <div class="space-y-2 sm:space-y-3">
                <div class="flex items-center justify-between">
                  <span class="text-xs sm:text-sm text-gray-600 flex-shrink-0">Model</span>
                  <span class="text-xs sm:text-sm font-medium text-gray-900 truncate ml-2 min-w-0">{device.modelName}</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-xs sm:text-sm text-gray-600 flex-shrink-0">State</span>
                  <span class="text-xs sm:text-sm font-medium {getStateColor(device)} truncate min-w-0">
                    {getDeviceCurrentState(device)}
                  </span>
                </div>
              </div>
              
              <div class="mt-3 sm:mt-4 pt-3 sm:pt-4 border-t border-gray-200/50">
                <div class="flex items-center justify-between text-xs text-gray-500">
                  <span class="truncate min-w-0 flex-1">ID: {device.deviceId.slice(-8)}</span>
                  <svg class="w-3 h-3 sm:w-4 sm:h-4 text-gray-400 group-hover:text-indigo-500 transition-colors duration-200 flex-shrink-0 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                  </svg>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {:else}
        <div class="bg-white/60 backdrop-blur-sm rounded-2xl border border-gray-200/50 shadow-lg overflow-hidden">
          <div class="divide-y divide-gray-200/50">
            {#each filteredDevices as device, index}
              <div 
                role="button"
                tabindex="0"
                on:click={() => handleDeviceClick(device)}
                on:keydown={(e) => e.key === 'Enter' && handleDeviceClick(device)}
                class="group p-6 hover:bg-white/80 transition-all duration-200 cursor-pointer"
                style="animation-delay: {index * 50}ms"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-3 sm:space-x-4 min-w-0 flex-1">
                    <div class="text-xl sm:text-2xl flex-shrink-0">{getDeviceIcon(device.type)}</div>
                    <div class="min-w-0 flex-1">
                      <h3 class="font-semibold text-gray-900 group-hover:text-indigo-600 transition-colors duration-200 truncate text-sm sm:text-base">
                        {device.name}
                      </h3>
                      <p class="text-xs sm:text-sm text-gray-500 truncate">{formatDeviceType(device.type)} • {device.modelName}</p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-2 sm:space-x-6 flex-shrink-0">
                    <div class="flex items-center space-x-1 sm:space-x-2">
                      <div class="w-2 h-2 sm:w-3 sm:h-3 rounded-full {getStatusColor(getDeviceStatus(device))} ring-1 sm:ring-2 ring-white shadow-lg"></div>
                      <span class="text-xs sm:text-sm font-medium {getStateColor(device)} hidden sm:inline truncate">
                        {getDeviceCurrentState(device)}
                      </span>
                    </div>
                    <div class="text-xs text-gray-500 hidden md:block">
                      ID: {device.deviceId.slice(-8)}
                    </div>
                    <svg class="w-4 h-4 sm:w-5 sm:h-5 text-gray-400 group-hover:text-indigo-500 transition-colors duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                    </svg>
                  </div>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  </main>
</div>

{#if selectedDevice}
  <DeviceDetails 
    device={selectedDevice} 
    onClose={handleCloseDetails}
  />
{/if}

<style>
  /* Custom animations for staggered entrance */
  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-fade-in-up {
    animation: fadeInUp 0.6s ease-out forwards;
  }

  /* Glassmorphism effect enhancement */
  .glass-effect {
    background: rgba(255, 255, 255, 0.25);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.18);
  }

  /* Smooth scroll behavior */
  html {
    scroll-behavior: smooth;
  }

  /* Custom scrollbar for webkit browsers */
  ::-webkit-scrollbar {
    width: 6px;
  }

  ::-webkit-scrollbar-track {
    background: rgba(0, 0, 0, 0.1);
    border-radius: 3px;
  }

  ::-webkit-scrollbar-thumb {
    background: linear-gradient(45deg, #6366f1, #8b5cf6);
    border-radius: 3px;
  }

  ::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(45deg, #4f46e5, #7c3aed);
  }

  /* Enhanced focus styles for accessibility */
  *:focus {
    outline: 2px solid #6366f1;
    outline-offset: 2px;
  }
</style> 