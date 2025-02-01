package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "pingfile",
	Short: "PingFile CLI to execute API requests from configuration files",
	Long:  `PingFile CLI helps execute API requests defined in JSON, YAML, or PKFILE formats.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to PingFile!")
		fmt.Println("Use 'pingfile run <file>' to execute API requests from a file.")
	},
	Version: version,
}

func init() {
	rootCmd.SetVersionTemplate(`{{printf "%s version %s\n" .Name .Version}}`)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
