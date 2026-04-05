package fileloader

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	var globalFileLoaderConfig GlobalFileLoaderConfig
	err := globalFileLoaderConfig.Load("config.test.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(globalFileLoaderConfig.Hosts) == 0 {
		t.Fatalf("No hosts loaded from config")
	}

	expectedHost := "backend_url"
	hostConfig, exists := globalFileLoaderConfig.GetItem(expectedHost)
	if !exists {
		t.Fatalf("Expected host %q not found in config", expectedHost)
	}

	if hostConfig.Path != "http://localhost:8080" {
		t.Fatalf("Expected path %q for host %q, got %q", "http://localhost:8080", expectedHost, hostConfig.Path)
	}

	if hostConfig.Timeout != 5000 {
		t.Fatalf("Expected timeout 5000 for host %q, got %d", expectedHost, hostConfig.Timeout)
	}

	if hostConfig.RetryCount != 3 {
		t.Fatalf("Expected retry count 3 for host %q, got %d", expectedHost, hostConfig.RetryCount)
	}

	if hostConfig.HealthCheckPath != "/health" {
		t.Fatalf("Expected health check path %q for host %q, got %q", "/health", expectedHost, hostConfig.HealthCheckPath)
	}

	if hostConfig.HealthCheckInterval != 10000 {
		t.Fatalf("Expected health check interval 10000 for host %q, got %d", expectedHost, hostConfig.HealthCheckInterval)
	}

	if hostConfig.HealthCheckTimeout != 5000 {
		t.Fatalf("Expected health check timeout 5000 for host %q, got %d", expectedHost, hostConfig.HealthCheckTimeout)
	}

	if hostConfig.HealthCheckCustomHeaders["X-Custom-Header"] != "CustomValue" {
		t.Fatalf("Expected health check custom header X-Custom-Header to be 'CustomValue' for host %q, got %q", expectedHost, hostConfig.HealthCheckCustomHeaders["X-Custom-Header"])
	}

	if hostConfig.HealthCheckExpectedStatus != 200 {
		t.Fatalf("Expected health check expected status 200 for host %q, got %d", expectedHost, hostConfig.HealthCheckExpectedStatus)
	}

	if len(hostConfig.InternalHeaders) == 0 {
		t.Fatalf("Expected additional headers for host %q, got none", expectedHost)
	}

	if hostConfig.PayloadSizeLimit < 0 {
		t.Fatalf("Expected non-negative payload size limit for host %q, got %d", expectedHost, hostConfig.PayloadSizeLimit)
	}

}
