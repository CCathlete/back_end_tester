package requestHandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
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

func (c *MyClient) MakeRequest(
	method string, requestURL string,
	requestBody io.Reader, withAUTH bool,
	userName, password, csrfToken string,
	headers map[string]string) ([]byte, map[string][]string) {
	// Checking the method we use.
	if isMethodValid := ValidateHttpMethod(method); !isMethodValid {
		fmt.Println("This type of request is not supported.")
		os.Exit(1)
	}
	// If there's a csrf token, it means that the request body
	// should be of type "application/x-www-form-urlencoded".
	if csrfToken != "" {
		form := url.Values{}
		form.Add("name", userName)
		form.Add("pass", password)
		form.Add("form_build_id", csrfToken)
		form.Add("form_id", "user_login")
		form.Add("op", `Log+in`)
		requestBody = strings.NewReader(form.Encode())
	}
	// Creating an http request 'object' with the URL inside.
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
	// request.Header.Set("Content-Type", "application/json")

	// Setting additional headers if needed.
	for headerName, headerContent := range headers {
		request.Header.Set(headerName, headerContent)
	}

	// Activating the request.
	response, err := c.Body.Do(request)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// Checking if the response was a redirect.
	if response.StatusCode >= 300 && response.StatusCode <= 399 {
		fmt.Println("We had a redirect! ", response.StatusCode)
		redirectUrl, err := response.Location()
		if err != nil {
			fmt.Println("Error getting redirect location: ",
				err)
			os.Exit(1)
		}

		// Creating a new GET request to follow the redirect.
		request.URL = redirectUrl
		response, err = c.Body.Do(request)
		if err != nil {
			fmt.Println("Error sending redirect request: ",
				err)
			os.Exit(1)
		}
	}
	defer response.Body.Close()

	// Processing the response.
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
	return responseByteSlice, response.Header
}

func (c *MyClient) ExtractCSRFToken(responseByteSlice []byte) (string, error) {
	// Gets the body of a response as byte slice, converts it to string,
	// extracts from it the csrf token using regexp and returns it.
	// I assume that the response if in HTML form since it's after a GET request.
	responseBodyString := string(responseByteSlice)
	pattern := `name="form_build_id" value="(.+?)"`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(responseBodyString)
	if len(matches) < 2 {
		return "", fmt.Errorf("CSRF token not found")
	}
	return matches[1], nil
}
