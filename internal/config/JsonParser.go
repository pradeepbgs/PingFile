package config

import (
	"encoding/json"
	"os"
)


type APIConfig struct {
	Name    string            `json:"name"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    map[string]interface{} `json:"body"`
}

func ParseJSON(filepath string) (*APIConfig,error) {
	data, err := os.ReadFile(filepath)
    if err != nil {
        return nil, err
    }

	var config APIConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config,nil
}