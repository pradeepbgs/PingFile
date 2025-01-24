package config

import (
	"bufio"
	"fmt"
	"os"
)


func ParsePKFile(filepath string) (*APIConfig,error)  {
	
	file,err := os.Open(filepath)
	if err != nil{
		return nil,  fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	return nil , nil
}