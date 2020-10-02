package entity

// Inspection is an inspection report
type Inspection struct {
	ID             string `json:"id"`
	Created        string `json:"created"`
	ProjectID      string `json:"projectID"`
	Username       string `json:"username"`
	StartTime      string `json:"startTime"`
	EndTime        string `json:"endTime"`
	InspectedParts string `json:"inspectedParts"`
}

// UpdateInspectionRequest is an inspection report
type UpdateInspectionRequest struct {
	ID             string `json:"id"`
	ProjectID      string `json:"projectID"`
	Username       string `json:"username"`
	StartTime      string `json:"startTime"`
	EndTime        string `json:"endTime"`
	InspectedParts string `json:"inspectedParts"`
}
