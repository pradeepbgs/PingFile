package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pradeepbgs/pingfile/internal/config"
)

func Execute (config *config.APIConfig) error {
	var bodyBytes []byte
	if len(config.Body) > 0 {
		bodyBytes, _ = json.Marshal(config.Body)
	}
	fmt.Println(config.Headers["Method"])
	req, err := http.NewRequest(config.Headers["Method"], config.URL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range config.Headers {
		if key != "Method" {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Print response details
	fmt.Println("RES",resp)
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", string(respBody))


	return nil
}