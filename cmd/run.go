package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/pradeepbgs/pingfile/internal/config"
	"github.com/pradeepbgs/pingfile/internal/runner"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Execute API requests from a file",
	Long: `The run command executes API requests defined in JSON, YAML, or PKFILE formats.`,
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]

		fmt.Println("--------------- >>>>")
		fmt.Printf("Running PingFile for: %s\n", filepath)
		fmt.Println("<<<<---------------")
		
		var apiConfig , err = config.Parser(filepath)
		if err != nil {
			log.Fatalf("Error parsing file: %v", err)
			os.Exit(1)
		}

		if err != nil {
			log.Fatalf("Error parsing file: %v", err)
		}

		err = runner.ExecuteAPI(apiConfig)
		if err != nil {
			log.Fatalf("Request execution failed: %v", err)
		}

		fmt.Println("API request executed successfully!")
	},
}


func init() {
	rootCmd.AddCommand(runCmd)
}