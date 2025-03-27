package salesforcetime

import (
	"fmt"

	"strings"
	"time"
)

type SalesforceTime time.Time

func (st *SalesforceTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from string
	str := string(data)
	str = strings.Trim(str, `"`)

	fmt.Printf("test %s\n", str)

	// Use layout matching exact Salesforce format
	t, err := time.Parse("2006-01-02T15:04:05.000-0700", str)
	if err != nil {
		// If that fails, try alternate format with +0000
		t, err = time.Parse("2006-01-02T15:04:05.000+0000", str)
		if err != nil {
			fmt.Printf("[%v] ********* error parsing time ************* %s\n", time.Now().Format(time.RFC3339), err.Error())
			*st = SalesforceTime{}
			//return err
			return nil
		}
	}

	*st = SalesforceTime(t)
	return nil
}

// Add Format method
func (st SalesforceTime) Format(layout string) string {
	return time.Time(st).Format(layout)
}

// Add Time method to convert back to time.Time
func (st SalesforceTime) Time() time.Time {
	return time.Time(st)
}

// IsZero reports whether t represents the zero time instant
func (st SalesforceTime) IsZero() bool {
	return time.Time(st).IsZero()
}

// Before reports whether the time instant t is before u
func (st SalesforceTime) Before(u SalesforceTime) bool {
	return time.Time(st).Before(time.Time(u))
}

// After reports whether the time instant t is after u
func (st SalesforceTime) After(u SalesforceTime) bool {
	return time.Time(st).After(time.Time(u))
}

// Equal reports whether t and u represent the same time instant
func (st SalesforceTime) Equal(u SalesforceTime) bool {
	return time.Time(st).Equal(time.Time(u))
}

// Sub returns the duration t-u
func (st SalesforceTime) Sub(u SalesforceTime) time.Duration {
	return time.Time(st).Sub(time.Time(u))
}

// MarshalJSON implements json.Marshaler
func (st SalesforceTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(st).Format("2006-01-02T15:04:05.000-0700") + `"`), nil
}

// String implements fmt.Stringer
func (st SalesforceTime) String() string {
	return time.Time(st).String()
}
