package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

func ParseYAML(filePath string) (*APIConfig, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var config APIConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}