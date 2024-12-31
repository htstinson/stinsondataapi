package model

import sftime "github.com/htstinson/stinsondataapi/api/salesforce/time"

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
