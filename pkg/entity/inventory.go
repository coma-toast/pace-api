package entity

// TODO: the "Shape" can be one of 18000+ shapes. See the iOS app "Steel".
// We are also allowed to import the "Steel" DB to this app

// Inventory is an inventory item
type Inventory struct {
	ID        string `json:"iD"`
	Created   string `json:"created"`
	ProjectID string `json:"projectID"`
	Stage     Stage  `json:"stage"`
	Size      int32  `json:"size"`
	Length    int32  `json:"length"`
	Grade     int32  `json:"grade"`
	Shape     string `json:"shape"`
	Passed    bool   `json:"passed"`
	Sequence  int32  `json:"sequence"`
	Priority  int32  `json:"priority"`
}

// Stage is what stage the inventory item is in
type Stage struct {
	Raw       bool `json:"raw"`
	InProcess bool `json:"inProcess"`
	OnHold    bool `json:"onHold"`
	Finished  bool `json:"finished"`
}

// UpdateInventoryRequest is an inventory item
type UpdateInventoryRequest struct {
	ID        string `json:"iD"`
	ProjectID string `json:"projectID"`
	Stage     Stage  `json:"stage"`
	Size      int32  `json:"size"`
	Length    int32  `json:"length"`
	Grade     int32  `json:"grade"`
	Shape     string `json:"shape"`
	Passed    bool   `json:"passed"`
	Sequence  int32  `json:"sequence"`
	Priority  int32  `json:"priority"`
}
