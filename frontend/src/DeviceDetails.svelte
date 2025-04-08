<script>
  export let device;
  export let onClose;

  let deviceState = null;
  let loading = false;
  let error = null;

  $: deviceType = device.type.toLowerCase();
  $: isDoorSensor = deviceType === 'doorsensor' || deviceType === 'door';
  $: isVibrationSensor = deviceType === 'vibrationsensor' || deviceType === 'vibration';

  // Battery level colors
  const getBatteryColor = (level) => {
    if (level >= 3) return 'text-green-500';
    if (level >= 2) return 'text-yellow-500';
    if (level >= 1) return 'text-orange-500';
    return 'text-red-500';
  };

  // Get online status color and icon
  const getOnlineStatus = (isOnline) => {
    return {
      color: isOnline ? 'text-green-500' : 'text-red-500',
      icon: isOnline ? `
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0" />
        </svg>
      ` : `
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      `
    };
  };

  // Get door status color and icon
  const getDoorStatus = (state) => {
    if (state === 'open') {
      return {
        color: 'text-red-500',
        icon: `
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        `,
        text: 'Open'
      };
    } else if (state === 'closed') {
      return {
        color: 'text-green-500',
        icon: `
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        `,
        text: 'Closed'
      };
    } else {
      return {
        color: 'text-yellow-500',
        icon: `
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        `,
        text: 'Error'
      };
    }
  };

  // Get vibration status color and icon
  const getVibrationStatus = (state) => {
    if (state === 'vibration') {
      return {
        color: 'text-red-500',
        icon: `
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
          </svg>
        `,
        text: 'Vibration Detected'
      };
    } else if (state === 'normal') {
      return {
        color: 'text-green-500',
        icon: `
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        `,
        text: 'Normal'
      };
    } else {
      return {
        color: 'text-yellow-500',
        icon: `
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        `,
        text: 'Error'
      };
    }
  };

  // Battery level icon
  const getBatteryIcon = (level) => {
    // Calculate fill percentage, leaving 2 units padding on each side
    const fill = Math.min(Math.max(level, 0), 4) * 25;
    const fillWidth = (fill / 100) * 14; // 14 is the available width (18 - 4 for padding)
    return `
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <!-- Battery body -->
        <rect x="2" y="6" width="18" height="12" rx="2" stroke-width="2" />
        <!-- Battery tip -->
        <path d="M20 10h2v4h-2" stroke-width="2" />
        <!-- Battery fill -->
        <rect x="4" y="8" width="${fillWidth}" height="8" fill="currentColor" rx="1" />
      </svg>
    `;
  };

  async function fetchDeviceState() {
    loading = true;
    error = null;
    try {
      const response = await fetch(`/api/devices/state?deviceId=${device.deviceId}&deviceType=${device.type}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      deviceState = await response.json();
      console.log('Device State Response:', deviceState);
      console.log('Device Type:', deviceType);
      console.log('State Path:', deviceState?.state?.state);
    } catch (err) {
      error = err.message;
      console.error('Error:', err);
    }
    loading = false;
  }

  // Fetch device state when component mounts
  fetchDeviceState();
</script>

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
  <div class="bg-white rounded-lg shadow-xl max-w-2xl w-full p-6">
    <div class="flex justify-between items-start mb-4">
      <h2 class="text-2xl font-bold text-gray-800">{device.name}</h2>
      <button 
        on:click={onClose}
        class="text-gray-500 hover:text-gray-700"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="space-y-2">
        <h3 class="font-semibold text-gray-700">Device Information</h3>
        <p><span class="text-gray-600">Type:</span> {deviceType}</p>
        <p><span class="text-gray-600">Model:</span> {device.modelName}</p>
        <div class="flex items-center space-x-2">
          <span class="text-gray-600">ID:</span>
          <span class="font-mono text-sm">{device.deviceId}</span>
        </div>
        <p><span class="text-gray-600">UDID:</span> {device.deviceUDID}</p>
      </div>

      <div class="space-y-2">
        <h3 class="font-semibold text-gray-700">Technical Details</h3>
        <p><span class="text-gray-600">Service Zone:</span> {device.serviceZone || 'Not specified'}</p>
        <p><span class="text-gray-600">Parent Device:</span> {device.parentDeviceId || 'None'}</p>
      </div>
    </div>

    {#if loading}
      <div class="mt-6 text-center">
        <p class="text-gray-600">Loading device state...</p>
      </div>
    {:else if error}
      <div class="mt-6">
        <p class="text-red-500">Error loading device state: {error}</p>
      </div>
    {:else if deviceState}
      <div class="mt-6">
        <h3 class="font-semibold text-gray-700 mb-2">Device State</h3>
        <div class="bg-gray-50 rounded-lg p-4">
          {#if deviceState.online !== undefined}
            <div class="flex items-center space-x-2">
              <span class="text-gray-600 mr-2">Status:</span>
              <div class={getOnlineStatus(deviceState.online).color}>
                {@html getOnlineStatus(deviceState.online).icon}
              </div>
              <span class="ml-2">{deviceState.online ? 'Online' : 'Offline'}</span>
            </div>
          {/if}
          {#if isDoorSensor && deviceState.state?.state !== undefined}
            <div class="flex items-center space-x-2 mt-4">
              <span class="text-gray-600 mr-2">Door:</span>
              <div class={getDoorStatus(deviceState.state.state).color}>
                {@html getDoorStatus(deviceState.state.state).icon}
              </div>
              <span class="ml-2">{getDoorStatus(deviceState.state.state).text}</span>
            </div>
          {/if}
          {#if isVibrationSensor && deviceState.state?.state !== undefined}
            <div class="flex items-center space-x-2 mt-4">
              <span class="text-gray-600 mr-2">Vibration:</span>
              <div class={getVibrationStatus(deviceState.state.state).color}>
                {@html getVibrationStatus(deviceState.state.state).icon}
              </div>
              <span class="ml-2">{getVibrationStatus(deviceState.state.state).text}</span>
            </div>
          {/if}
          {#each Object.entries(deviceState) as [key, value]}
            <div class="flex items-center space-x-2">
              {#if key === 'state' && value.battery !== undefined}
                <div class="flex items-center">
                  <span class="text-gray-600 mr-2">Battery:</span>
                  <div class={getBatteryColor(value.battery)}>
                    {@html getBatteryIcon(value.battery)}
                  </div>
                  <span class="ml-2">{value.battery}/4</span>
                </div>
              {/if}
              {#if key !== 'state' && key !== 'online'}
                <p class="text-gray-600">
                  <span class="font-medium">{key}:</span> {JSON.stringify(value, null, 2)}
                </p>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <div class="mt-6 flex justify-end">
      <button 
        on:click={onClose}
        class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700"
      >
        Close
      </button>
    </div>
  </div>
</div> 