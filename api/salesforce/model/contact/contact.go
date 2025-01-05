package salesforce

import (
	"fmt"
	"time"

	salesforcetime "github.com/htstinson/stinsondataapi/api/salesforce/auth"
)

// Contact represents a Salesforce Contact record
type Contact struct {
	Id                 string               `json:"id,omitempty"`
	AccountId          string               `json:"AccountId,omitempty"`
	FirstName          string               `json:"FirstName,omitempty"`
	LastName           string               `json:"LastName"`       // Required
	Name               string               `json:"Name,omitempty"` // Read-only, auto-generated
	Email              string               `json:"Email,omitempty"`
	Phone              string               `json:"Phone,omitempty"`
	MobilePhone        string               `json:"MobilePhone,omitempty"`
	Title              string               `json:"Title,omitempty"`
	Department         string               `json:"Department,omitempty"`
	Description        string               `json:"Description,omitempty"`
	MailingStreet      string               `json:"MailingStreet,omitempty"`
	MailingCity        string               `json:"MailingCity,omitempty"`
	MailingState       string               `json:"MailingState,omitempty"`
	MailingPostalCode  string               `json:"MailingPostalCode,omitempty"`
	MailingCountry     string               `json:"MailingCountry,omitempty"`
	OtherStreet        string               `json:"OtherStreet,omitempty"`
	OtherCity          string               `json:"OtherCity,omitempty"`
	OtherState         string               `json:"OtherState,omitempty"`
	OtherPostalCode    string               `json:"OtherPostalCode,omitempty"`
	OtherCountry       string               `json:"OtherCountry,omitempty"`
	Fax                string               `json:"Fax,omitempty"`
	AssistantName      string               `json:"AssistantName,omitempty"`
	AssistantPhone     string               `json:"AssistantPhone,omitempty"`
	LeadSource         string               `json:"LeadSource,omitempty"`
	Birthdate          *salesforcetime.Salesforce.Time `json:"Birthdate,omitempty"`
	CreatedDate        *time.Time           `json:"CreatedDate,omitempty"`      // Read-only
	LastModifiedDate   *time.Time           `json:"LastModifiedDate,omitempty"` // Read-only
	SystemModstamp     *time.Time           `json:"SystemModstamp,omitempty"`   // Read-only
	LastActivityDate   *time.Time           `json:"LastActivityDate,omitempty"` // Read-only
	IsDeleted          bool                 `json:"IsDeleted,omitempty"`        // Read-only
	DoNotCall          bool                 `json:"DoNotCall,omitempty"`
	HasOptedOutOfEmail bool                 `json:"HasOptedOutOfEmail,omitempty"`
	HasOptedOutOfFax   bool                 `json:"HasOptedOutOfFax,omitempty"`
	OwnerId            string               `json:"OwnerId,omitempty"`
	CreatedById        string               `json:"CreatedById,omitempty"`      // Read-only
	LastModifiedById   string               `json:"LastModifiedById,omitempty"` // Read-only
	RecordTypeId       string               `json:"RecordTypeId,omitempty"`
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
