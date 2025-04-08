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

<main class="container mx-auto px-4 py-8">
  <h1 class="text-3xl font-bold text-gray-800 mb-8">YoLink Devices</h1>

  {#if loading}
    <div class="flex justify-center items-center h-32">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>
  {:else if error}
    <div class="bg-red-50 border-l-4 border-red-500 p-4 mb-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-500" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <p class="text-sm text-red-700">{error}</p>
        </div>
      </div>
    </div>
  {:else}
    <div class="flex justify-end mb-4">
      <button 
        on:click={fetchDevices}
        class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 transition-colors duration-200 flex items-center space-x-2"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
        </svg>
        <span>Refresh Devices</span>
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each devices as device}
        <div 
          on:click={() => handleDeviceClick(device)}
          class="bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer p-6 border border-gray-100"
        >
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-semibold text-gray-800">{device.name}</h2>
            <span class="text-sm text-gray-500">{device.type}</span>
          </div>
          <div class="space-y-2">
            <p class="text-gray-600">
              <span class="font-medium">Model:</span> {device.modelName}
            </p>
            <p class="text-gray-600">
              <span class="font-medium">ID:</span> {device.deviceId}
            </p>
          </div>
        </div>
      {/each}
    </div>
  {/if}
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