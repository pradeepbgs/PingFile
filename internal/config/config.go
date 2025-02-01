package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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



func ParsePKFile(filepath string) (*APIConfig,error)  {
	
	data,err := os.ReadFile(filepath)
	if err != nil{
		return nil,  fmt.Errorf("error opening file: %v", err)
	}
	
	var config APIConfig
	if err = json.Unmarshal(data,&config); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &config , nil
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


type CookieData struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Expires  string `json:"expires,omitempty"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"http_only"`
}

func ParseCookie (filename string) ([]*http.Cookie , error) {
	file,err := os.Open(filename)
	if err != nil {
		return nil,err
	}
	defer file.Close()

	var cookieList []CookieData
	if err := json.NewDecoder(file).Decode(&cookieList); err != nil {
		return nil,err
	}

	var cookies []*http.Cookie
	for _, cookieData := range cookieList {

		if cookieData.Name == "" || cookieData.Value == "" {
			continue
		}

		var expires time.Time
		if cookieData.Expires != ""{
			expires ,_ = time.Parse(time.RFC1123, cookieData.Expires)
		}

		cookies = append(cookies, &http.Cookie{
			Name: cookieData.Name,
			Value: cookieData.Value,
			Path: cookieData.Path,
			Domain: cookieData.Domain,
			Expires: expires,
			Secure: cookieData.Secure,
			HttpOnly: cookieData.HttpOnly,
		})
	}

	return cookies,nil
	
}


func SaveCookies(filename string, cookies []*http.Cookie) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return fmt.Errorf("failed to open or create file: %w", err)
	}
	defer file.Close()

	existingCookie, _ := ParseCookie("root.cookie.pkfile")

	cookieMap := make(map[string]*http.Cookie)
	for _, c := range existingCookie {
		cookieMap[c.Name] = c
	}

	for _, c := range cookies {
		cookieMap[c.Name] = c
	}

	var updatedCookie []*http.Cookie
	for _, c := range cookieMap {
		updatedCookie = append(updatedCookie, c)
	}

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(updatedCookie)
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
		return fmt.Errorf("unsupported file extension: %s , please user .josn , .yaml or .pkfile", FileExtension)
	}

}