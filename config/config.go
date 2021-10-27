package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-articles-controller
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	Debug                      bool          `envconfig:"DEBUG"`
	SiteDomain                 string        `envconfig:"SITE_DOMAIN"`
	PatternLibraryAssetsPath   string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	APIRouterURL               string        `envconfig:"API_ROUTER_URL"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	cfg, err := get()
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		cfg.PatternLibraryAssetsPath = "http://localhost:9000/dist"
	} else {
		cfg.PatternLibraryAssetsPath = "//cdn.ons.gov.uk/sixteens/f80be2c"
	}
	return cfg, nil
}

func get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   ":26500",
		Debug:                      false,
		SiteDomain:                 "localhost",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		APIRouterURL:               "http://localhost:23200/v1",
	}

	return cfg, envconfig.Process("", cfg)
}
