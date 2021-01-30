package entity

// If you want historical data, leave it.
// If you want phone number update, for example, use a contact id
// Project is a construction project
type Project struct {
	ID                    string `json:"id"`
	Created               string `json:"created"`
	Deleted               bool   `json:"deleted"`
	Name                  string `json:"name"`
	StartDate             string `json:"startDate"`
	DueDate               string `json:"dueDate"`
	Address               string `json:"address"`
	City                  string `json:"city"`
	State                 string `json:"state"`
	Zip                   int32  `json:"zip"`
	ProjectManager        string `json:"projectManager"`
	ClientID              string `json:"clientID"`
	EORNameID             string `json:"eORNameID"`
	DetailerNameID        string `json:"detailerNameID"`
	InspectionLabID       string `json:"inspectionLabID"`
	SteelErectorNameID    string `json:"steelErectorNameID"`
	SteelFabricatorNameID string `json:"steelFabricatorNameID"`
	GeneralContractorID   string `json:"generalContractorID"`
	PrimaryContactNameID  string `json:"primaryContactNameID"`
	PrimaryContactPhone   string `json:"primaryContactPhone"`
	PrimaryContactEmail   string `json:"primaryContactEmail"`
	SquareFootage         int32  `json:"squareFootage"`
	WeightInTons          int32  `json:"weightInTons"`
}

// UpdateProjectRequest is a construction project
type UpdateProjectRequest struct {
	Name                  string `json:"name"`
	Deleted               bool   `json:"deleted"`
	StartDate             string `json:"startDate"`
	DueDate               string `json:"dueDate"`
	Address               string `json:"address"`
	City                  string `json:"city"`
	State                 string `json:"state"`
	Zip                   int32  `json:"zip"`
	ProjectManager        string `json:"projectManager"`
	ClientID              string `json:"clientID"`
	EORNameID             string `json:"eORNameID"`
	DetailerNameID        string `json:"detailerNameID"`
	InspectionLabID       string `json:"inspectionLabID"`
	SteelErectorNameID    string `json:"steelErectorNameID"`
	SteelFabricatorNameID string `json:"steelFabricatorNameID"`
	GeneralContractorID   string `json:"generalContractorID"`
	PrimaryContactNameID  string `json:"primaryContactNameID"`
	PrimaryContactPhone   string `json:"primaryContactPhone"`
	PrimaryContactEmail   string `json:"primaryContactEmail"`
	SquareFootage         int32  `json:"squareFootage"`
	WeightInTons          int32  `json:"weightInTons"`
}
