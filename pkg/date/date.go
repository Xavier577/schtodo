package date

import (
	"encoding/json"
	"fmt"
	"time"
)

// DateTime wraps around time.DateTime to provide custom JSON unmarshaling
type DateTime struct {
	time.Time
}

// Ensure Time implements the json.Unmarshaler interface
var _ json.Unmarshaler = &DateTime{}

// UnmarshalJSON parses the JSON string in the custom format YYYY-MM-DDTHH:MM
func (mt *DateTime) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}

	// Parse the string into time.Time using the correct format
	t, err := time.Parse("2006-01-02T15:04", s)
	if err != nil {
		return fmt.Errorf("invalid time format: %v", err)
	}

	mt.Time = t
	return nil
}

// MarshalJSON converts the Time back to JSON
func (mt DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.Format("2006-01-02T15:04"))
}

// String returns the string representation of the time
func (mt DateTime) String() string {
	return mt.Format("2006-01-02T15:04")
}
