package helpers

import (
	"net/http"
)

// Supported Headers:
// - X-Forwarded-For (client IP address)
// - X-Forwarded-Host (original host header)
// - X-Forwarded-Proto (original request protocol)
// - X-Forwarded-Port (original request port)
// - X-Real-Ip (client IP address)
func SetInternalHeaders(r *http.Request, headers []string) {
	for _, header := range headers {
		switch header {
		case "X-Forwarded-For":
			r.Header.Set("X-Forwarded-For", r.RemoteAddr)
		case "X-Forwarded-Host":
			r.Header.Set("X-Forwarded-Host", r.Host)
		case "X-Forwarded-Proto":
			r.Header.Set("X-Forwarded-Proto", r.URL.Scheme)
		case "X-Forwarded-Port":
			r.Header.Set("X-Forwarded-Port", r.URL.Port())
		case "X-Real-Ip":
			r.Header.Set("X-Real-Ip", GetIP(r))
		}
	}
}
