package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Credentials struct{
	Type string 	`json:"type" yaml:"type"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Token    string `json:"token" yaml:"token"`
}

type FileItem struct{
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
}

type APIConfig struct {
	Name         string                 `json:"name" yaml:"name"`
	SaveResponse bool                   `json:"saveResponse" yaml:"saveResponse"`
	FilePath     string                 `json:"filePath" yaml:"filePath"`
	IncludeCookie     *bool                  `json:"includeCookie" yaml:"includeCookie"`
	IncludeCredentials bool              `json:"includeCredentials" yaml:"includeCredentials"`
	URL          string                 `json:"url" yaml:"url"`
	Headers      map[string]string      `json:"headers" yaml:"headers"`
	Body         map[string]interface{} `json:"body" yaml:"body"`
	File      []FileItem             `json:"file" yaml:"file"`
	Credentials  *Credentials            `json:"credentials" yaml:"credentials"`
}

func GetFileExtension(filePath string) string {
    return strings.ToLower(filepath.Ext(filePath))
}

func replaceEnvVars(value string) string {
	resolved := os.ExpandEnv(value)
	return resolved
}

func Parser(filepath string) (*APIConfig,error) {
	if filepath == "" {
		log.Fatal("File path is empty")
		return nil, errors.New("file path is empty")
	}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil,fmt.Errorf("file not found: %s", filepath)
	}

	ext := GetFileExtension(filepath)
	var config *APIConfig
	var err error

    switch ext {
    case ".json":
        config, err = ParseJSON(filepath)
    case ".yaml", ".yml":
        config , err =  ParseYAML(filepath)
    case ".pkfile":
        config,err = ParsePKFile(filepath)
    default:
        return nil, fmt.Errorf("unsupported file format %s; supported formats: pkfile, yaml, json", ext)
	}

	if err != nil {
		return nil, err
	}

	if config.Credentials != nil {
		config.Credentials.Username = replaceEnvVars(config.Credentials.Username)
		config.Credentials.Password = replaceEnvVars(config.Credentials.Password)
		config.Credentials.Token = replaceEnvVars(config.Credentials.Token)
	}

	return config,nil
}