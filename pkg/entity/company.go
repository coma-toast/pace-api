package entity

// Company is a contact company
type Company struct {
	ID             string
	Created        string
	Name           string
	PrimaryContact string
	Contacts       []string
	Phone          string
	Email          string
	Address        string
	City           string
	State          string
	Zip            int32
}
