package config

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func GetFileExtension(filePath string) string {
    return strings.ToLower(filepath.Ext(filePath))
}

func Parser(filepath string) (*APIConfig,error) {
	if filepath == "" {
		log.Fatal("File path is empty")
		return nil, fmt.Errorf("File path is empty")
	}
	
	ext := GetFileExtension(filepath)
    switch ext {
    case ".json":
        return ParseJSON(filepath)
    case ".yaml", ".yml":
        return ParseYAML(filepath)
    case ".pkfile":
        return ParsePKFile(filepath)
    default:
        return nil, errors.New("unsupported file format, we only support pkfile,yaml & json")
	}
}