package tests

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestRetrieveSinglePost(t *testing.T) {
	url := "http://localhost/posts/1"
	executeHttpGetCall(url)
}

func TestRetrievAllPosts(t *testing.T) {
	url := "http://localhost/posts"
	executeHttpGetCall(url)
}

func TestRetrieveOnwerPost(t *testing.T) {
	url := "http://localhost/posts/ownerId/1"
	executeHttpGetCall(url)
}

func executeHttpGetCall(urlString string) {
	// Create a new GET request
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		log.Fatalf("Error creating POST request: %v", err)
	}

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error performing POST request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	// Convert the byte slice to a string
	bodyString := string(bodyBytes)
	log.Println("Response is:\n", bodyString)
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	log.Println("Resource deleted successfully!")
}
