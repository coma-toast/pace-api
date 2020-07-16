package project

import (
	"time"

	"github.com/coma-toast/pace-api/pkg/contact"
)

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
	ProjectManager      contact.Contact
	ClientName          contact.Contact
	EORName             contact.Contact
	DetailerName        contact.Contact
	InspectionLab       contact.Contact
	SteelErectorName    contact.Contact
	SteelFabricatorName contact.Contact
	GeneralContractor   contact.Contact
	PrimaryContactName  contact.Contact
	PrimaryContactPhone string
	PrimaryContactEmail string
	SquareFootage       int32
	WeightInTons        int32
}
