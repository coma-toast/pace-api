package entity

// If you want historical data, leave it.
// If you want phone number update, for eample, use a contact id
// Project is a construction project
type Project struct {
	ID                  string
	Created             string
	Name                string
	StartDate           string
	DueDate             string
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

// UpdateProjectRequest is a construction project
type UpdateProjectRequest struct {
	Name                string
	StartDate           string
	DueDate             string
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
