package contact

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/coma-toast/pace-api/pkg/entity"
	helper "github.com/coma-toast/pace-api/pkg/utils"
	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

// DatabaseProvider is a contact.Provider the uses a database
type DatabaseProvider struct {
	Database *firestore.Client
}

// ErrContactNotFound if no Contacts are found
var ErrContactNotFound = errors.New("Contact not found")

// // GetByContactname gets a Contact by Contactname
// func (d *DatabaseProvider) GetByContactname(Contactname string) (entity.Contact, error) {
// 	return d.getByContactname(Contactname)
// }

// GetAll gets a Contact by Contactname
func (d *DatabaseProvider) GetAll() ([]entity.Contact, error) {
	return d.getAll()
}

func (d *DatabaseProvider) getAll() ([]entity.Contact, error) {
	var contacts []entity.Contact

	allContactData, err := d.Database.Collection("contacts").Documents(context.TODO()).GetAll()
	if err != nil {
		return []entity.Contact{}, err
	}

	for _, ContactData := range allContactData {
		var Contact entity.Contact
		err := ContactData.DataTo(&Contact)
		if err != nil {
			return []entity.Contact{}, fmt.Errorf("ERROR: GetAll(): Firestore.DataTo() error %w", err)
		}
		contacts = append(contacts, Contact)
	}

	return contacts, nil
}

// AddContact is to update a Contact record
func (d *DatabaseProvider) AddContact(newContactData entity.Contact) (entity.Contact, error) {
	contactRef, err := d.addContact(newContactData)
	if err != nil {
		return entity.Contact{}, err
	}
	rollbar.Info(fmt.Sprintf("Adding new Contact %s %s", newContactData.FirstName, newContactData.LastName))
	updatedContactData, err := d.getByContactID(contactRef.ID)
	if err != nil {
		return entity.Contact{}, err
	}

	return updatedContactData, nil
}

// UpdateContact is to update a Contact record
func (d *DatabaseProvider) UpdateContact(newContactData entity.Contact) (entity.Contact, error) {
	currentContactData, err := d.getByCompanyName(newContactData.Company)
	if err != nil {
		return entity.Contact{}, err
	}
	rollbar.Info(fmt.Sprintf("Updating ContactID %s. \nOld Data: %v \nNew Data: %v", currentContactData.ID, currentContactData, newContactData))
	updatedContact := entity.Contact{
		ID:          currentContactData.ID,
		Created:     currentContactData.Created,
		FirstName:   newContactData.FirstName,
		LastName:    newContactData.LastName,
		Role:        newContactData.Role,
		Contactname: newContactData.Contactname,
		Password:    currentContactData.Password,
		Email:       newContactData.Email,
		Phone:       newContactData.Phone,
		TimeZone:    newContactData.TimeZone,
		DarkMode:    newContactData.DarkMode,
	}

	err = d.setByContactID(currentContactData.ID, updatedContact)
	if err != nil {
		return entity.Contact{}, err
	}
	updatedContactData, err := d.getByContactID(updatedContact.ID)
	if err != nil {
		return entity.Contact{}, err
	}

	return updatedContactData, nil
}

// DeleteContact is to update a Contact record
func (d *DatabaseProvider) DeleteContact(Contact entity.Contact) error {
	ContactData, err := d.getByContactname(Contact.Contactname)
	if err != nil {
		return err
	}

	err = d.deleteByContactID(ContactData.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted Contact %s", ContactData.Contactname))

	return nil
}

func (d *DatabaseProvider) addContact(ContactData entity.Contact) (entity.Contact, error) {
	existingContact, _ := d.GetByContactname(ContactData.Contactname)
	if (entity.Contact{}) != existingContact {
		return entity.Contact{}, fmt.Errorf("Error adding Contact %s: Contactname already exists", ContactData.Contactname)
	}
	newUUID := uuid.New().String()
	newContactData := entity.Contact{
		ID:          newUUID,
		Created:     time.Now().String(),
		FirstName:   ContactData.FirstName,
		LastName:    ContactData.LastName,
		Role:        ContactData.Role,
		Contactname: ContactData.Contactname,
		Password:    helper.Hash(ContactData.Password, newUUID),
		Email:       ContactData.Email,
		Phone:       ContactData.Phone,
		TimeZone:    ContactData.TimeZone,
		DarkMode:    ContactData.DarkMode,
	}
	addContactResult, err := d.Database.Collection("Contacts").Doc(newUUID).Set(context.TODO(), newContactData)
	if err != nil {
		return entity.Contact{}, fmt.Errorf("Error setting Contact %s by ID: %s", newContactData.Contactname, err)
	}
	rollbar.Info(fmt.Sprintf("Contact %s added at %s.", newContactData.Contactname, addContactResult))

	newContact, err := d.getByContactID(newContactData.ID)
	if err != nil {
		return entity.Contact{}, fmt.Errorf("Error getting newly created Contact %s by ID: %s", newContactData.Contactname, err)
	}

	return newContact, nil
}

func (d *DatabaseProvider) getByContactID(ContactID string) (entity.Contact, error) {
	var Contact entity.Contact
	ContactData, err := d.Database.Collection("Contacts").Doc(ContactID).Get(context.TODO())
	if err != nil {
		return entity.Contact{}, fmt.Errorf("Error getting Contact %s by ID: %s", ContactID, err)
	}
	ContactData.DataTo(&Contact)

	return Contact, nil
}

func (d *DatabaseProvider) setByContactID(ContactID string, ContactData entity.Contact) error {
	_, err := d.Database.Collection("Contacts").Doc(ContactID).Set(context.TODO(), ContactData)
	if err != nil {
		return fmt.Errorf("Error setting Contact %s by ID: %s", ContactID, err)
	}

	return nil
}

func (d *DatabaseProvider) deleteByContactID(ContactID string) error {
	result, err := d.Database.Collection("Contacts").Doc(ContactID).Delete(context.TODO())
	if err != nil {
		return fmt.Errorf("Error deleting Contact %s by ID: %s", ContactID, err)
	}
	log.Printf("Deleting Contact %s: %v", ContactID, result)

	return nil
}

// func (d *DatabaseProvider) getByContactname(Contactname string) (entity.Contact, error) {
// 	var Contact entity.Contact

// 	Contacts := d.Database.Collection("Contacts").Where("Contactname", "==", Contactname).Documents(context.TODO())
// 	allMatchingContacts, err := Contacts.GetAll()
// 	if err != nil {
// 		return entity.Contact{}, err
// 	}
// 	for _, fbContact := range allMatchingContacts {
// 		err = fbContact.DataTo(&Contact)
// 		if err != nil {
// 			return entity.Contact{}, fmt.Errorf("ERROR: Contact error - Firestore.DataTo() error %w, for Contact %s", err, Contactname)
// 		}
// 		return Contact, nil
// 		// data = append(data, fbContact.Data())
// 	}

// 	return entity.Contact{}, ErrContactNotFound
// }
