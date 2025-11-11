package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		return nil
	}

	// Try parsing with milliseconds
	t, err := time.Parse("2006-01-02T15:04:05.999999", s)
	if err != nil {
		// Try without milliseconds
		t, err = time.Parse("2006-01-02T15:04:05", s)
		if err != nil {
			return fmt.Errorf("invalid time format: %s", s)
		}
	}
	ct.Time = t
	return nil
}

// BSON format: { "time": ISODate(...) }
func (ct *CustomTime) UnmarshalBSON(data []byte) error {

	var timeStr string
	if err := bson.Unmarshal(data, &timeStr); err == nil {
		// Parse the ISO 8601 string format without milliseconds: 2025-09-09T07:05:09
		t, err := time.Parse("2006-01-02T15:04:05", timeStr)
		if err != nil {
			return fmt.Errorf("failed to parse time string '%s': %w", timeStr, err)
		}
		ct.Time = t
		return nil
	}

	var wrapper struct {
		Time time.Time `bson:"time"`
	}
	if err := bson.Unmarshal(data, &wrapper); err != nil {
		return err
	}
	ct.Time = wrapper.Time
	return nil
}

// Optional: for debugging
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Time.Format(time.RFC3339Nano))
}
