package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

// Set the URL of a downstream client
var clientHost = "https://127.0.0.1:8080/"
var adminToken = "exampleToken"

// Set the names of all downstream services
var validServices = map[string]bool {
	"exampleServiceA": true,
    "exampleServiceB": true,
}

func sendJob(service string, body []byte) error {
	if !validServices[service] {
		return error(fmt.Errorf("Service %s is not a valid Service", service))
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf(clientHost + service), bytes.NewBuffer(body))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("admintoken", adminToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return error(fmt.Errorf("Error sending request to %s due to: %v", url, err))
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return error(fmt.Errorf("Error reading response from %s due to %v", url, err))
	}
	return err
}
