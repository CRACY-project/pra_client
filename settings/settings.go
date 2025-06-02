package settings

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	configFilePath = "/etc/jimber/settings_pra.json" // Change to your actual path
)

type Settings struct {
	Token string `json:"token"`
	// Add more fields as needed
}

// LoadSettings reads and parses the config file.
func LoadSettings() (*Settings, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	return &s, nil
}
