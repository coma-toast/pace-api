package contact

import "time"

// Contact is a non-user contact
type Contact struct {
	ID        string
	Created   time.Time
	FirstName string
	LastName  string
	Username  string
	Password  string
	Email     string
	Phone     string
	TimeZone  time.Location
}

// Company is a contact company
type Company struct {
	ID             string
	PrimaryContact Contact
	Contacts       []Contact
	Phone          string
	Email          string
	Created        time.Time
	Address        string
	City           string
	State          string
	Zip            int32
}
