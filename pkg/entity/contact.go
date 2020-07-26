package entity

// Contact is a non-user contact
type Contact struct {
	ID        string `json:"id"`
	Created   string `json:"created"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Company   string `json:"company"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Timezone  string `json:"timezone"`
}

// Company is a contact company
type Company struct {
	ID             string
	PrimaryContact Contact
	Contacts       []Contact
	Phone          string
	Email          string
	Created        string
	Address        string
	City           string
	State          string
	Zip            int32
}
