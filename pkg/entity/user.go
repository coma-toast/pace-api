package entity

// User is the user data
type User struct {
	ID        int
	Created   string
	FirstName string
	LastName  string
	Role      string
	Username  string
	Password  string
	Email     string
	Phone     string
	TimeZone  string
	DarkMode  bool
}
