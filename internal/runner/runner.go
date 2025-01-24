package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
	"github.com/pradeepbgs/pingfile/internal/config"
)

func ExecuteAPI(config *config.APIConfig) error {
	// Prepare request body
	var bodyBytes []byte
	if len(config.Body) > 0 {
		bodyBytes, _ = json.Marshal(config.Body)
	}

	// Create the HTTP request
	req, err := http.NewRequest(config.Headers["Method"], config.URL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers to the request
	for key, value := range config.Headers {
		if key != "Method" {
			req.Header.Set(key, value)
		}
	}

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Define colors for output
	statusColor := color.New(color.FgGreen).SprintFunc() // Green for status
	headerColor := color.New(color.FgCyan).SprintFunc()  // Cyan for headers
	bodyColor := color.New(color.FgYellow).SprintFunc()  // Yellow for body
	errorColor := color.New(color.FgRed).SprintFunc()    // Red for errors

	// Print the response details in a cool way
	fmt.Println()
	fmt.Printf("%s: %s\n", statusColor("Status Code"), resp.Status)
	fmt.Println(headerColor("\nHeaders:"))
	for key, values := range resp.Header {
		fmt.Printf("  %s: %s\n", key, values)
	}
	fmt.Println(bodyColor("\nBody:"))
	fmt.Printf("%s\n\n", string(respBody))

	// Return error if status code is not 2xx
	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s: %s", errorColor("Error"), resp.Status)
	}

	return nil
}
