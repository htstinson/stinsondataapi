package model

import sftime "github.com/htstinson/stinsondataapi/api/salesforce/time"

func Transform(full Account) NewAccount {

	transformed := NewAccount{
		Attributes:        full.Attributes,
		Name:              full.Name,
		Industry:          full.Industry,
		Description:       full.Description,
		Phone:             full.Phone,
		Fax:               full.Fax,
		Website:           full.Website,
		LastActivityDate:  full.LastActivityDate,
		MasterRecordId:    full.MasterRecordId,
		AccountType:       full.AccountType,
		ParentId:          full.ParentId,
		BillingCity:       full.BillingCity,
		BillingState:      full.BillingState,
		BillingPostalCode: full.BillingPostalCode,
		BillingCountry:    full.BillingCountry,
		AnnualRevenue:     full.AnnualRevenue,
		NumberOfEmployees: full.NumberOfEmployees,
		OwnerId:           full.OwnerId,
		AccountSource:     full.AccountSource,
	}

	return transformed
}

type AccountQueryResponse struct {
	TotalSize int       `json:"totalSize"`
	Done      bool      `json:"done"`
	Records   []Account `json:"records"`
}

type Account struct {
	Attributes        AccountAttributes      `json:"attributes"`
	Id                string                 `json:"Id"`
	Name              string                 `json:"Name,omitempty"`
	Industry          *string                `json:"Industry,omitempty"` // Using pointer since it can be null
	Description       string                 `json:"Description,omitempty"`
	Phone             string                 `json:"Phone,omitempty"`
	Fax               *string                `json:"Fax,omitempty"`
	Website           *string                `json:"Website,omitempty"`
	LastModifiedDate  *sftime.SalesforceTime `json:"LastModifiedDate,omitempty"`
	CreatedDate       *sftime.SalesforceTime `json:"CreatedDate,omitempty"`
	LastActivityDate  *sftime.SalesforceTime `json:"LastActivityDate,omitempty"`
	LastViewedDate    *sftime.SalesforceTime `json:"LastViewedDate,omitempty"`
	IsDeleted         *bool                  `json:"IsDeleted,omitempty"`
	MasterRecordId    *string                `json:"MasterRecordId,omitempty"`
	AccountType       *string                `json:"Type,omitempty"`
	ParentId          *string                `json:"ParentId,omitempty"`
	BillingStreet     *string                `json:"BillingStreet,omitempty"`
	BillingCity       *string                `json:"BillingCity,omitempty"`
	BillingState      *string                `json:"BillingState,omitempty"`
	BillingPostalCode *string                `json:"BillingPostalCode,omitempty"`
	BillingCountry    *string                `json:"BillingCountry,omitempty"`
	AnnualRevenue     *int64                 `json:"AnnualRevenue,omitempty"`
	NumberOfEmployees *int16                 `json:"NumberOfEmployees,omitempty"`
	OwnerId           *string                `json:"OwnerId,omitempty"`
	CreatedById       *string                `json:"CreatedById,omitempty"`
	LastModifiedById  *string                `json:"LastModifiedById,omitempty"`
	AccountSource     *string                `json:"AccountSource,omitempty"`
}

type NewAccount struct {
	Attributes        AccountAttributes      `json:"attributes"`
	Name              string                 `json:"Name,omitempty"`
	Industry          *string                `json:"Industry,omitempty"` // Using pointer since it can be null
	Description       string                 `json:"Description,omitempty"`
	Phone             string                 `json:"Phone,omitempty"`
	Fax               *string                `json:"Fax,omitempty"`
	Website           *string                `json:"Website,omitempty"`
	LastActivityDate  *sftime.SalesforceTime `json:"LastActivityDate,omitempty"`
	MasterRecordId    *string                `json:"MasterRecordId,omitempty"`
	AccountType       *string                `json:"Type,omitempty"`
	ParentId          *string                `json:"ParentId,omitempty"`
	BillingStreet     *string                `json:"BillingStreet,omitempty"`
	BillingCity       *string                `json:"BillingCity,omitempty"`
	BillingState      *string                `json:"BillingState,omitempty"`
	BillingPostalCode *string                `json:"BillingPostalCode,omitempty"`
	BillingCountry    *string                `json:"BillingCountry,omitempty"`
	AnnualRevenue     *int64                 `json:"AnnualRevenue,omitempty"`
	NumberOfEmployees *int16                 `json:"NumberOfEmployees,omitempty"`
	OwnerId           *string                `json:"OwnerId,omitempty"`
	AccountSource     *string                `json:"AccountSource,omitempty"`
}

type AccountAttributes struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type ContactQueryResponse struct {
	TotalSize      int       `json:"totalSize"`
	Done           bool      `json:"done"`
	Records        []Contact `json:"records"`
	NextRecordsUrl string    `json:"nextRecordsUrl"`
}

// Contact represents a Salesforce Contact record
type Contact struct {
	Id                 string                 `json:"id,omitempty"`
	AccountId          string                 `json:"AccountId,omitempty"`
	FirstName          string                 `json:"FirstName,omitempty"`
	LastName           string                 `json:"LastName"`       // Required
	Name               string                 `json:"Name,omitempty"` // Read-only, auto-generated
	Email              string                 `json:"Email,omitempty"`
	Phone              string                 `json:"Phone,omitempty"`
	MobilePhone        string                 `json:"MobilePhone,omitempty"`
	Title              string                 `json:"Title,omitempty"`
	Department         string                 `json:"Department,omitempty"`
	Description        string                 `json:"Description,omitempty"`
	MailingStreet      string                 `json:"MailingStreet,omitempty"`
	MailingCity        string                 `json:"MailingCity,omitempty"`
	MailingState       string                 `json:"MailingState,omitempty"`
	MailingPostalCode  string                 `json:"MailingPostalCode,omitempty"`
	MailingCountry     string                 `json:"MailingCountry,omitempty"`
	OtherStreet        string                 `json:"OtherStreet,omitempty"`
	OtherCity          string                 `json:"OtherCity,omitempty"`
	OtherState         string                 `json:"OtherState,omitempty"`
	OtherPostalCode    string                 `json:"OtherPostalCode,omitempty"`
	OtherCountry       string                 `json:"OtherCountry,omitempty"`
	Fax                string                 `json:"Fax,omitempty"`
	AssistantName      string                 `json:"AssistantName,omitempty"`
	AssistantPhone     string                 `json:"AssistantPhone,omitempty"`
	LeadSource         string                 `json:"LeadSource,omitempty"`
	Birthdate          *sftime.SalesforceTime `json:"Birthdate,omitempty"`
	CreatedDate        *sftime.SalesforceTime `json:"CreatedDate,omitempty"`        // Read-only
	LastModifiedDate   *sftime.SalesforceTime `json:"LastModifiedDate,omitempty"`   // Read-only
	SystemModstamp     *sftime.SalesforceTime `json:"SystemModstamp,omitempty"`     // Read-only
	LastActivityDate   *sftime.SalesforceTime `json:"LastActivityDate,omitempty"`   // Read-only
	IsDeleted          bool                   `json:"IsDeleted,omitempty"`          // Read-only
	DoNotCall          bool                   `json:"DoNotCall,omitempty"`          // custom needs __c
	HasOptedOutOfEmail bool                   `json:"HasOptedOutOfEmail,omitempty"` // custom needs __c
	HasOptedOutOfFax   bool                   `json:"HasOptedOutOfFax,omitempty"`   // custom needs __c
	OwnerId            string                 `json:"OwnerId,omitempty"`
	CreatedById        string                 `json:"CreatedById,omitempty"`      // Read-only
	LastModifiedById   string                 `json:"LastModifiedById,omitempty"` // Read-only

	LinkedIn_Profile__c     string `json:"LinkedIn_Profile__c,omitempty"`
	Facebook_Friend__c      bool   `json:"Facebook_Friend__c,omitempty"`
	Type__c                 string `json:"Type__c,omitempty"`
	Scheduling_Site__c      string `json:"Scheduling_Site__c,omitempty"`
	Name__c                 string `json:"Name__c,omitempty"`
	Suffix__c               string `json:"Suffix__c,omitempty"`
	Birth_Year__c           string `json:"Birth_Year__c,omitempty"`
	Middle_Name__c          string `json:"Middle_Name__c string,omitempty"`
	Non_Standard_Address__c string `json:"Non_Standard_Address__c,omitempty"`
	HasEmailPermission__c   bool   `json:"HasEmailPermission__c,omitempty"`
	HasPhonePermission__c   bool   `json:"HasPhonePermission__c,omitempty"`
	HasSMSPermission__c     bool   `json:"HasSMSPermission__c,omitempty"`
	PreferredName__c        string `json:"PreferredName__c,omitempty"`

	HomePhone              string `json:"HomePhone,omitempty"`
	MasterRecordId         string `json:"MasterRecordId,omitempty"`
	Salutation             string `json:"Salutation,omitempty"`
	OtherLatitude          string `json:"OtherLatitude,omitempty"`
	OtherLongitude         string `json:"OtherLongitude,omitempty"`
	OtherGeocodeAccuracy   string `json:"OtherGeocodeAccuracy,omitempty"`
	OtherAddress           string `json:"OtherAddress,omitempty"`
	OtherPhone             string `json:"OtherPhone,omitempty"`
	MailingLatitude        string `json:"MailingLatitude,omitempty"`
	MailingLongitude       string `json:"MailingLongitude,omitempty"`
	MailingGeocodeAccuracy string `json:"MailingGeocodeAccuracy,omitempty"`
	MailingAddress         string `json:"MailingAddress,omitempty"`
	ReportsToId            string `json:"ReportsToId,omitempty"`
	LastCURequestDate      string `json:"LastCURequestDate,omitempty"`
	LastCUUpdateDate       string `json:"LastCUUpdateDate,omitempty"`
	LastReferencedDate     string `json:"LastReferencedDate,omitempty"`
	EmailBouncedReason     string `json:"EmailBouncedReason,omitempty"`
	EmailBouncedDate       string `json:"EmailBouncedDate,omitempty"`
	IsEmailBounced         string `json:"IsEmailBounced ,omitempty"`
	PhotoURL               string `json:"PhotoURL,omitempty"`
	Jigsaw                 string `json:"Jigsaw,omitempty"`
	JigsawContactId        string `json:"JigsawContactId,omitempty"`
	IndividualId           string `json:"IndividualId,omitempty"`
	IsPriorityRecord       string `json:"IsPriorityRecord,omitempty"`
}
