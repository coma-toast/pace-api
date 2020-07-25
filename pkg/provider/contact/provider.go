package contact

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with contact data
type Provider interface {
	// GetBy(contactname string) (entity.Contact, error)
	GetAll() ([]entity.Contact, error)
	Addcontact(entity.Contact) (entity.Contact, error)
	Updatecontact(entity.Contact) (entity.Contact, error)
	Deletecontact(contactname entity.Contact) error
}
