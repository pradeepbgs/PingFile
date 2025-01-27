package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pradeepbgs/pingfile/internal/config"
	"gopkg.in/yaml.v3"
)

func SaveCookies(filename string, cookies []*http.Cookie) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return fmt.Errorf("failed to open or create file: %w", err)
	}
	defer file.Close()

	existingCookie, _ := config.ParseCookie("root.cookie.pkfile")

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
	FileExtension := config.GetFileExtension(filename)

	switch FileExtension {
	case ".json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", " ")
		return encoder.Encode(data)
	
	case ".yaml":
		encoder := yaml.NewEncoder(file)
		// encoder.SetIndent(1)
		return encoder.Encode(data)
	
	case ".pkfile":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", " ")
		return encoder.Encode(data)
	
	default:
		return fmt.Errorf("unsupported file extension: %s , please user .josn , .yaml or .pkfile", FileExtension)
	}
	 
}

func ExecuteAPI(apiConfig *config.APIConfig, saveResponses bool, cookie []*http.Cookie) error {
	if apiConfig.Headers["Method"] == "" {
		return fmt.Errorf("HTTP method not specified")
	}
	if apiConfig.URL == "" {
		return fmt.Errorf("URL not specified")
	}

	// Prepare request body
	var bodyBytes []byte
	if len(apiConfig.Body) > 0 {
		bodyBytes, _ = json.Marshal(apiConfig.Body)
	}

	// Create the HTTP request
	req, err := http.NewRequest(apiConfig.Headers["Method"], apiConfig.URL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if apiConfig.IncludeCredentials && apiConfig.Credentials != nil {
		switch apiConfig.Credentials.Type {
		case "basic":
			req.SetBasicAuth(apiConfig.Credentials.Username, apiConfig.Credentials.Password)
		case "bearer":
			req.Header.Set("Authorization", "Bearer "+apiConfig.Credentials.Token)
		default:
			return fmt.Errorf("unsupported credential type: %s", apiConfig.Credentials.Type)
		}
	}

	// Add headers to the request
	for key, value := range apiConfig.Headers {
		if key != "Method" {
			req.Header.Set(key, value)
		}
	}

	// add cookies in req header
	includeCookie := true
	if apiConfig.IncludeCookie != nil {
		includeCookie = *apiConfig.IncludeCookie
	}

	if includeCookie && cookie != nil {
		for _, c := range cookie {
			req.AddCookie(c)
		}
	}

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Define colors for output
	statusColor := color.New(color.FgGreen).SprintFunc()
	headerColor := color.New(color.FgCyan).SprintFunc()
	bodyColor := color.New(color.FgYellow).SprintFunc()
	errorColor := color.New(color.FgRed).SprintFunc()

	// Print the response details in a cool way
	fmt.Println()
	fmt.Printf("%s: %s\n", statusColor("Status Code"), resp.Status)
	fmt.Println(headerColor("\nHeaders:"))
	for key, values := range resp.Header {
		fmt.Printf("  %s: %s\n", key, values)
	}
	fmt.Println(bodyColor("\nBody:"))

	var responseBodyBytes bytes.Buffer
	chunk := make([]byte, 4096)

	for {
		n, err := resp.Body.Read(chunk)
		if n > 0 {
			responseBodyBytes.Write(chunk[:n])
			fmt.Print(string(chunk))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
	}

	if saveResponses || apiConfig.SaveResponse {
		requestDetails := map[string]interface{}{
			"URL":     apiConfig.URL,
			"Headers": apiConfig.Headers,
			"Body":    apiConfig.Body,
		}
		responseDetails := map[string]interface{}{
			"Status":  resp.Status,
			"Headers": resp.Header,
			"Body":    responseBodyBytes.String(),
		}

		timestamp := time.Now().Format("20060102_150405")
		saveFilePath := fmt.Sprintf("response_%s_%s_%s.pkfile", apiConfig.Headers["Method"], strings.ReplaceAll(apiConfig.URL, "/", "_"), timestamp)

		if apiConfig.FilePath != "" {
			saveFilePath = apiConfig.FilePath
		}

		err := SaveResponseToFile(saveFilePath, requestDetails, responseDetails)
		if err != nil {
			return fmt.Errorf("\nfailed to save response: %w", err)
		}
		fmt.Println("\nResponse saved to response.json")
	}

	cookies := resp.Cookies()
	if len(cookies) > 0 {
		SaveCookies("root.cookie.pkfile", cookies)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s: %s", errorColor("Error"), resp.Status)
	}

	return nil
}
