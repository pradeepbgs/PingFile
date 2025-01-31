package main

import (

	"github.com/joho/godotenv"
	"github.com/pradeepbgs/pingfile/cmd"
)

func main(){
	godotenv.Load()
	// if err != nil {
	// 	log.Println("No .env file found or error loading .env file")
	// }
	cmd.Execute()
}