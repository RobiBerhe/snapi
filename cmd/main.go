package main

import (
	"log"
	"os"
	"snapi/internal"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("you need to pass the test file")
	}
	path := args[0]
	log.Println("file path :> ", path)
	tests := internal.NewTestSpecJson(path).ReadJSON()
	runTests := internal.Test(tests)
	runTests.Run()
}
