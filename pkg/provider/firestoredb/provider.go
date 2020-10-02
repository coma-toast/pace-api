package firestoredb

// Provider is for working with contact data
type Provider interface {
	GetAll(target interface{}) error
	GetByID(ID string, target interface{}) error
	GetFirstBy(path string, op string, value string, target interface{}) error
	Set(ID string, data interface{}) error
	Delete(ID string) error
}
