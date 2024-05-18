package main

import (
	"back_end_tester/requestHandlers"
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type myReader struct {
	num int
}

func (mr myReader) Read(p []byte) (int, error) {
	return 42, nil
}

func main() {
	myClient := requestHandlers.MyClient{
		Port: 3333,
		Body: http.Client{
			Timeout: 30 * time.Second,
		},
	}
	// fmt.Printf("Making a GET request to the main site.\n")
	// requestURL := "https://online-qa-test.ccdemo.site"
	// myClient.MakeRequest(http.MethodGet, requestURL,
	// 	nil, true, "qa-test", "1z2a6iTzNmKPvHga", nil)
	// fmt.Println("\n")

	fmt.Printf("Making a POST request to the main site.\n")
	requestURL := "https://online-qa-test.ccdemo.site"
	var headers map[string]string
	// headers = {"Host": "CCCCCC"}
	jsonByteSlice := []byte(`{"client_message": "hemlo, server fren."}`)
	// Creating an io reader 'object' with our message in it."
	bodyReader := bytes.NewReader(jsonByteSlice)
	// bodyReader := myReader{}
	myClient.MakeRequest(http.MethodPost, requestURL, bodyReader,
		true, "qa-test", "1z2a6iTzNmKPvHga", headers)
	fmt.Println("\n")

	fmt.Printf("Making a GET request to the civiCRM system, with authentication.\n")
	requestURL = "https://online-qa-test.ccdemo.site/civicrm/dashboard"
	myClient.MakeRequest(http.MethodGet, requestURL,
		nil, true, "civicrm_user", "civicrm_user", nil)
	fmt.Println("\n")
}
