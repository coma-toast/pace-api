package entity

// Company is a contact company
type Company struct {
	ID             string `json:"id"`
	Created        string `json:"created"`
	Name           string `json:"name"`
	PrimaryContact string `json:"primaryContact"`
	Contacts       string `json:"contacts"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	City           string `json:"city"`
	State          string `json:"state"`
	Zip            int32  `json:"zip"`
	Favorite       bool   `json:"favorite"`
	Deleted        bool   `json:"deleted"`
	Instance       string `json:"instance"`
}
