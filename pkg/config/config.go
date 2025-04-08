package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadConfig(configPath string) (*GatewayConfig, error) {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("config file does not exist at %s\n", configPath)
		os.Exit(1)
	}

	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer jsonFile.Close()

	config := &GatewayConfig{}
	decoder := json.NewDecoder(jsonFile)
	decoder.DisallowUnknownFields() // This will cause an error if JSON contains fields not in struct

	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Additional validation
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *GatewayConfig) error {
	// Validate required fields
	if config.Gateway.Port <= 0 {
		return fmt.Errorf("invalid port number: %d", config.Gateway.Port)
	}

	if len(config.Services) == 0 {
		return fmt.Errorf("no services defined in configuration")
	}

	//TODO: Validate each service
	// for i, service := range config.Services {
	// 	if service.Name == "" {
	// 		return fmt.Errorf("service at index %d has no name", i)
	// 	}
	// 	if service.URL == "" {
	// 		return fmt.Errorf("service '%s' has no URL", service.Name)
	// 	}
	// 	if service.Timeout <= 0 {
	// 		return fmt.Errorf("service '%s' has invalid timeout: %d", service.Name, service.Timeout)
	// 	}
	// }

	return nil
}
