package contact

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with contact data
type Provider interface {
	// GetBy(contactname string) (entity.Contact, error)
	GetAll() ([]entity.Contact, error)
	AddContact(entity.Contact) (entity.Contact, error)
	UpdateContact(entity.Contact) (entity.Contact, error)
	DeleteContact(contactname entity.Contact) error
}
