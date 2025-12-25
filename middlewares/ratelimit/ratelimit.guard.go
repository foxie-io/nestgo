package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/foxie-io/ng"
)

// Ensure Guard implements the required interfaces
var _ interface {
	ng.ID
	ng.Guard
} = (*Guard)(nil)

type (
	// Config holds the rate limit configuration.
	Config struct {
		// Limit is the maximum number of requests allowed in the window.
		Limit int

		// Window is the duration for which the rate limit applies.
		Window time.Duration

		// Identifier generates a unique identifier for the client making the request.
		Identifier func(ctx context.Context) string

		// ErrorHandler is called when the rate limit is exceeded.
		ErrorHandler func(ctx context.Context) error

		// SetHeaderHandler is called to set rate limit headers in the response.
		SetHeaderHandler func(ctx context.Context, key, value string)

		// MetadataKey is the key used to store the config in route metadata.
		MetadataKey string
	}

	// clientData holds the rate limit data for a specific client.
	clientData struct {

		// reqs is the number of requests made by the client in the current window.
		reqs int

		// resetAt is the time when the rate limit window resets.
		resetAt time.Time

		// limit is the maximum number of requests allowed in the window.
		limit int
	}

	// Guard is the rate limiting middleware.
	Guard struct {
		ng.DefaultID[Guard]
		clients map[string]*clientData
		mutex   sync.RWMutex
		config  *Config
	}
)

// DefaultConfig provides default settings for the rate limiter.
var DefaultConfig = &Config{
	Limit:  100,
	Window: time.Minute,
	Identifier: func(ctx context.Context) string {
		return "default-client"
	},
	ErrorHandler: func(ctx context.Context) error {
		return errors.New("rate limit exceeded")
	},
}

// New creates a new Guard instance
func New(config *Config) *Guard {
	if config == nil {
		config = DefaultConfig
	}

	guard := &Guard{
		config:  overrideOptional(config, DefaultConfig),
		clients: make(map[string]*clientData),
	}

	guard.startCleanup(time.Minute)
	return guard
}

// Allow checks if the request is allowed under the rate limit.
func (g *Guard) Allow(ctx context.Context) error {
	config := g.config

	// Check if there is a route-specific rate limit configuration
	if routeConfig, ok := GetConfig(ctx, config.MetadataKey); ok {
		config = overrideOptional(routeConfig, g.config)
	}

	// Generate a unique identifier for the client
	id := config.Identifier(ctx)

	// Lock the mutex to ensure thread-safe access to the clients map
	g.mutex.Lock()
	defer g.mutex.Unlock()

	now := time.Now()
	client, exists := g.clients[id]

	// If the client does not exist or their rate limit window has expired, reset their data
	if !exists || now.After(client.resetAt) {
		client = &clientData{
			limit:   config.Limit,
			reqs:    1,
			resetAt: now.Add(config.Window),
		}
		g.clients[id] = client
	} else {
		client.reqs++
	}

	// Set rate limit headers
	if config.SetHeaderHandler != nil {
		remaining := client.limit - client.reqs
		if remaining < 0 {
			remaining = 0
		}
		config.SetHeaderHandler(ctx, "Limit", fmt.Sprintf("%d", client.limit))
		config.SetHeaderHandler(ctx, "Remaining", fmt.Sprintf("%d", remaining))
		config.SetHeaderHandler(ctx, "Reset", client.resetAt.Format(time.RFC3339))
	}

	hasReachedLimit := client.reqs > config.Limit
	if hasReachedLimit {
		client.reqs--
		return config.ErrorHandler(ctx)
	}

	return nil
}

// cleanupExpiredClients removes clients whose rate limit window has expired
func (g *Guard) cleanupExpiredClients() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	now := time.Now()
	for id, client := range g.clients {
		if now.After(client.resetAt) {
			delete(g.clients, id)
		}
	}
}

func (g *Guard) startCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			g.cleanupExpiredClients()
		}
	}()
}

func overrideOptional(config *Config, defaultConfig *Config) *Config {
	if config == nil {
		return defaultConfig
	}

	if config.Limit == 0 {
		config.Limit = defaultConfig.Limit
	}

	if config.Window == 0 {
		config.Window = defaultConfig.Window
	}

	if config.Identifier == nil {
		config.Identifier = defaultConfig.Identifier
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultConfig.ErrorHandler
	}

	if config.MetadataKey == "" {
		config.MetadataKey = defaultConfig.MetadataKey
	}

	if config.SetHeaderHandler == nil {
		config.SetHeaderHandler = defaultConfig.SetHeaderHandler
	}

	return config
}
