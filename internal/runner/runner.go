package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/fatih/color"
	"github.com/pradeepbgs/pingfile/internal/config"
)



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
	var req *http.Request
	var err error

	hasFile := len(apiConfig.File) > 0
	
	if hasFile {
		var requestBody bytes.Buffer
		// body := &bytes.Buffer{}
		writer := multipart.NewWriter(&requestBody)

		for key, value := range apiConfig.Body {
			err := writer.WriteField(key,fmt.Sprintf("%v",value))
			if err != nil{
				return fmt.Errorf("failed to write field: %w", err)
			}
		}	
		// attach files
		for _,fileItem := range apiConfig.File{
			file,err := os.Open(fileItem.Path)
			if err != nil{
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			part, err := writer.CreateFormFile(fileItem.Name, fileItem.Path)
			if err != nil{
				return fmt.Errorf("failed to create form file: %w", err)
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return fmt.Errorf("failed to copy file content: %w", err)
			}
		}

		if err := writer.Close(); err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}

		req, err = http.NewRequest(apiConfig.Headers["Method"], apiConfig.URL, &requestBody)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

	} else {
		// esle ping normal json  request	
		req, err = http.NewRequest(apiConfig.Headers["Method"], apiConfig.URL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
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

	// Add other headers to the request
	for key, value := range apiConfig.Headers {
		if key != "Method" {
			req.Header.Set(key, value)
		}
	}

	// add cookies if enabled
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

		err := config.SaveResponseToFile(saveFilePath, requestDetails, responseDetails)
		if err != nil {
			return fmt.Errorf("\nfailed to save response: %w", err)
		}
		fmt.Println("\nResponse saved to response.json")
	}

	cookies := resp.Cookies()
	if len(cookies) > 0 {
		config.SaveCookies("root.cookie.pkfile", cookies)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s: %s", errorColor("Error"), resp.Status)
	}

	return nil
}
