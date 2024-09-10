package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type CustomTime time.Time

// UnmarshalJSON parses time in HH:mm format
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b[1 : len(b)-1]) // remove quotes
	parsedTime, err := time.Parse("15:04", s)
	if err != nil {
		return errors.New("invalid time format, expected HH:mm")
	}
	*ct = CustomTime(parsedTime)
	return nil
}

// MarshalJSON serializes time to HH:mm format
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(ct).Format("15:04") + `"`), nil
}

// ToTime returns the time.Time value of CustomTime
func (ct CustomTime) ToTime() time.Time {
	return time.Time(ct)
}

// IsZero checks if CustomTime is zero value
func (ct CustomTime) IsZero() bool {
	return time.Time(ct).IsZero()
}

func Debug(obj any) {
	raw, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(raw))
}

func LocalTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

func ConvertStringTimeToTime(t string) time.Time {
	layout := "2006-01-02 15:04:05.999 -0700 MST"
	result, err := time.Parse(layout, t)
	if err != nil {
	    log.Printf("Error: Parse time failed: %s", err.Error())
	}
	return result
}

func ParseTimeOnly(timeStr string) CustomTime {
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return CustomTime{}
	}
	return CustomTime(parsedTime)
}
