package http_handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexisg24/groxy/core/config"
)

func TestHandleProxyRequest_Success(t *testing.T) {
	err := config.GlobalConfig.Load("../file-loader/config.test.yaml")
	if err != nil {
		panic(err)
	}
	hostConfig, exists := config.GlobalConfig.GetItem("backend_url")
	if !exists {
		t.Fatalf("expected hostConfig to exist")
	}

	expectedResponseBody := "Test request body"
	expectedHeaderValue := "test-value"
	headerKey := "X-Test-Header"

	// Build backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(headerKey) != expectedHeaderValue {
			t.Errorf("Expected header %s to be %q, got %q", headerKey, expectedHeaderValue, r.Header.Get(headerKey))
		}
		w.Header().Set(headerKey, expectedHeaderValue)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedResponseBody))
	}))
	defer backend.Close()

	// Build proxy request
	req := httptest.NewRequest("GET", "/", nil)
	// Set some headers
	req.Header.Set(headerKey, expectedHeaderValue)

	// Build response recorder to capture the response from the proxy handler
	rec := httptest.NewRecorder()

	// Call the proxy handler
	handlerOpts := HttpHandler{
		Host:    backend.URL,
		Res:     rec,
		Http:    req,
		Configs: hostConfig,
	}

	// Set another url on HttpHandler.Http to test that the handler uses the Host field instead of the original request URL
	handlerOpts.Http.URL.Path = "/test-path"

	HandleProxyRequest(handlerOpts)

	// Check the response
	response := rec.Result()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Compare the response body and headers
	if string(body) != expectedResponseBody {
		t.Fatalf("Expected response body %q, got %q", expectedResponseBody, string(body))
	}

	// Check status code
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Check custom header value
	if response.Header.Get(headerKey) != expectedHeaderValue {
		t.Fatalf("Expected header %s to be %q, got %q", headerKey, expectedHeaderValue, response.Header.Get(headerKey))
	}
}

func TestHandleProxyRequest_BackendError(t *testing.T) {
	err := config.GlobalConfig.Load("../file-loader/config.test.yaml")
	if err != nil {
		panic(err)
	}
	hostConfig, exists := config.GlobalConfig.GetItem("backend_url")
	if !exists {
		t.Fatalf("expected hostConfig to exist")
	}
	// Build proxy request
	req := httptest.NewRequest("GET", "/", nil)
	// Build response recorder to capture the response from the proxy handler
	rec := httptest.NewRecorder()
	// Call the proxy handler with an invalid backend URL to trigger an error
	handlerOpts := HttpHandler{
		Host:    "http://invalid-backend",
		Res:     rec,
		Http:    req,
		Configs: hostConfig,
	}
	HandleProxyRequest(handlerOpts)
	response := rec.Result()
	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}
}
