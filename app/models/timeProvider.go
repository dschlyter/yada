package models

import (
	"time"
)

var mockTime *time.Time = nil

func SetMockTime(newTime string) (err error) {
	parsed, err := time.Parse("2006-01-02T15:04:05", newTime)
	if err == nil {
		mockTime = &parsed
	}
	return
}

func ClearMockTime() {
	mockTime = nil
}

func GetTime() time.Time {
	if mockTime == nil {
		return time.Now()
	} else {
		return *mockTime
	}
}
