package salesforce

import (
	common "api/internal/commonweb"
	"api/internal/salesforce/auth"
	"encoding/json"
	"log"
)

type Salesforce struct {
	Creds  *auth.SalesforceCreds
	logger *log.Logger
}

func New(logger *log.Logger) (Salesforce, error) {

	var salesforce = Salesforce{
		logger: logger,
	}

	var SalesforceCreds = &auth.SalesforceCreds{}
	salesforceCreds, err := common.GetSecretString("Salesforce", "us-west-2")
	if err != nil {
		logger.Println("Salesforce Creds", err.Error())
		return salesforce, err
	}
	json.Unmarshal(salesforceCreds, &SalesforceCreds)

	salesforce.Creds = SalesforceCreds

	return salesforce, nil
}
