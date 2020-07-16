package user

import "time"

// User is the user data
type User struct {
	ID        string
	Created   time.Time
	FirstName string
	LastName  string
	Role      []string
	Username  string
	Password  string
	Email     string
	Phone     string
	TimeZone  time.Location
	DarkMode  bool
}
