package models

import "time"

type ServerStatus struct {
	Status      string
	StatusColor int
	TestType    string
	Access      string
	StartDate   time.Time
	EndDate     time.Time
	Duration    time.Duration
	Error bool
	ErrorStatus string
}