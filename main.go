package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/pradeepbgs/pingfile/cmd"
)

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file")
	}
	cmd.Execute()
}