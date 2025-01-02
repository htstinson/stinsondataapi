package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"fmt"
	"log"
	"net/http"
	"net/url"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
	auth "github.com/htstinson/stinsondataapi/api/salesforce/auth"
	"github.com/htstinson/stinsondataapi/api/salesforce/model"

	"github.com/gorilla/mux"
)

type SalesforceHandler struct {
	Auth   *auth.SalesforceAuth
	logger *log.Logger
}

func New(creds *auth.SalesforceCreds) (*SalesforceHandler, error) {

	var SalesforceHandler = &SalesforceHandler{}

	authResponse, err := auth.SalesForceLogin(creds)
	if err != nil {
		fmt.Println(err.Error())
		return SalesforceHandler, err
	}

	SalesforceHandler.Auth = &auth.SalesforceAuth{
		AccessToken: authResponse.AccessToken,
		InstanceURL: creds.InstanceURL,
	}

	SalesforceHandler.logger = log.New(os.Stdout, "[API] ", log.LstdFlags)

	return SalesforceHandler, err
}

func (h *SalesforceHandler) Get(endpoint string, query string) ([]byte, error) {

	// Construct full URL
	baseurl := h.Auth.InstanceURL + endpoint

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
	req.Header.Set("Authorization", "Bearer "+h.Auth.AccessToken)
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

func (h *SalesforceHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {

	query := `SELECT Id, Name, Industry, Description, Phone, Fax, Website, LastModifiedDate, CreatedDate, LastActivityDate,	LastViewedDate, IsDeleted, MasterRecordId, Type, ParentId, BillingStreet, BillingCity, BillingState, BillingPostalCode, BillingCountry, AnnualRevenue, NumberOfEmployees, OwnerId, CreatedById, LastModifiedById, AccountSource FROM Account ORDER BY Name LIMIT 200`

	data, err := h.Get("/services/data/v59.0/query?q=", query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	response := model.AccountQueryResponse{}

	err = json.Unmarshal(data, &response)
	if err != nil {
		h.logger.Println(err.Error())
	}

	common.RespondJSON(w, http.StatusOK, response.Records)
}

func (h *SalesforceHandler) GetAccount(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	account, err := h.GetAccountById(id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "account lookup error")
		return
	}

	if account.Id == "" {
		common.RespondError(w, http.StatusNotFound, "account not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, account)
}

func (h *SalesforceHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	h.logger.Println(id)

	var Account model.Account // this is for new or updated accounts

	if err := json.NewDecoder(r.Body).Decode(&Account); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	currentAccount, err := h.GetAccountById(id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if currentAccount.Id == "" {
		common.RespondError(w, http.StatusNotFound, "Account not found")
		return
	}

	endpoint := fmt.Sprintf("/services/data/v59.0/sobjects/Account/%s", currentAccount.Id)

	transformedAccount := model.Transform(Account)

	_, err = h.SalesforcePatch(endpoint, transformedAccount)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error updating account")
		return
	}

	common.RespondJSON(w, http.StatusOK, Account)
}

func (h *SalesforceHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	var bodyBytes bytes.Buffer
	_, err := bodyBytes.ReadFrom(r.Body)
	if err != nil {
		h.logger.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	// Display the body
	h.logger.Printf("Body: %s", bodyBytes.String())

	// Restore the body for further processing
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes.Bytes()))

	var account *model.NewAccount
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	response, err := h.SalesforcePost("/services/data/v62.0/sobjects/Account", account)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	h.logger.Println(string(response))

	common.RespondJSON(w, http.StatusOK, "test complete")
}

func (h *SalesforceHandler) GetAccountById(id string) (model.Account, error) {

	query := fmt.Sprintf(`SELECT Id, Name, Industry, Description, Phone, Fax, Website, LastModifiedDate, CreatedDate, LastActivityDate,	LastViewedDate, IsDeleted, MasterRecordId, Type, ParentId, BillingStreet, BillingCity, BillingState, BillingPostalCode, BillingCountry, AnnualRevenue, NumberOfEmployees, OwnerId, CreatedById, LastModifiedById, AccountSource FROM Account Where Id = '%s' LIMIT 200`, id)

	data, err := h.Get("/services/data/v59.0/query?q=", query)
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

func (h *SalesforceHandler) SalesforcePatch(endpoint string, payload interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Construct full URL
	url := h.Auth.InstanceURL + endpoint

	// Create request
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+h.Auth.AccessToken)
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

func (h *SalesforceHandler) SalesforcePost(endpoint string, payload interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Construct full URL
	url := h.Auth.InstanceURL + endpoint

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+h.Auth.AccessToken)
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
