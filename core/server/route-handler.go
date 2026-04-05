package server

import (
	"fmt"
	"net/http"
	"strings"

	logger "log"

	fileloader "github.com/alexisg24/groxy/core/file-loader"
	"github.com/alexisg24/groxy/core/helpers"
	http_handler "github.com/alexisg24/groxy/core/http-handler"
)

// Global scope on fowarding function to allow for easier testing
var ProxyRequestFunc = http_handler.HandleProxyRequest

// RouteHandler handles the root path and forwards the request asynchronously.
func RouteHandler(w http.ResponseWriter, r *http.Request, hostConfig *fileloader.FileLoaderItem) {
	// Log the incoming request (IP - METHOD - PATH - String of headers - Body Size)
	logger.Printf("Received request from %s - %s %s - Headers: %v - Body Size: %d bytes", helpers.GetIP(r), r.Method, r.URL.Path, r.Header, r.ContentLength)

	// Check if the route is healty or return 502
	if !hostConfig.IsHealthy {
		logger.Printf("Backend %s is unhealthy, returning 502", hostConfig.Name)
		http.Error(w, "Bad Gateway: Backend is unhealthy", http.StatusBadGateway)
		return
	}
	// Check rate limit
	if hostConfig.RateLimiter == nil {
		logger.Printf("No rate limiter configured for backend %s, allowing all requests", hostConfig.Name)
		http.Error(w, "Internal Server Error: No rate limiter configured", http.StatusInternalServerError)
		return
	}

	limiter := hostConfig.RateLimiter.GetLimiter(helpers.GetIP(r))
	if !limiter.Allow() {
		logger.Printf("Rate limit exceeded for backend %s, returning 429", hostConfig.Name)
		http.Error(w, "Too Many Requests: Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Add size limit to request body
	r.Body = http.MaxBytesReader(w, r.Body, int64(hostConfig.PayloadSizeLimit))

	// Check if request body is larger than the payload size limit
	if r.ContentLength > int64(hostConfig.PayloadSizeLimit) {
		logger.Printf("Request body size %d exceeds payload size limit for backend %s, returning 413", r.ContentLength, hostConfig.Name)
		http.Error(w, "Payload Too Large: Request body exceeds the configured payload size limit", http.StatusRequestEntityTooLarge)
		return
	}

	// Set additional headers from config
	helpers.SetInternalHeaders(r, hostConfig.InternalHeaders)

	// Build the URL to forward the request to
	prefix := "/" + hostConfig.Name
	suffix := r.URL.Path
	if after, ok := strings.CutPrefix(suffix, prefix); ok {
		suffix = after
	}
	if suffix == "" {
		suffix = "/"
	}
	parsedUrl := strings.TrimRight(hostConfig.Path, "/") + suffix
	fmt.Printf("Forwarding request to %s\n", parsedUrl)
	handlerOpts := http_handler.HttpHandler{
		Host:    parsedUrl,
		Res:     w,
		Http:    r,
		Configs: hostConfig,
	}

	ProxyRequestFunc(handlerOpts)
}
