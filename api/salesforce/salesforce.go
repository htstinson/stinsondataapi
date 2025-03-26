package salesforce

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/htstinson/stinsondataapi/api/salesforce/auth"
	"github.com/htstinson/stinsondataapi/api/salesforce/handler"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
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
		fmt.Printf("[%v] Error: Salesforce Creds %s\n", time.Now().Format(time.RFC3339), err.Error())
		return salesforce, err
	}
	json.Unmarshal(salesforceCreds, &SalesforceCreds)

	salesforce.Creds = SalesforceCreds

	fmt.Println(salesforce.Creds.AccessToken)
	fmt.Println(salesforce.Creds.ClientId)
	fmt.Println(salesforce.Creds.ClientSecret)
	fmt.Println(salesforce.Creds.InstanceURL)
	fmt.Println(salesforce.Creds.LoginURL)
	fmt.Println(salesforce.Creds.Password)
	fmt.Println(salesforce.Creds.Username)

	handler, err := handler.New(SalesforceCreds)
	if err != nil {
		fmt.Printf("[%v] Error: %s\n", time.Now().Format(time.RFC3339), err.Error())
		return salesforce, err
	}

	salesforce.Handler = handler
	salesforce.logger = logger

	return salesforce, nil
}
