package entity

// Company is a contact company
type Company struct {
	ID             string
	PrimaryContact string
	Contacts       []string
	Phone          string
	Email          string
	Created        string
	Address        string
	City           string
	State          string
	Zip            int32
}
