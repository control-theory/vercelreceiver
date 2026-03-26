package vercelreceiver

import (
	"fmt"
	"net"
)

type Config struct {
	// Single endpoint for all data types (host:port format)
	Endpoint string `mapstructure:"endpoint"`

	// Default secret for all endpoints (can be overridden per signal type)
	Secret string `mapstructure:"secret"`

	// Default signature algorithm for all endpoints (can be overridden per signal type)
	SignatureAlgorithm string `mapstructure:"signature_algorithm"`

	// Per-signal configuration (optional overrides for routes and secrets)
	Logs          SignalConfig `mapstructure:"logs"`
	Traces        SignalConfig `mapstructure:"traces"`
	SpeedInsights SignalConfig `mapstructure:"speed_insights"`
	WebAnalytics  SignalConfig `mapstructure:"web_analytics"`
}

// SignalConfig defines configuration for a specific signal type (logs, traces, metrics, analytics)
// Both route and secret are optional and will use defaults if not provided
type SignalConfig struct {
	Route  string `mapstructure:"route"`  // Optional: Override default route
	Secret string `mapstructure:"secret"` // Optional: Override default secret
}

// Validate validates the configuration
func (cfg *Config) Validate() error {
	// Validate endpoint if provided
	if cfg.Endpoint != "" {
		if err := validateEndpoint(cfg.Endpoint); err != nil {
			return fmt.Errorf("invalid endpoint: %w", err)
		}
	}

	// Set default routes if not provided
	if cfg.Logs.Route == "" {
		cfg.Logs.Route = "/logs"
	}
	if cfg.Traces.Route == "" {
		cfg.Traces.Route = "/traces"
	}
	if cfg.SpeedInsights.Route == "" {
		cfg.SpeedInsights.Route = "/speed-insights"
	}
	if cfg.WebAnalytics.Route == "" {
		cfg.WebAnalytics.Route = "/analytics"
	}

	// Set default signature algorithm if not provided
	if cfg.SignatureAlgorithm == "" {
		cfg.SignatureAlgorithm = "sha1"
	}

	// Validate allowed values for signature algorithm
	if !isValidSignatureAlgorithm(cfg.SignatureAlgorithm) {
		return fmt.Errorf("invalid signature_algorithm: %s", cfg.SignatureAlgorithm)
	}

	return nil
}

// isValidSignatureAlgorithm checks if the algorithm is allowed
func isValidSignatureAlgorithm(algo string) bool {
	switch algo {
	case "sha1", "sha256":
		return true
	default:
		return false
	}
}

// GetSignatureAlgorithm returns the signature algorithm (global for all signals)
func (cfg *Config) GetSignatureAlgorithm() string {
	if cfg.SignatureAlgorithm == "" {
		return "sha1"
	}
	return cfg.SignatureAlgorithm
}

// GetLogsSecret returns the secret for logs (per-signal override or default)
func (cfg *Config) GetLogsSecret() string {
	if cfg.Logs.Secret != "" {
		return cfg.Logs.Secret
	}
	return cfg.Secret
}

// GetTracesSecret returns the secret for traces (per-signal override or default)
func (cfg *Config) GetTracesSecret() string {
	if cfg.Traces.Secret != "" {
		return cfg.Traces.Secret
	}
	return cfg.Secret
}

// GetSpeedInsightsSecret returns the secret for speed insights (per-signal override or default)
func (cfg *Config) GetSpeedInsightsSecret() string {
	if cfg.SpeedInsights.Secret != "" {
		return cfg.SpeedInsights.Secret
	}
	return cfg.Secret
}

// GetWebAnalyticsSecret returns the secret for web analytics (per-signal override or default)
func (cfg *Config) GetWebAnalyticsSecret() string {
	if cfg.WebAnalytics.Secret != "" {
		return cfg.WebAnalytics.Secret
	}
	return cfg.Secret
}

// validateEndpoint validates that the endpoint is in the correct host:port format
func validateEndpoint(endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("endpoint cannot be empty")
	}

	// Check if it's a valid host:port format
	host, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		return fmt.Errorf("failed to split endpoint into 'host:port' pair: %w", err)
	}

	if host == "" {
		// Empty host is valid for Go net.Listen (means listen on all interfaces)
		// Only reject if both host and port are empty
		if port == "" {
			return fmt.Errorf("host cannot be empty when port is also empty")
		}
	}

	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	// Validate that the port is a valid number
	if _, err := net.LookupPort("tcp", port); err != nil {
		return fmt.Errorf("invalid port '%s': %w", port, err)
	}

	return nil
}
