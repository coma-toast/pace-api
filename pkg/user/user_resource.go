package user

import "time"

// User is the user data
type User struct {
	id        string
	Created   time.Time
	FirstName string
	LastName  string
	Username  string
	Email     string
	Phone     string
	TimeZone  time.TimeZone
	DarkMode  bool
}
