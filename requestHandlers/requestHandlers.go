package requestHandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Port int
	Body http.Server
}

type MyClient struct {
	Port int
	Body http.Client
}

func (c *MyClient) GetRequest(requestURL string) {
	response, err := http.Get(requestURL)
	// If there"s a problem, exit the program with return value = 1.
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	time.Sleep(50 * time.Millisecond) // Might not be needed, check.
	fmt.Println("client: got response! YEY!")
	fmt.Printf("client: status code %d\n", response.StatusCode)
	// Reading the response's body.
	responseByteSlice, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("client: couldn't read the response: %s\n", err)
		os.Exit(1)
	}
	// Print the json response with proper indentation.
	var byteBuffer bytes.Buffer
	json.Indent(&byteBuffer, responseByteSlice, "", "\t")
	if err != nil {
		fmt.Printf("couldn't pretty print json: %s\n", err)
	}
	fmt.Printf("\nclient: response body:\n")
	byteBuffer.WriteTo(os.Stdout)
}

func ValidateHttpMethod(method string) bool {
	switch method {
	case http.MethodGet:
		return true
	case http.MethodPost:
		return true
	default:
		return false
	}
}

func (c *MyClient) MakeRequest(method string, requestURL string,
	requestBody io.Reader, withAUTH bool, userName, password string,
	headers map[string]string) {
	// Checking the method we use.
	if isMethodValid := ValidateHttpMethod(method); !isMethodValid {
		fmt.Println("This type of request is not supported.")
		os.Exit(1)
	}
	// Create an http request 'object' with the URL inside.
	request, err := http.NewRequest(method, requestURL, requestBody)
	if withAUTH {
		request.SetBasicAuth(userName, password)
	}
	// If there"s a problem, exit the program with return value = 1.
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	/* Setting the content type of the request's body to be json media
	   type in our requests header. */
	request.Header.Set("Content-Type", "application/json")

	// Setting additional headers if needed.
	for headerName, headerContent := range headers {
		request.Header.Set(headerName, headerContent)
	}

	// Activate the request.
	response, err := c.Body.Do(request)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	time.Sleep(50 * time.Millisecond) // Might not be needed, check.
	fmt.Println("client: got response! YEY!")
	fmt.Printf("client: status code %d\n", response.StatusCode)
	// Reading the response's body.
	responseByteSlice, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("client: couldn't read the response: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println(string(responseByteSlice))
	}
	// Print the json response with proper indentation.
	var byteBuffer bytes.Buffer
	json.Indent(&byteBuffer, responseByteSlice, "", "\t")
	if err != nil {
		fmt.Printf("couldn't pretty print json: %s\n", err)
	}
	fmt.Printf("\nclient: response body:\n")
	byteBuffer.WriteTo(os.Stdout)
}
