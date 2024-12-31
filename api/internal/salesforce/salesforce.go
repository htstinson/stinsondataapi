package salesforce

import (
	common "api/internal/commonweb"
	"api/internal/salesforce/auth"
	"api/internal/salesforce/handler"
	"encoding/json"
	"log"
)

type Salesforce struct {
	Creds   *auth.SalesforceCreds
	Handler *handler.SalesforceHandler
}

func New(logger *log.Logger) (Salesforce, error) {

	var salesforce = Salesforce{}

	var SalesforceCreds = &auth.SalesforceCreds{}
	salesforceCreds, err := common.GetSecretString("Salesforce", "us-west-2")
	if err != nil {
		logger.Println("Salesforce Creds", err.Error())
		return salesforce, err
	}
	json.Unmarshal(salesforceCreds, &SalesforceCreds)

	salesforce.Creds = SalesforceCreds

	handler, err := handler.New(SalesforceCreds)
	if err != nil {
		logger.Println(err.Error())
		return salesforce, err
	}

	salesforce.Handler = handler

	return salesforce, nil
}
