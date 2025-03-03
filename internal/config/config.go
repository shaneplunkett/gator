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

func (c Config) SetUser() {

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
		return nil, err
	}

	return &data, nil
}
