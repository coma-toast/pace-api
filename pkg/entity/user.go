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

// FirestoreUser is the user data
type FirestoreUser struct {
	ID        string `firestore:"id"`
	Created   string `firestore:"created"`
	FirstName string `firestore:"firstname"`
	LastName  string `firestore:"lastname"`
	Role      string `firestore:"role"`
	Username  string `firestore:"username"`
	Password  string `firestore:"-"`
	Email     string `firestore:"email"`
	Phone     string `firestore:"phone"`
	TimeZone  string `firestore:"timezone"`
	DarkMode  bool   `firestore:"darkmode"`
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
