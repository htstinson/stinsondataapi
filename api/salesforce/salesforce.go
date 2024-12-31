package salesforce

import (
	"api/salesforce/auth"
	"api/salesforce/handler"
	"encoding/json"
	"log"
	"os"

	common "github.com/htstinson/stinsondataapi/api/internal/commonweb"
)

type Salesforce struct {
	Creds   *auth.SalesforceCreds
	Handler *handler.SalesforceHandler
	logger  *log.Logger
}

func New() (Salesforce, error) {

	var salesforce = Salesforce{}
	var logger = log.New(os.Stdout, "[API] ", log.LstdFlags)
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
	salesforce.logger = logger

	return salesforce, nil
}