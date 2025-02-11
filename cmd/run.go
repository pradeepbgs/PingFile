package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/pradeepbgs/pingfile/src/config"
	"github.com/pradeepbgs/pingfile/src/runner"
	"github.com/spf13/cobra"
)

func execMultithreaded(filepath string, saveResponses bool, cookies []*http.Cookie) {
	var apiConfig, err = config.Parser(filepath)
	if err != nil {
		log.Printf("Error parsing file: %v", err)
		return
	}

	switch v := apiConfig.(type) {
	case *config.APIConfig:
		buffer, err := runner.ExecuteAPI(v, saveResponses, cookies, filepath)
		if err != nil {
			log.Printf("Request execution failed: %v", err)
			return
		}
		fmt.Print(buffer)
		
	case *config.GroupApiConfig:
		for i := range v.APIs {
			api_Config := &v.APIs[i]
			if !strings.HasPrefix(api_Config.URL, "http://") && !strings.HasPrefix(api_Config.URL, "https://") {
				api_Config.URL = v.BaseUrl + api_Config.URL
			}

			if api_Config.Run != nil && !*api_Config.Run {
				fmt.Printf("\nRunning config is disabled for %s file, skipping execution.\n", api_Config.URL)
				continue
			}

			buffer, err := runner.ExecuteAPI(api_Config, saveResponses, cookies, filepath)
			if err != nil {
				log.Printf("Request execution failed: %v", err)
				return
			}
			fmt.Print(buffer)
		}
	default:
		log.Println("Unknown config type")
	}
}

func execSequentially(filepath string, saveResponses bool, cookies []*http.Cookie) {
	var apiConfig, err = config.Parser(filepath)
	if err != nil {
		log.Printf("Error parsing file: %v", err)
		return
	}

	switch v := apiConfig.(type) {
	case *config.APIConfig:
		buffer, err := runner.ExecuteAPI(v, saveResponses, cookies, filepath)
		if err != nil {
			log.Printf("Request execution failed: %v", err)
			return
		}
		fmt.Print(buffer)
	case *config.GroupApiConfig:
		for i := range v.APIs {
			api_Config := &v.APIs[i]

			if !strings.HasPrefix(api_Config.URL, "http://") && !strings.HasPrefix(api_Config.URL, "https://") {
				api_Config.URL = v.BaseUrl + api_Config.URL
			}

			if api_Config.Run != nil && !*api_Config.Run {
				fmt.Printf("\nRunning config is disabled for %s file, skipping execution.\n", api_Config.URL)
				continue
			}

			buffer, err := runner.ExecuteAPI(api_Config, saveResponses, cookies, filepath)
			if err != nil {
				log.Printf("Request execution failed: %v", err)
				return
			}
			fmt.Print(buffer)
		}
	default:
		log.Println("Unknown config type")
	}
}

var runCmd = &cobra.Command{
	Use:   "run [files]",
	Short: "Execute API requests from a file",
	Long:  "The run command executes API requests defined in JSON, YAML, or PKFILE formats.",

	Run: func(cmd *cobra.Command, args []string) {
		runtime.GOMAXPROCS(runtime.NumCPU())
		filepaths := args

		greenColor := color.New(color.FgGreen).SprintFunc()
		BlueColor := color.New(color.FgCyan).SprintFunc()
		fmt.Println(BlueColor("--------------- >>>>"))
		fmt.Println(greenColor("Running PingFile "))
		fmt.Println(BlueColor("<<<<---------------"))

		saveResponses, _ := cmd.Flags().GetBool("save")

		var multiThread bool
		if cmd.Flags().Changed("multithread") {
			multiThread, _ = cmd.Flags().GetBool("multithread")
		}

		cookies, err := config.ParseCookie("root.cookie.pkfile")
		if err != nil {
			log.Printf("Error parsing cookies: %v", err)
			cookies = nil
		}

		if multiThread {
			workers, _ := cmd.Flags().GetInt("workers")
			if workers <= 0 {
				workers = 4
			}
			var wg sync.WaitGroup
			jobQueue := make(chan string, len(filepaths))

			for i := 0; i < workers; i++ {
				go func() {
					for file := range jobQueue {
						execMultithreaded(file, saveResponses, cookies)
						wg.Done()
					}
				}()
			}

			for _, filepath := range filepaths {
				wg.Add(1)
				jobQueue <- filepath
			}

			close(jobQueue)
			wg.Wait()

		} else {
			for _, filepath := range filepaths {
				execSequentially(filepath, saveResponses, cookies)
			}
		}
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the PingFile binary to a directory in your PATH",
	Run: func(cmd *cobra.Command, args []string) {
		installBinary()
	},
}

func installBinary() {
	binaryPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}

	destDir := ""
	if runtime.GOOS == "linux" {
		destDir = "/usr/local/bin"
	} else if runtime.GOOS == "windows" {
		destDir = filepath.Join(os.Getenv("USERPROFILE"), "bin")
		if err := os.MkdirAll(destDir, 0755); err != nil {
			log.Fatalf("Error creating directory %s: %v", destDir, err)
		}
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "macos" {
		destDir = filepath.Join(os.Getenv("HOME"), "bin")
		if err := os.MkdirAll(destDir, 0755); err != nil {
			log.Fatalf("Error creating directory %s: %v", destDir, err)
		}
	}

	if destDir == "" {
		log.Fatal("Unsupported os")
	}

	destPath := filepath.Join(destDir, "pingfile")
	err = os.Rename(binaryPath, destPath)
	if err != nil {
		log.Fatal("Error installing Binary", err)
	}
	fmt.Printf("PingFile installed to %s\n", destPath)
	fmt.Println("Make sure the directory is in your PATH.")
}

func init() {
	runCmd.Flags().BoolP("multithread", "m", false, "Run API requests concurrently (multi-threaded)")
	runCmd.Flags().BoolP("save", "s", false, "Save API response and request details")
	runCmd.Flags().IntP("workers", "w", 4, "Number of concurrent workers (default: 4)")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(installCmd)
}
