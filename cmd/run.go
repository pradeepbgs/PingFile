package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/pradeepbgs/pingfile/internal/config"
	"github.com/pradeepbgs/pingfile/internal/runner"
	"github.com/spf13/cobra"
)

func exec(filepath string, wg *sync.WaitGroup,saveResponses bool,cookies []*http.Cookie) {
	fmt.Println("--------------- >>>>")
	fmt.Printf("Running PingFile for: %s\n", filepath)
	fmt.Println("<<<<---------------")
	
	var apiConfig, err = config.Parser(filepath)

	if err != nil {
		log.Printf("Error parsing file: %v", err)
		defer wg.Done()
		return
	}

	err = runner.ExecuteAPI(apiConfig,saveResponses,cookies)
	if err != nil {
		log.Printf("Request execution failed: %v", err)
		wg.Done()
		return
	}

	fmt.Println("\nAPI request executed successfully for:", filepath)
	defer wg.Done()
}

func execSequentially(filepath string,saveResponses bool,cookies []*http.Cookie) {
	fmt.Println("--------------- >>>>")
	fmt.Printf("Running PingFile for: %s\n", filepath)
	fmt.Println("<<<<---------------")

	var apiConfig, err = config.Parser(filepath)
	if err != nil {
		log.Printf("Error parsing file: %v", err)
		return
	}

	err = runner.ExecuteAPI(apiConfig,saveResponses,cookies)
	if err != nil {
		log.Printf("Request execution failed: %v", err)
		return
	}

	fmt.Println("\nAPI request executed successfully for:", filepath)
}

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Execute API requests from a file",
	Long:  `The run command executes API requests defined in JSON, YAML, or PKFILE formats.`,
	
	Run: func(cmd *cobra.Command, args []string) {
		filepaths := args

		saveResponses , _ := cmd.Flags().GetBool("save")

		var multiThread bool
		if cmd.Flags().Changed("multithread") {
			multiThread, _ = cmd.Flags().GetBool("multithread")
		}	
		
		cookies, err := config.ParseCookie("cookie.pkfile")
		if err != nil {
			log.Printf("Error parsing cookies: %v", err)
			cookies = nil
		}

		if multiThread {
			var wg sync.WaitGroup
			for _, filepath := range filepaths {
				wg.Add(1)
				go func(file string) {
					exec(file, &wg, saveResponses, cookies)
				}(filepath)
			}

			wg.Wait()

		} else {
			for _, filepath := range filepaths {
				execSequentially(filepath,saveResponses,cookies)
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
		log.Fatal("Error installing Binary")
	}
	fmt.Printf("PingFile installed to %s\n", destPath)
	fmt.Println("Make sure the directory is in your PATH.")
}

func init() {
	runCmd.Flags().BoolP("multithread", "m", false, "Run API requests concurrently (multi-threaded)")
	runCmd.Flags().BoolP("save", "s", false, "Save API response and request details")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(installCmd)
}
