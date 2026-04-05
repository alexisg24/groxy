package ratelimiter

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IpRateLimiter struct {
	ips   map[string]*ipEntry
	mu    *sync.Mutex
	rate  rate.Limit
	burst int
	ttl   time.Duration
}

func NewIpRateLimiter(limit rate.Limit, burst int, ttl time.Duration) *IpRateLimiter {
	rateLimiter := &IpRateLimiter{
		ips:   make(map[string]*ipEntry),
		mu:    &sync.Mutex{},
		rate:  limit,
		burst: burst,
		ttl:   ttl, // Set TTL for IP entries
	}
	fmt.Printf("Initialized IpRateLimiter with Rate: %v, Burst: %d, TTL: %v\n", limit, burst, ttl)
	go rateLimiter.janitor() // Start the janitor goroutine

	return rateLimiter
}

func (rl *IpRateLimiter) AddIp(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.ips[ip] = &ipEntry{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

func (rl *IpRateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if entry, exists := rl.ips[ip]; exists {
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.ips[ip] = &ipEntry{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

func (rl *IpRateLimiter) janitor() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, entry := range rl.ips {
			if now.Sub(entry.lastSeen) > rl.ttl {
				delete(rl.ips, ip)
			}
		}
		rl.mu.Unlock()
	}
}
