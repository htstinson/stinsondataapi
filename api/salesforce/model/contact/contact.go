package salesforce

import (
	"fmt"

	salesforcetime "github.com/htstinson/stinsondataapi/api/salesforce/time"
)

// Contact represents a Salesforce Contact record
type Contact struct {
	Id                 string                         `json:"id,omitempty"`
	AccountId          string                         `json:"AccountId,omitempty"`
	FirstName          string                         `json:"FirstName,omitempty"`
	LastName           string                         `json:"LastName"`       // Required
	Name               string                         `json:"Name,omitempty"` // Read-only, auto-generated
	Email              string                         `json:"Email,omitempty"`
	Phone              string                         `json:"Phone,omitempty"`
	MobilePhone        string                         `json:"MobilePhone,omitempty"`
	Title              string                         `json:"Title,omitempty"`
	Department         string                         `json:"Department,omitempty"`
	Description        string                         `json:"Description,omitempty"`
	MailingStreet      string                         `json:"MailingStreet,omitempty"`
	MailingCity        string                         `json:"MailingCity,omitempty"`
	MailingState       string                         `json:"MailingState,omitempty"`
	MailingPostalCode  string                         `json:"MailingPostalCode,omitempty"`
	MailingCountry     string                         `json:"MailingCountry,omitempty"`
	OtherStreet        string                         `json:"OtherStreet,omitempty"`
	OtherCity          string                         `json:"OtherCity,omitempty"`
	OtherState         string                         `json:"OtherState,omitempty"`
	OtherPostalCode    string                         `json:"OtherPostalCode,omitempty"`
	OtherCountry       string                         `json:"OtherCountry,omitempty"`
	Fax                string                         `json:"Fax,omitempty"`
	AssistantName      string                         `json:"AssistantName,omitempty"`
	AssistantPhone     string                         `json:"AssistantPhone,omitempty"`
	LeadSource         string                         `json:"LeadSource,omitempty"`
	Birthdate          *salesforcetime.SalesforceTime `json:"Birthdate,omitempty"`
	CreatedDate        *salesforcetime.SalesforceTime `json:"CreatedDate,omitempty"`        // Read-only
	LastModifiedDate   *salesforcetime.SalesforceTime `json:"LastModifiedDate,omitempty"`   // Read-only
	SystemModstamp     *salesforcetime.SalesforceTime `json:"SystemModstamp,omitempty"`     // Read-only
	LastActivityDate   *salesforcetime.SalesforceTime `json:"LastActivityDate,omitempty"`   // Read-only
	IsDeleted          bool                           `json:"IsDeleted,omitempty"`          // Read-only
	DoNotCall          bool                           `json:"DoNotCall,omitempty"`          // custom needs __c
	HasOptedOutOfEmail bool                           `json:"HasOptedOutOfEmail,omitempty"` // custom needs __c
	HasOptedOutOfFax   bool                           `json:"HasOptedOutOfFax,omitempty"`   // custom needs __c
	OwnerId            string                         `json:"OwnerId,omitempty"`
	CreatedById        string                         `json:"CreatedById,omitempty"`      // Read-only
	LastModifiedById   string                         `json:"LastModifiedById,omitempty"` // Read-only

	LinkedIn_Profile__c     string `json:"LinkedIn_Profile__c,omitempty"`
	Facebook_Friend__c      string `json:"Facebook_Friend__c,omitempty"`
	Type__c                 string `json:"Type__c,omitempty"`
	Scheduling_Site__c      string `json:"Scheduling_Site__c,omitempty"`
	Name__c                 string `json:"Name__c,omitempty"`
	Suffix_c                string `json:"Suffix__c,omitempty"`
	Birth_Year__c           string `json:"Birth_Year__c,omitempty"`
	Middle_Name__c          string `json:"Middle_Name__c string,omitempty"`
	Non_Standard_Address__c string `json:"Non_Standard_Address__c,omitempty"`
	HasEmailPermission__c   string `json:"HasEmailPermission__c,omitempty"`
	HasPhonePermission__c   string `json:"HasPhonePermission__c,omitempty"`
	HasSMSPermission__c     string `json:"HasSMSPermission__c,omitempty"`
	PreferredName__c        string `json:"PreferredName__c,omitempty"`

	Smagicinteract__SMSOptOut__c string `json:"smagicinteract__SMSOptOut__c,omitempty"`
	Smagicinteract__Contact__c   string `json:"smagicinteract__Contact__c,omitempty"`

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
	LastViewDate           string `json:"LastViewDate,omitempty"`
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

// NewContact creates a new Contact with required fields
func NewContact(lastName string) *Contact {
	return &Contact{
		LastName: lastName,
	}
}

// Validate checks if required fields are present
func (c *Contact) Validate() error {
	if c.LastName == "" {
		return fmt.Errorf("LastName is required")
	}
	return nil
}

// Example usage:
/*
   contact := NewContact("Smith")
   contact.FirstName = "John"
   contact.Email = "john.smith@example.com"

   // Validate before sending to API
   if err := contact.Validate(); err != nil {
       // Handle error
   }
*/
