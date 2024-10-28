package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type API struct {
	Name    string      `json:"name"`
	Method  string      `json:"method"`
	Route   string      `json:"route"`
	Payload interface{} `json:"payload"`
	Expects interface{} `json:"expects"`
}

type TestSpecJSON struct {
	Tests struct {
		BaseURL string `json:"base_url"`
		Apis    []API  `json:"apis"`
	} `json:"tests"`
	file *os.File
}

func NewTestSpecJson(filePath string) *TestSpecJSON {
	ts := &TestSpecJSON{}
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("couldn't open file", err)
	}
	ts.file = file
	return ts
}

func (ts *TestSpecJSON) ReadJSON() *TestSpecJSON {
	// var spec TestSpecJSON
	decoder := json.NewDecoder(ts.file)
	if err := decoder.Decode(&ts); err != nil {
		log.Fatalf("error decoding json file:%v", err)
	}

	return ts
}

type Fails struct {
	Name string
	Msg  string
}

type TestAPI struct {
	TestSpecJSON
}

func Test(ts *TestSpecJSON) *TestAPI {
	return &TestAPI{TestSpecJSON: *ts}
}

func (ta *TestAPI) call(api *API) []byte {
	base := ta.TestSpecJSON.Tests.BaseURL
	var req *http.Request
	var data []byte
	var err error

	apiMethod := strings.ToUpper(api.Method)
	if apiMethod == "POST" || apiMethod == "PUT" || apiMethod == "PATCH" {
		data, err = json.Marshal(api.Payload)
		if err != nil {
			log.Fatalf("couldn't marshal values to json %v", err)
		}
		log.Println("data : ", data)
	} else if strings.ToUpper(api.Method) == "GET" {
		data = nil
	}
	req, err = http.NewRequest(strings.ToUpper(api.Method), fmt.Sprintf("%v%v", base, api.Route), bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("error creating a request object for %v, fails with error %v", api.Name, err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("request failed :%v", err)
	}
	defer res.Body.Close()
	log.Println("status code : ", res.StatusCode)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading response body %v", err)
	}
	// log.Printf("response body : %v", string(body))
	return body
}

func (ta *TestAPI) passExpects(body []byte, api *API) error {
	// log.Println("the api body is expected to be : ", api.Expects)
	exData, exerr := json.Marshal(api.Expects)
	if exerr != nil {
		return exerr
	}

	var j1, j2 interface{}
	if err := json.Unmarshal(exData, &j1); err != nil {
		log.Fatal("error :", err)
	}
	if err := json.Unmarshal(body, &j2); err != nil {
		log.Fatal("error :", err)
	}
	// log.Println("j1 :> ", j1, " and j2 :> ", j2)
	if reflect.DeepEqual(j1, j2) {
		log.Println("EQUAL..")
		return nil
	}
	diff := cmp.Diff(j1, j2)
	return errors.New(diff)
}

func (ta *TestAPI) Run() {
	apis := ta.TestSpecJSON.Tests.Apis
	for _, api := range apis {
		log.Println("Name : ", api.Name)
		log.Println("Method : ", api.Method)
		log.Println("Route : ", api.Route)
		response := ta.call(&api)
		log.Printf("the response for %v is :>> %v ", api.Name, string(response))
		err := ta.passExpects(response, &api)
		if err != nil {
			log.Fatalf("API test fails at:%v with reason :%v\n", api.Name, err.Error())
		}
	}
}
