package entity

// Contact is a non-user contact
type Contact struct {
	ID        string
	Created   string
	FirstName string
	LastName  string
	Username  string
	Password  string
	Email     string
	Phone     string
	Timezone  string
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
