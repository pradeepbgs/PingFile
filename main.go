package main

import (
	"fmt"

	"github.com/pradeepbgs/pingfile/internal/config"
	"github.com/pradeepbgs/pingfile/internal/runner"
)

func main(){
	file := "api.json"

	result,err := config.Parser(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Headers)
	runner.Execute(result)
}