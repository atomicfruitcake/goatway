package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Set the URL of a downstream client
var clientHost = "https://127.0.0.1:8080/"
var adminToken = "exampleToken"

// Set the names of all downstream services
var validServices = map[string]bool {
	"exampleServiceA": true,
    "exampleServiceB": true,
}

func SendJob(service string, body []byte) error {
	if !validServices[service] {
		return error(fmt.Errorf("Service %s is not a valid Service", service))
	}
	url := fmt.Sprintf(clientHost + service)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
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
