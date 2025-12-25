package ratelimit

import (
	"context"
	"fmt"

	"github.com/foxie-io/ng"
)

// WithConfig sets the rate limit config in the route metadata.
func WithConfig(config *Config) ng.Option {
	key := fmt.Sprintf("RATE_LIMIT:%s", config.MetadataKey)
	return ng.WithMetadata(key, config)
}

// SkipRateLimit is used to skip rate limiting for a specific route.
func SkipRateLimit() ng.Option {
	return ng.WithSkip(ng.DefaultID[Guard]{})
}

// GetConfig gets the rate limit config from context metadata.
func GetConfig(ctx context.Context, metadateKey string) (*Config, bool) {
	key := fmt.Sprintf("RATE_LIMIT:%s", metadateKey)

	data, ok := ng.GetContext(ctx).Route().Core().Metadata(key)
	if !ok {
		return nil, false
	}

	config, ok := data.(*Config)
	if !ok {
		return nil, false
	}

	return config, true
}
