package entity

// Company is a contact company
type Company struct {
	ID             string `json:"id"`
	Created        string `json:"created"`
	Name           string `json:"name"`
	PrimaryContact string `json:"primaryContact"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	City           string `json:"city"`
	State          string `json:"state"`
	Zip            string `json:"zip"`
	Favorite       bool   `json:"favorite"`
	Deleted        bool   `json:"deleted"`
	Instance       string `json:"instance"`
}
