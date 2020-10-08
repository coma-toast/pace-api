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
	Favorite  bool   `json:"favorite"`
	Deleted   bool   `json:"deleted"`
	Instance  string `json:"instance"`
}
