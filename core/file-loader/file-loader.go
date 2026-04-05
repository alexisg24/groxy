package fileloader

import (
	"os"
	"time"

	ratelimiter "github.com/alexisg24/groxy/core/rate-limiter"
	"golang.org/x/time/rate"
	"gopkg.in/yaml.v3"
)

type FileLoaderItem struct {
	Name                          string            `yaml:"name"`
	Path                          string            `yaml:"path"`
	InternalHeaders               []string          `yaml:"internal_headers"`
	Timeout                       int               `yaml:"timeout"`
	HealthCheckPath               string            `yaml:"health_check_path"`
	HealthCheckInterval           int               `yaml:"health_check_interval"`
	HealthCheckTimeout            int               `yaml:"health_check_timeout"`
	HealthCheckCustomHeaders      map[string]string `yaml:"health_check_custom_headers"`
	HealthCheckExpectedStatus     int               `yaml:"health_check_expected_status"`
	RetryCount                    int               `yaml:"retry_count"`
	IsHealthy                     bool
	HasHealthCheck                bool
	RateLimitMaxRequestsPerSecond float64 `yaml:"rate_limit_max_requests_per_second"`
	RateLimitBurst                int     `yaml:"rate_limit_burst"`
	RateLimitTTLMinutes           int     `yaml:"rate_limit_ttl_minutes"`
	RateLimiterT                  *ratelimiter.IpRateLimiter
	RateLimiter                   *ratelimiter.IpRateLimiter
	CustomRequestHeaders          map[string]string `yaml:"custom_request_headers"`
	PayloadSizeLimit              int               `yaml:"payload_size_limit"`
}

type internalFileLoaderItem struct {
	Hosts []FileLoaderItem `yaml:"hosts"`
	Port  int              `yaml:"port"`
}

type GlobalFileLoaderConfig struct {
	Hosts map[string]*FileLoaderItem `yaml:"hosts"`
	Port  int                        `yaml:"port"`
}

func (f *GlobalFileLoaderConfig) Load(path string) error {
	// Read the file content
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var items internalFileLoaderItem

	// Unmarshal the YAML content into the struct
	err = yaml.Unmarshal(data, &items)
	if err != nil {
		return err
	}

	// Set the port from the loaded config
	f.Port = items.Port

	// Convert the slice of items into a map for easier access
	f.Hosts = make(map[string]*FileLoaderItem)
	for _, item := range items.Hosts {
		// Default values for missing fields can be set here if needed
		if item.Timeout == 0 {
			item.Timeout = 5000
		}
		if item.HealthCheckInterval == 0 {
			item.HealthCheckInterval = 5000
		}
		if item.HealthCheckTimeout == 0 {
			item.HealthCheckTimeout = 5000
		}
		if item.HealthCheckExpectedStatus == 0 {
			item.HealthCheckExpectedStatus = 200
		}
		if item.RetryCount == 0 {
			item.RetryCount = 3
		}
		if item.HealthCheckPath != "" {
			item.HasHealthCheck = true
		}
		if item.RateLimitTTLMinutes == 0 {
			item.RateLimitTTLMinutes = 5
		}
		if item.RateLimitMaxRequestsPerSecond > 0 && item.RateLimitBurst > 0 {
			maxPerSecond := item.RateLimitMaxRequestsPerSecond
			burst := item.RateLimitBurst
			ttl := time.Duration(item.RateLimitTTLMinutes) * time.Minute
			item.RateLimiter = ratelimiter.NewIpRateLimiter(rate.Limit(maxPerSecond), burst, ttl)
		} else {
			item.RateLimitMaxRequestsPerSecond = 100
			item.RateLimitBurst = 20
			maxPerSecond := item.RateLimitMaxRequestsPerSecond
			burst := item.RateLimitBurst
			ttl := time.Duration(item.RateLimitTTLMinutes) * time.Minute
			item.RateLimiter = ratelimiter.NewIpRateLimiter(rate.Limit(maxPerSecond), burst, ttl)

		}
		item.IsHealthy = true // Assume healthy at startup, health checks will update this

		if item.PayloadSizeLimit == 0 {
			item.PayloadSizeLimit = 1048576 // Default to 1MB if not set
		}

		if item.CustomRequestHeaders == nil {
			item.CustomRequestHeaders = make(map[string]string)
		}

		f.Hosts[item.Name] = &item
	}

	return nil
}

func (f *GlobalFileLoaderConfig) GetItem(name string) (*FileLoaderItem, bool) {
	item, exists := f.Hosts[name]
	return item, exists
}

func (f *GlobalFileLoaderConfig) GetAllItems() map[string]*FileLoaderItem {
	return f.Hosts
}

func (f *GlobalFileLoaderConfig) GetPort() int {
	return f.Port
}
