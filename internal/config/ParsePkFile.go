package config

import (
	"encoding/json"
	"fmt"
	"os"
)


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