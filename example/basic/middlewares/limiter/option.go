package limiter

import (
	"context"

	"github.com/foxie-io/ng"
)

type configKey struct{}

func WithConfig(config *Config) ng.Option {
	return ng.WithMetadata(configKey{}, config)
}

func GetConfig(ctx context.Context) (config *Config, ok bool) {
	rc := ng.GetContext(ctx)
	metadata, exists := rc.Route().Core().Metadata(configKey{})
	if !exists {
		return nil, false
	}

	if ratelimitConfig, ok := metadata.(*Config); ok {
		return ratelimitConfig, true
	}

	return nil, false
}
