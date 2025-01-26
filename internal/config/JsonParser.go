package config

import (
	"encoding/json"
	"os"
)



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