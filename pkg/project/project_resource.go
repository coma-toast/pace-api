package project

import "time"

// TODO: contact structs
// If you want historical data, leave it.
// If you want phone number update, for eample, use a contact id
// Project is a construction project
type Project struct {
	ProjectID           string
	Created             time.Time
	StartDate           time.Time
	DueDate             time.Time
	Address             string
	City                string
	State               string
	Zip                 int32
	ProjectManager      string
	ClientName          string
	EORName             string
	DetailerName        string
	InspectionLab       string
	SteelErectorName    string
	SteelFabricatorName string
	GeneralContractor   string
	PrimaryContactName  string
	PrimaryContactPhone string
	PrimaryContactEmail string
	SquareFootage       int32
	WeightInTons        int32
}
