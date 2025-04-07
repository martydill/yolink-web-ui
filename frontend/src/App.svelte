<script>
  import './app.css';
  import DeviceDetails from './DeviceDetails.svelte';
  
  let devices = [];
  let loading = false;
  let error = null;
  let selectedDevice = null;

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

  // Fetch devices when component mounts
  fetchDevices();
</script>

<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <h1 class="text-3xl font-bold text-center text-indigo-600 mb-8">YoLink Devices</h1>
  
  <div class="bg-white rounded-lg shadow-md p-6">
    <button 
      on:click={fetchDevices} 
      disabled={loading}
      class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 disabled:bg-gray-400 disabled:cursor-not-allowed mb-4"
    >
      {loading ? 'Loading...' : 'Refresh Devices'}
    </button>
    
    {#if error}
      <p class="text-red-500 mb-4">{error}</p>
    {/if}

    {#if devices.length > 0}
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {#each devices as device}
          <div 
            on:click={() => handleDeviceClick(device)}
            class="bg-white rounded-lg shadow-md p-4 hover:shadow-lg transition-shadow cursor-pointer"
          >
            <h3 class="text-xl font-semibold text-gray-800 mb-2">{device.name}</h3>
            <p class="text-gray-600">Type: {device.type}</p>
            <p class="text-gray-600">ID: {device.deviceId}</p>
            <p class="text-gray-600">Model: {device.modelName}</p>
          </div>
        {/each}
      </div>
    {:else if !loading}
      <p class="text-gray-600 text-center">No devices found</p>
    {/if}
  </div>
</main>

{#if selectedDevice}
  <DeviceDetails 
    device={selectedDevice} 
    onClose={handleCloseDetails}
  />
{/if}

<style>
  main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }

  h1 {
    color: #646cff;
    margin-bottom: 2rem;
    text-align: center;
  }

  .card {
    padding: 2em;
    border-radius: 8px;
    background-color: #f9f9f9;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  button {
    border-radius: 8px;
    border: 1px solid transparent;
    padding: 0.6em 1.2em;
    font-size: 1em;
    font-weight: 500;
    font-family: inherit;
    background-color: #646cff;
    color: white;
    cursor: pointer;
    transition: background-color 0.25s;
    margin-bottom: 1rem;
  }

  button:hover {
    background-color: #535bf2;
  }

  button:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }

  .error {
    color: #ff4444;
    margin: 1rem 0;
  }

  .devices-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 1rem;
    margin-top: 1rem;
  }

  .device-card {
    padding: 1rem;
    border-radius: 8px;
    background-color: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s;
  }

  .device-card:hover {
    transform: translateY(-2px);
  }

  .device-card h3 {
    margin: 0 0 0.5rem 0;
    color: #333;
  }

  .device-card p {
    margin: 0.25rem 0;
    color: #666;
  }
</style> 