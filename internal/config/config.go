package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg Config) SetUser(u string) error {
	f, err := Read()
	if err != nil {
		return err
	}
	f.CurrentUserName = u
	if err = write(f); err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Unable to get Home Directory")
	}
	return fmt.Sprintf("%v/%v", dir, configFileName), nil
}

func Read() (*Config, error) {
	p, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	f, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("Unable to read config file")
	}
	var data Config
	if err = json.Unmarshal(f, &data); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Config JSON")
	}

	return &data, nil
}

func write(cfg *Config) error {
	jsonBlob, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Failed to Marshal Write")
	}
	p, _ := getConfigFilePath()
	if err = os.WriteFile(p, jsonBlob, 0644); err != nil {
		return fmt.Errorf("Failed to write to file")
	}

	return nil
}
