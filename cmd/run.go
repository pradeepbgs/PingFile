package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/pradeepbgs/pingfile/internal/config"
	"github.com/pradeepbgs/pingfile/internal/runner"
	"github.com/spf13/cobra"
)

func exec(filepath string, wg *sync.WaitGroup) {
	var apiConfig, err = config.Parser(filepath)
	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
		os.Exit(1)
	}

	err = runner.ExecuteAPI(apiConfig)
	if err != nil {
		log.Fatalf("Request execution failed: %v", err)
	}

	fmt.Println("API request executed successfully for:", filepath)
	defer wg.Done()
}

func execSequentially(filepath string) {
	var apiConfig, err = config.Parser(filepath)
	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
		os.Exit(1)
	}

	err = runner.ExecuteAPI(apiConfig)
	if err != nil {
		log.Fatalf("Request execution failed: %v", err)
	}

	fmt.Println("API request executed successfully for:", filepath)
}

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Execute API requests from a file",
	Long:  `The run command executes API requests defined in JSON, YAML, or PKFILE formats.`,
	Run: func(cmd *cobra.Command, args []string) {
		filepaths := args
		multiThread, _ := cmd.Flags().GetBool("multithread")

		if multiThread {
			var wg sync.WaitGroup
			for _, filepath := range filepaths {
				fmt.Println("--------------- >>>>")
					fmt.Printf("Running PingFile for: %s\n", filepath)
					fmt.Println("<<<<---------------")
				wg.Add(1)
				go func(file string) {					
					exec(file, &wg)
				}(filepath)
			}
			
			wg.Wait()
			
		} else {

		

			for _, filepath := range filepaths {
				fmt.Println("--------------- >>>>")
				fmt.Printf("Running PingFile for: %s\n", filepath)
				fmt.Println("<<<<---------------")
				
				execSequentially(filepath)
				
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
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(installCmd)
}
