package main

import (
	"encoding/json"
	"log"
	"os"
)

type Test struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		User    string `yaml:"user"`
		Passwod string `yaml:"Passwod"`
		Name    string `yaml:"Name"`
	}
}

type Api struct {
	Tests struct {
		BaseURL string `yaml:"baseURL"`

		Api map[string]interface{}
	}
}

// Api     []struct {
// 	Name   string `yaml:"name"`
// 	Route  string `yaml:"route"`
// 	Method string `yaml:"method"`
// }

type T struct {
	Tests struct {
		BaseURL string `json:"base_url"`
		Apis    []struct {
			Name    string      `json:"name"`
			Method  string      `json:"method"`
			Route   string      `json:"route"`
			Payload interface{} `json:"payload"`
		} `json:"apis"`
	} `json:"tests"`
}

func main() {

	file, err := os.Open("test.json")
	if err != nil {
		log.Fatalf("error opening YAML file:%v", err)
	}
	defer file.Close()

	// var test T
	var test T
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&test)
	if err != nil {
		log.Fatalf("error decoding json file%v", err)
	}

	// use the config values

	log.Println("test out base url:> ", test.Tests.BaseURL)
	log.Println("test out apis:> ", test.Tests.Apis[0].Method)
	log.Println("test out apis:> ", test.Tests.Apis[1].Method)

}
