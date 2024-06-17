package types

import (
	"encoding/json"
	"time"
)

const (
	millisInSecond = 1000
	secondsInDay   = 86400
)

// CustpTime wraps time.Time to handle different JSON time formats.
type CustpTime struct {
	time.Time
}

// MarshalJSON customizes the JSON representation of CustpTime to be an integer.
// Converts time to milliseconds since epoch.
func (ct CustpTime) MarshalJSON() ([]byte, error) {
	// Convert the time to milliseconds since epoch
	millis := ct.UnixNano() / int64(time.Millisecond)
	return json.Marshal(millis) // Return the JSON representation
}

func (ct *CustpTime) UnmarshalJSON(b []byte) error {
	s := string(b)

	// Check if the value is numeric (either days or milliseconds)
	var numericValue int64
	if err := json.Unmarshal(b, &numericValue); err == nil {
		// Determine if we're dealing with days or milliseconds by magnitude
		if numericValue > 1e10 { // Assuming any large number is a timestamp in milliseconds
			t := time.Unix(0, numericValue*int64(time.Millisecond)).UTC()
			ct.Time = t
			return nil
		} else { // Days since the Unix epoch
			t := time.Unix(numericValue*int64(secondsInDay), 0).UTC()
			ct.Time = t
			return nil
		}
	}

	// Try parsing as a string date if not numeric
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// ToTimePointer returns the embedded time.Time as a *time.Time
func (ct *CustpTime) ToTimePointer() *time.Time {
	if ct == nil {
		return nil
	}
	return &ct.Time
}

func serializeDate(record *int) interface{} {
	if record == nil {
		return nil // Return nil if the pointer is nil
	}
	daysSinceEpoch := *record
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC) // Epoch reference
	date := epoch.AddDate(0, 0, daysSinceEpoch)          // Add the number of days to the epoch
	return map[string]interface{}{"int.date": date}
}
