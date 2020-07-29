package firestoredb

import (
	firestore "cloud.google.com/go/firestore/apiv1"
)

// Provider is for working with contact data
type Provider interface {
	// GetBy(contactname string) (entity.Contact, error)
	GetAll(collection string) ([]*firestore.DocumentSnapshot, error)
	// AddItem(??) (*[]firestore.DocumentSnapshot, error)
	// UpdateItem(entity.Contact) (entity.Contact, error)
	// DeleteItem(contactname entity.Contact) error
}
