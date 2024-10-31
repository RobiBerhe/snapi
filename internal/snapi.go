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

type Expects struct {
	StatusCode int         `json:"status"`
	Body       interface{} `json:"body"`
	Exclude    []string    `json:"exclude"`
}

type API struct {
	Name       string      `json:"name"`
	Method     string      `json:"method"`
	Route      string      `json:"route"`
	Payload    interface{} `json:"payload"`
	Expects    Expects     `json:"expects"`
	StatusCode int         `json:"status"`
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

func (ta *TestAPI) call(api *API) ([]byte, int) {
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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading response body %v", err)
	}
	return body, res.StatusCode
}

func (ta *TestAPI) passExpects(body []byte, api *API) error {
	exData, exerr := json.Marshal(api.Expects.Body)
	log.Println("excluded fields will be :> ", api.Expects.Exclude)
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
	if reflect.DeepEqual(j1, j2) {
		return nil
	}
	opts := cmp.FilterPath(func(p cmp.Path) bool {
		for _, field := range api.Expects.Exclude {
			if strings.Contains(p.GoString(), field) {
				return true
			}
		}
		return false
	}, cmp.Ignore())
	diff := cmp.Diff(j1, j2, opts)
	if diff != "" {
		return errors.New(diff)
	}
	return nil
}

func (ta *TestAPI) PassStatus(status int, apiStatus int) error {
	if apiStatus <= 0 {
		return nil
	}
	if status == apiStatus {
		return nil
	}
	return fmt.Errorf("expected status %d but found %v", apiStatus, status)
}

func (ta *TestAPI) Run() {
	apis := ta.TestSpecJSON.Tests.Apis
	log.Println("-----------TESTS RUNNING-----------")
	for _, api := range apis {
		log.Printf("Now checking :%v\n", api.Name)
		response, status := ta.call(&api)
		err := ta.passExpects(response, &api)
		if err != nil {
			log.Fatalf("API test fails at:'%v'  >>:%v\n", api.Name, err.Error())
		}
		if err := ta.PassStatus(status, api.Expects.StatusCode); err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Println("all tests have passed")
}
