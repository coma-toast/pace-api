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
