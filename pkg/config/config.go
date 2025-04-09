package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// InitConfig identifies the config file path and verifies it exists
func InitConfig(cfgFile string) {
	if cfgFile == "" {
		cfgFile = "./config.json"
		fmt.Printf("warning: config file not specified, using default path: %s\n", cfgFile)
	}

	// Verify config file exists
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Printf("config file %s does not exist\n", cfgFile)
		os.Exit(1)
	}

	fmt.Printf("using config file: %s\n", cfgFile)
}

func LoadConfig(configPath string) (*GatewayConfig, error) {
	// Check if config file exists
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

	// TODO: Validate each service
	// for i, service := range config.Services {
	//   if service.Name == "" {
	//     return fmt.Errorf("service at index %d has no name", i)
	//   }
	//   if service.URL == "" {
	//     return fmt.Errorf("service '%s' has no url", service.Name)
	//   }
	//   if service.Timeout <= 0 {
	//     return fmt.Errorf("service '%s' has invalid timeout: %d", service.Name, service.Timeout)
	//   }
	// }

	return nil
}
