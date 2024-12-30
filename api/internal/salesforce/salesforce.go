package salesforce

import (
	"api/internal/model"
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
	fmt.Println("SalesforcePost")

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	fmt.Println(string(jsonData))

	// Construct full URL
	url := auth.InstanceURL + endpoint

	fmt.Println("url", url)

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

func GetAccount(auth SalesforceAuth, id string) (model.Account, error) {

	fmt.Println("GetAccount")

	query := fmt.Sprintf(`SELECT Id, Name, Industry, Description, Phone, Fax, Website, LastModifiedDate, CreatedDate, LastActivityDate,	LastViewedDate, IsDeleted, MasterRecordId, Type, ParentId, BillingStreet, BillingCity, BillingState, BillingPostalCode, BillingCountry, AnnualRevenue, NumberOfEmployees, OwnerId, CreatedById, LastModifiedById, AccountSource FROM Account Where Id = '%s' LIMIT 200`, id)

	fmt.Println(query)
	fmt.Println()

	data, err := SalesforceGet(auth, "/services/data/v59.0/query?q=", query, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return model.Account{}, err
	}

	response := model.AccountQueryResponse{}

	err = json.Unmarshal(data, &response)
	if err != nil {
		// Handle error
	}

	return response.Records[0], nil

}

// {{_endpoint}}/services/data/v{{version}}/sobjects/:SOBJECT_API_NAME/:RECORD_ID

func SalesforcePatch(auth SalesforceAuth, endpoint string, payload interface{}) ([]byte, error) {
	fmt.Println("SalesforcePatch")

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	fmt.Println(string(jsonData))

	// Construct full URL
	url := auth.InstanceURL + endpoint

	fmt.Println("url", url)

	// Create request
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
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

func UpdateAccount(auth SalesforceAuth, newAccount model.NewAccount) error {

	fmt.Println("salesforce UpdateAccount")

	return nil
}

func SalesForceLogin(SalesforceCreds model.SalesforceCreds) (*SalesforceAuthResponse, error) {
	// Your Salesforce credentials
	var (
		clientID     = SalesforceCreds.ClientId
		clientSecret = SalesforceCreds.ClientSecret
		username     = SalesforceCreds.Username
		password     = SalesforceCreds.Password
		accessToken  = SalesforceCreds.AccessToken
		loginURL     = "https://login.salesforce.com"
	)

	password += accessToken

	var auth *SalesforceAuthResponse
	auth, err := GetSalesforceToken(clientID, clientSecret, username, password, loginURL)
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		fmt.Println(clientID, clientSecret, username, password, loginURL)
	}

	return auth, err

}
