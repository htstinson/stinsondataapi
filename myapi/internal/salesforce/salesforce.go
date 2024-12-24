package salesforce

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SalesforceAuth struct {
	AccessToken string
	InstanceURL string
}

type SalesforceAuthResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	TokenType   string `json:"token_type"`
}

func SalesforceGet(auth SalesforceAuth, endpoint string, query string, payload interface{}) ([]byte, error) {

	// Construct full URL
	baseurl := auth.InstanceURL + endpoint

	// Create URL with encoded query parameter
	u, err := url.Parse(baseurl)
	if err != nil {
		// Handle error
	}

	q := u.Query()
	q.Set("q", query)
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func SalesforcePost(auth SalesforceAuth, endpoint string, payload interface{}) ([]byte, error) {
	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Construct full URL
	url := auth.InstanceURL + endpoint

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func GetSalesforceToken(clientID, clientSecret, username, password, loginURL string) (*SalesforceAuthResponse, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("username", username)
	data.Set("password", password)

	// Create request
	req, err := http.NewRequest(
		"POST",
		loginURL+"/services/oauth2/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create the TLS config
	config := &tls.Config{
		RootCAs: nil, // This will use system root certificates
	}

	// Create a custom transport with the TLS config
	transport := &http.Transport{
		TLSClientConfig: config,
	}

	// Create a client with the custom transport
	client := &http.Client{
		Transport: transport,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var authResponse SalesforceAuthResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &authResponse, nil
}
