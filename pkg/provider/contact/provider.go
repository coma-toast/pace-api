package contact

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with contact data
type Provider interface {
	// GetBy(contactname string) (entity.Contact, error)
	GetAll() ([]entity.Contact, error)
	Add(entity.Contact) (entity.Contact, error)
	Update(entity.Contact) (entity.Contact, error)
	Delete(contactname entity.Contact) error
}
