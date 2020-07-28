package entity

// User is the user data
type User struct {
	ID        string `json:"id"`
	Created   string `json:"created"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	TimeZone  string `json:"timezone"`
	DarkMode  bool   `json:"darkmode"`
}

// UpdateUserRequest is a passwordless user entity
type UpdateUserRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	TimeZone  string `json:"timezone"`
	DarkMode  bool   `json:"darkmode"`
}
