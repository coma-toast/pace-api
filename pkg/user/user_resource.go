package user

import "time"

// User is the user data
type User struct {
	FirstName string
	LastName  string
	Username  string
	Email     string
	DarkMode  bool
	Created   time.Time
}
