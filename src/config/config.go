package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Credentials struct {
	Type     string `json:"type" yaml:"type"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Token    string `json:"token" yaml:"token"`
}

type FileItem struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
}

type APIConfig struct {
	Name               string                 `json:"name" yaml:"name"`
	Description        string                 `json:"description" yaml:"description"`
	Run 			   *bool					  `json:"run" yaml:"run"`
	SaveResponse       bool                   `json:"saveResponse" yaml:"saveResponse"`
	FilePath           string                 `json:"filePath" yaml:"filePath"`
	IncludeCookie      *bool                  `json:"includeCookie" yaml:"includeCookie"`
	IncludeCredentials bool                   `json:"includeCredentials" yaml:"includeCredentials"`
	URL                string                 `json:"url" yaml:"url"`
	Headers            map[string]string      `json:"headers" yaml:"headers"`
	Body               map[string]interface{} `json:"body" yaml:"body"`
	File               []FileItem             `json:"file" yaml:"file"`
	Credentials        *Credentials           `json:"credentials" yaml:"credentials"`
}

type GroupApiConfig struct {
	Name        string      `json:"name" yaml:"name"`
	Description string      `json:"description" yaml:"description"`
	Version     string      `json:"version" yaml:"version"`
	BaseUrl     string      `json:"baseUrl" yaml:"baseUrl"`
	APIs        []APIConfig `json:"apis" yaml:"apis"`
}

func GetFileExtension(filePath string) string {
	return strings.ToLower(filepath.Ext(filePath))
}

func replaceEnvVars(value string) string {
	return os.ExpandEnv(value)
}

func Parser(filepath string) (interface{}, error) {
	if filepath == "" {
		return nil, errors.New("file path is empty")
	}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filepath)
	}

	ext := GetFileExtension(filepath)
	var result interface{}
	var err error

	switch ext {
	case ".json":
		result, err = ParseJSON(filepath)
	case ".yaml", ".yml":
		result, err = ParseYAML(filepath)
	case ".pkfile":
		result, err = ParsePKFile(filepath)
	default:
		return nil, fmt.Errorf("unsupported file format %s; supported formats: pkfile, yaml, json", ext)
	}

	if err != nil {
		return nil, err
	}

	if group, ok := result.(*GroupApiConfig); ok {
		if len(group.APIs) == 0 {
			return nil, fmt.Errorf("group API configuration is empty")
		}
		for i := range group.APIs {
			if group.APIs[i].Credentials == nil {
				group.APIs[i].Credentials = &Credentials{}
			}
			group.APIs[i].Credentials.Username = replaceEnvVars(group.APIs[i].Credentials.Username)
			group.APIs[i].Credentials.Password = replaceEnvVars(group.APIs[i].Credentials.Password)
			group.APIs[i].Credentials.Token = replaceEnvVars(group.APIs[i].Credentials.Token)
		}
	} else if api, ok := result.(*APIConfig); ok {
		if api.Credentials == nil {
			api.Credentials = &Credentials{}
		}
		api.Credentials.Username = replaceEnvVars(api.Credentials.Username)
		api.Credentials.Password = replaceEnvVars(api.Credentials.Password)
		api.Credentials.Token = replaceEnvVars(api.Credentials.Token)
	}

	return result, nil
}

func ParseYAML(filePath string) (interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var groupConfig GroupApiConfig
	if err := yaml.Unmarshal(data, &groupConfig); err == nil && len(groupConfig.APIs) > 0 {
		return &groupConfig, nil
	}

	var config APIConfig
	if err := yaml.Unmarshal(data, &config); err == nil {
		return &config, nil
	}

	return nil, fmt.Errorf("invalid YAML configuration format")
}

func ParsePKFile(filepath string) (interface{}, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	var groupConfig GroupApiConfig
	if err := json.Unmarshal(data, &groupConfig); err == nil && len(groupConfig.APIs) > 0 {
		return &groupConfig, nil
	}

	var apiConfig APIConfig
	if err := json.Unmarshal(data, &apiConfig); err == nil {
		return &apiConfig, nil
	}

	return nil, fmt.Errorf("invalid configuration file format")
}

func ParseJSON(filepath string) (interface{}, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var groupConfig GroupApiConfig
	if err := json.Unmarshal(data, &groupConfig); err == nil && len(groupConfig.APIs) > 0 {
		return &groupConfig, nil
	}

	var apiConfig APIConfig
	if err := json.Unmarshal(data, &apiConfig); err == nil {
		return &apiConfig, nil
	}

	return nil, fmt.Errorf("invalid configuration file format")
}

type CookieData struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Expires  string `json:"expires,omitempty"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"http_only"`
}

func ParseCookie(filename string) ([]*http.Cookie, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cookieList []CookieData
	if err := json.NewDecoder(file).Decode(&cookieList); err != nil {
		return nil, err
	}

	var cookies []*http.Cookie
	for _, cookieData := range cookieList {
		if cookieData.Name == "" || cookieData.Value == "" {
			continue
		}

		var expires time.Time
		if cookieData.Expires != "" {
			expires, _ = time.Parse(time.RFC1123, cookieData.Expires)
		}

		cookies = append(cookies, &http.Cookie{
			Name:     cookieData.Name,
			Value:    cookieData.Value,
			Path:     cookieData.Path,
			Domain:   cookieData.Domain,
			Expires:  expires,
			Secure:   cookieData.Secure,
			HttpOnly: cookieData.HttpOnly,
		})
	}

	return cookies, nil
}

func SaveCookies(filename string, cookies []*http.Cookie) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %w", err)
	}
	defer file.Close()

	existingCookies, _ := ParseCookie(filename)

	cookieMap := make(map[string]*http.Cookie)
	for _, c := range existingCookies {
		cookieMap[c.Name] = c
	}

	for _, c := range cookies {
		cookieMap[c.Name] = c
	}

	var updatedCookies []*http.Cookie
	for _, c := range cookieMap {
		updatedCookies = append(updatedCookies, c)
	}

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(updatedCookies)
}

func SaveResponseToFile(filename string, requestDetails map[string]interface{}, responseDetails map[string]interface{}) error {
	data := map[string]interface{}{
		"request":  requestDetails,
		"response": responseDetails,
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	FileExtension := GetFileExtension(filename)

	switch FileExtension {
	case ".json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", " ")
		return encoder.Encode(data)

	case ".yaml":
		encoder := yaml.NewEncoder(file)
		return encoder.Encode(data)

	case ".pkfile":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", " ")
		return encoder.Encode(data)

	default:
		return fmt.Errorf("unsupported file extension: %s, please use .json, .yaml, or .pkfile", FileExtension)
	}
}
