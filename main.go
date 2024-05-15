package main

import (
	"back_end_tester/requestHandlers"
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func main() {
	myClient := requestHandlers.MyClient{
		Port: 3333,
		Body: http.Client{
			Timeout: 30 * time.Second,
		},
	}
	fmt.Printf("Making a GET request - Test 1: cat facts.\n")
	requestURL := fmt.Sprintf(
		"https://cat-fact.herokuapp.com/facts?animal_type=cat")
	myClient.MakeRequest(http.MethodGet, requestURL, nil)
	fmt.Println("\n")

	fmt.Printf(
		"Making a GET request - Test 2: my user without authentication.\n")
	requestURL = "https://cat-fact.herokuapp.com/users/me"
	myClient.MakeRequest(http.MethodGet, requestURL, nil)
	fmt.Println("\n")

	fmt.Printf("Making a POST request - Test 5\n.")
	requestURL = "https://cat-fact.herokuapp.com/me"
	jsonByteSlice := []byte(`{"client_message": "hemlo, server fren."}`)
	// Creating an io reader 'object' with our message in it."
	bodyReader := bytes.NewReader(jsonByteSlice)
	myClient.MakeRequest(http.MethodPost, requestURL, bodyReader)
}
