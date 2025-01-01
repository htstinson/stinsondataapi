package auth

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SalesforceCreds struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	AccessToken  string `json:"accessToken"`
	InstanceURL  string `json:"instanceURL"`
	LoginURL     string `json:"loginURL"`
}

type SalesforceAuth struct {
	AccessToken string
	InstanceURL string
}

type SalesforceAuthResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	TokenType   string `json:"token_type"`
}

func SalesForceLogin(SalesforceCreds *SalesforceCreds) (*SalesforceAuthResponse, error) {
	// Your Salesforce credentials
	var (
		clientID     = SalesforceCreds.ClientId
		clientSecret = SalesforceCreds.ClientSecret
		username     = SalesforceCreds.Username
		password     = SalesforceCreds.Password
		accessToken  = SalesforceCreds.AccessToken
		loginURL     = SalesforceCreds.LoginURL
	)

	password += accessToken

	var auth *SalesforceAuthResponse
	auth, err := GetSalesforceToken(clientID, clientSecret, username, password, loginURL)
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
	}

	return auth, err
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
