package main

import (
	"log"
	"os"
	"snapi/internal/snapi"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("you need to pass the test file")
	}
	path := args[0]
	log.Println("file path :> ", path)
	tests := snapi.NewTestSpecJson(path).ReadJSON()
	runTests := snapi.Test(tests)
	runTests.Run()
}
