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

	fmt.Printf(
		"Making a POST request with the initial login " +
			"credentials to the main site.\n")
	requestURL := "https://online-qa-test.ccdemo.site"
	var headers map[string]string
	// headers = {"Host": "CCCCCC"}
	bodyByteSlice := []byte(`{"client_message": "hemlo, server fren."}`)
	// Creating an io reader 'object' with our message in it."
	bodyReader := bytes.NewReader(bodyByteSlice)
	// bodyReader := myReader{}
	responseByteSlice, _ := myClient.MakeRequest(http.MethodPost, requestURL, bodyReader,
		true, "qa-test", "1z2a6iTzNmKPvHga", "", headers)
	csrfToken, err := myClient.ExtractCSRFToken(responseByteSlice)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("\n\n")

	fmt.Printf("Making a POST request to the main site.\n")
	requestURL = "https://online-qa-test.ccdemo.site/user"
	// bodyByteSlice = []byte(`name=civicrm_user&pass=civicrm_user&form_build_id=form-b6RPma1-XvAH27PP0_HdTD8hZ0GZi2V5o8yzXP33la8&form_id=user_login&op=Log+in`)
	bodyByteSlice = []byte(``)
	// Creating an io reader 'object' with our message in it."
	bodyReader = bytes.NewReader(bodyByteSlice)
	headers = map[string]string{
		"Host":                      "online-qa-test.ccdemo.site",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":           "en-US,en;q=0.5",
		"Accept-Encoding":           "gzip, deflate, br",
		"Content-Type":              "application/x-www-form-urlencoded",
		"Content-Length":            "127",
		"Origin":                    "https://online-qa-test.ccdemo.site",
		"Connection":                "keep-alive",
		"Referer":                   "https://online-qa-test.ccdemo.site/user",
		"Cookie":                    "has_js=1",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "same-origin",
		"Sec-Fetch-User":            "?1",
		"priority":                  "u=1",
		"name":                      "value",
	}
	myClient.MakeRequest(http.MethodPost, requestURL, bodyReader,
		true, "civicrm_user", "civicrm_user", csrfToken, headers)
	fmt.Print("\n\n")

	fmt.Printf("Making a GET request to the civiCRM system, with authentication.\n")
	requestURL = "https://online-qa-test.ccdemo.site/civicrm/dashboard"
	myClient.MakeRequest(http.MethodGet, requestURL,
		nil, true, "civicrm_user", "civicrm_user", "", nil)
	fmt.Print("\n\n")
}
