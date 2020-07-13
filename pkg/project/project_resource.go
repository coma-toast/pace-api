package project

import "time"

// Project is a construction project
type Project struct {
	Created             time.Time
	StartDate           time.Time
	DueDate             time.Time
	Address             string
	City                string
	State               string
	Zip                 int32
	PrimaryContactName  string
	PrimaryContactPhone string
	PrimaryContactEmail string
}
