package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/obriena/dockerdevtemplate/domain"
)

func TestPost(t *testing.T) {
	// URL of the resource to delete
	url := "http://localhost/posts/"

	post := domain.Post{
		Content:  "Have a good day!",
		Title:    "Greeting",
		Language: "English",
		OwnerId:  1,
	}

	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(post)
	if err != nil {
		log.Fatalf("Error creating JSON for request: %v", err)
	}
	// Create a new DELETE request
	req, err := http.NewRequest("POST", url, reqBodyBytes)
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

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	log.Println("Resource deleted successfully!")
}
