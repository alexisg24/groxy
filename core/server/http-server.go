package server

import (
	"fmt"
	logger "log"
	"net/http"
	"time"

	"github.com/alexisg24/groxy/core/config"
	"github.com/alexisg24/groxy/core/helpers"
	ratelimiter "github.com/alexisg24/groxy/core/rate-limiter"
)

func StartHttpServer() {
	port := config.GlobalConfig.GetPort()
	if port == 0 {
		port = 6969
	}
	// Setup internal rate limiter
	rateLimitConfig := ratelimiter.NewIpRateLimiter(10, 20, 5*time.Minute)

	// Generate a handler for every registered host
	for _, hostConfig := range config.GlobalConfig.GetAllItems() {
		path := fmt.Sprintf("/%s/", hostConfig.Name)
		logger.Printf("Route: %s, registered. On path: %s", hostConfig.Name, hostConfig.Path)
		// Handle subroutes
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			RouteHandler(w, r, hostConfig)
		})

		// Redirect root path to the host path
		http.HandleFunc(fmt.Sprintf("/%s", hostConfig.Name), func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, path, http.StatusMovedPermanently)
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		limiter := rateLimitConfig.GetLimiter(helpers.GetIP(r))
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests: Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		http.NotFound(w, r)
	})

	logger.Printf("Starting HTTP server on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Printf("Error starting HTTP server: %v", err)
		panic(err)
	}
}
