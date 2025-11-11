package api

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"
)

// AuthClient is a reusable HTTP client with basic authentication
var AuthClient *http.Client

func init() {
	// Create a client with custom transport to skip TLS verification if needed
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Use only for testing
	}

	AuthClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

// NewAuthenticatedRequest creates a new HTTP request with basic auth headers
func NewAuthenticatedRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Get credentials from environment variables
	username := os.Getenv("API_USERNAME")
	password := os.Getenv("API_PASS")

	// Set basic authentication
	req.SetBasicAuth(username, password)
	return req, nil
}
