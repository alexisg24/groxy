package helpers

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xr := r.Header.Get("X-Real-Ip"); xr != "" {
		return strings.TrimSpace(xr)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	fmt.Printf("Extracted IP: %s from RemoteAddr: %s\n", host, r.RemoteAddr)
	return host
}
