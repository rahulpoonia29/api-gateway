package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
)

// LoadConfig loads and validates the gateway configuration from the specified file
func LoadConfig(configPath string, logger *slog.Logger) (*GatewayConfig, error) {
	if configPath == "" {
		configPath = "./config.json"
		logger.Info("Using default config path", "path", configPath)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", configPath)
	}

	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer jsonFile.Close()

	config := &GatewayConfig{}
	decoder := json.NewDecoder(jsonFile)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if err := validateConfig(config, logger); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *GatewayConfig, logger *slog.Logger) error {
	if err := validateGatewaySettings(&config.Gateway); err != nil {
		return fmt.Errorf("gateway settings validation failed: %w", err)
	}

	if len(config.Services) == 0 {
		return fmt.Errorf("no services defined in configuration")
	}

	for i, service := range config.Services {
		if err := validateServiceConfig(&service, i, logger); err != nil {
			return err
		}
	}

	return nil
}

func validateGatewaySettings(settings *GatewaySettings) error {
	if settings.Port <= 0 || settings.Port > 65535 {
		return fmt.Errorf("invalid port number: %d (must be between 1-65535)", settings.Port)
	}

	if settings.LogLevel != "" {
		switch settings.LogLevel {
		case Debug, Info, Warn, Error:
		default:
			return fmt.Errorf("invalid log level: %s (must be one of: debug, info, warn, error)", settings.LogLevel)
		}
	} else {
		settings.LogLevel = Info
	}

	return nil
}

func validateServiceConfig(service *ServiceConfig, index int, logger *slog.Logger) error {
	if service.Name == "" {
		return fmt.Errorf("service at index %d has no name", index)
	}

	if service.Proxy.ListenPath == "" {
		return fmt.Errorf("service '%s' has no listen path", service.Name)
	}

	if !strings.HasPrefix(service.Proxy.ListenPath, "/") {
		return fmt.Errorf("service '%s' listen path must start with a '/'", service.Name)
	}

	if !service.Enabled {
		logger.Debug("Service is disabled", "service", service.Name)
	}

	return validateUpstreamConfig(&service.Proxy.Upstream, service.Name)
}

func validateUpstreamConfig(upstream *UpstreamConfig, serviceName string) error {
	if len(upstream.Targets) == 0 {
		return fmt.Errorf("service '%s' has no upstream targets", serviceName)
	}

	for i, target := range upstream.Targets {
		if _, err := url.Parse(target); err != nil {
			return fmt.Errorf("service '%s' has invalid target URL at index %d: %s",
				serviceName, i, err.Error())
		}
	}

	switch upstream.Balancing {
	case RoundRobin, LeastConn, IPHash:
	case "":
		upstream.Balancing = RoundRobin
	default:
		return fmt.Errorf("service '%s' has invalid balancing strategy: %s (must be one of: roundrobin, least_conn, ip_hash)",
			serviceName, upstream.Balancing)
	}

	return nil
}
