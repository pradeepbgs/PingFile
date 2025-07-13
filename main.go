package main

import (
	"github.com/joho/godotenv"
	"github.com/pradeepbgs/pingfile/cmd"
)

func main() {
	godotenv.Load()

	cmd.Execute()
}
