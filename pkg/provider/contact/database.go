package contact

import (
	"errors"
	"fmt"
	"time"

	"github.com/coma-toast/pace-api/pkg/entity"
	"github.com/coma-toast/pace-api/pkg/provider/firestoredb"
	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

// DatabaseProvider is a contact.Provider the uses a database
type DatabaseProvider struct {
	SharedProvider *firestoredb.DatabaseProvider
}

// ErrContactNotFound if no Contacts are found
var ErrContactNotFound = errors.New("Contact not found")

// GetAll gets a Contact by ID
func (d *DatabaseProvider) GetAll() ([]entity.Contact, error) {
	var allContactData []entity.Contact

	err := d.SharedProvider.GetAll(&allContactData)
	if err != nil {
		return []entity.Contact{}, err
	}

	return allContactData, nil
}

// // GetByName gets a Contact by name
// func (d *DatabaseProvider) GetByName(contactName string) (entity.Contact, error) {
// 	var contact entity.Contact
// 	err := d.SharedProvider.GetFirstBy("Name", "==", contactName, &contact)
// 	if err != nil {
// 		return entity.Contact{}, err
// 	}

// 	return contact, nil
// }

// Add is to update a Contact record
func (d *DatabaseProvider) Add(newContactData entity.Contact) (entity.Contact, error) {
	rollbar.Info(fmt.Sprintf("Adding new Contact to DB %s %s", newContactData.FirstName, newContactData.LastName))

	// There can be multiple contacts with the same name.
	// existingContact, err := d.SharedProvider.Get(newContactData.Name)
	// if (entity.Contact{}) != existingContact {
	// 	return entity.Contact{}, fmt.Errorf("Error adding Contact %s: Contactname already exists. ID: %s", newContactData.Contactname, existingContact.ID)
	// }

	newUUID := uuid.New().String()
	newContactData = entity.Contact{
		ID:        newUUID,
		Created:   time.Now().String(),
		FirstName: newContactData.FirstName,
		LastName:  newContactData.LastName,
		Company:   newContactData.Company,
		Email:     newContactData.Email,
		Phone:     newContactData.Phone,
		Timezone:  newContactData.Timezone,
		Instance:  newContactData.Instance,
	}
	err := d.SharedProvider.Set(newContactData.ID, newContactData)
	if err != nil {
		return entity.Contact{}, fmt.Errorf("Error setting Contact %s %s by ID: %w", newContactData.FirstName, newContactData.LastName, err)
	}

	var newContact = entity.Contact{}
	err = d.SharedProvider.GetByID(newContactData.ID, &newContact)
	if err != nil {
		return entity.Contact{}, fmt.Errorf("Error getting newly created Contact %s %s by ID: %s", newContactData.FirstName, newContactData.LastName, err)
	}

	rollbar.Info(fmt.Sprintf("Contact %s %s added.", newContactData.FirstName, newContactData.LastName))

	return newContact, nil
}

// Update is to update a Contact record
func (d *DatabaseProvider) Update(newContactData entity.Contact) (entity.Contact, error) {
	var currentContactData = entity.Contact{}
	err := d.SharedProvider.GetByID(newContactData.ID, &currentContactData)
	if err != nil {
		return entity.Contact{}, err
	}
	rollbar.Info(fmt.Sprintf("Updating ContactID %s. \nOld Data: %v \nNew Data: %v", currentContactData.ID, currentContactData, newContactData))
	updatedContact := entity.Contact{
		ID:        currentContactData.ID,
		Created:   currentContactData.Created,
		FirstName: newContactData.FirstName,
		LastName:  newContactData.LastName,
		Company:   newContactData.Company,
		Email:     newContactData.Email,
		Phone:     newContactData.Phone,
		Timezone:  newContactData.Timezone,
		Favorite:  newContactData.Favorite,
		Deleted:   newContactData.Deleted,
		Instance:  newContactData.Instance,
	}

	err = d.SharedProvider.Set(currentContactData.ID, updatedContact)
	if err != nil {
		return entity.Contact{}, err
	}

	var updatedContactData = entity.Contact{}
	err = d.SharedProvider.GetByID(updatedContact.ID, &updatedContactData)
	if err != nil {
		return entity.Contact{}, err
	}

	return updatedContactData, nil
}

// Delete deletes a contact
func (d *DatabaseProvider) Delete(contact entity.Contact) error {
	rollbar.Info(fmt.Sprintf("Deleting Contact from DB: %s", contact.ID))
	var currentContact entity.Contact

	err := d.SharedProvider.GetByID(contact.ID, &currentContact)
	if (entity.Contact{}) == currentContact {
		return fmt.Errorf("Contact not found")
	}

	err = d.SharedProvider.Delete(contact.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted Contact %s %s: %s", contact.FirstName, contact.LastName, contact.ID))

	return nil
}
