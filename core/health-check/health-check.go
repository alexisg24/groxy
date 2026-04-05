package healthcheck

import (
	"fmt"
	"net/http"
	"time"

	fileloader "github.com/alexisg24/groxy/core/file-loader"
)

func PerformHealthCheck(hostKey string, hostConfig *fileloader.FileLoaderItem) {
	client := &http.Client{
		Timeout: time.Duration(hostConfig.HealthCheckTimeout) * time.Millisecond,
	}

	// Build the health check URL
	healthCheckURL := hostConfig.Path + hostConfig.HealthCheckPath
	// Create the health check request
	req, err := http.NewRequest("GET", healthCheckURL, nil)
	if err != nil {
		hostConfig.IsHealthy = false
		return
	}

	// Set custom headers for the health check request
	for header, value := range hostConfig.HealthCheckCustomHeaders {
		req.Header.Set(header, value)
	}
	// Perform the health check request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Health check for host %s: %v (error: %v)\n", hostKey, hostConfig.IsHealthy, err)

		hostConfig.IsHealthy = false
		return
	}
	defer resp.Body.Close()

	// Check if the response status code matches the expected status
	if resp.StatusCode == hostConfig.HealthCheckExpectedStatus {
		hostConfig.IsHealthy = true
	} else {
		hostConfig.IsHealthy = false
	}
	fmt.Printf("Health check for host %s: %v (status code: %d)\n", hostKey, hostConfig.IsHealthy, resp.StatusCode)
}

func InitializeHealthCheck(hostConfigs map[string]*fileloader.FileLoaderItem) {
	for key, item := range hostConfigs {
		if !item.HasHealthCheck {
			continue
		}

		// Add in a task every health check interval to perform the health check
		checkIntervals := time.NewTicker(time.Duration(item.HealthCheckInterval) * time.Millisecond)
		go func(hostKey string, hostConfig *fileloader.FileLoaderItem) {
			for range checkIntervals.C {
				PerformHealthCheck(hostKey, hostConfig)
			}
		}(key, item)
	}
}
