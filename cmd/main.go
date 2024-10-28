package main

import (
	"log"
	"os"
	"snapi/internal"
)

func main() {
	args := os.Args[1:]
	log.Println("os args :> ", args)
	path := args[0]
	log.Println("file path :> ", path)
	tests := internal.NewTestSpecJson(path).ReadJSON()
	runTests := internal.Test(tests)
	runTests.Run()
}
