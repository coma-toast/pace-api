package contact

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/coma-toast/pace-api/pkg/entity"
	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

// DatabaseProvider is a contact.Provider the uses a database
type DatabaseProvider struct {
	Database *firestore.Client
}

// ErrContactNotFound if no Contacts are found
var ErrContactNotFound = errors.New("Contact not found")

// // GetByID gets a Contact by ID
// func (d *DatabaseProvider) GetByID(ID string) (entity.Contact, error) {
// 	return d.getByID(ID)
// }

// GetAll gets a Contact by ID
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
	currentContactData, err := d.getByContactID(newContactData.ID)
	// * dev code currentContactData, err := d.getByNameAndCompany(newContactData.FirstName, newContactData.LastName, newContactData.Company)
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
func (d *DatabaseProvider) DeleteContact(contact entity.Contact) error {
	contactData, err := d.getByContactID(contact.ID)
	if err != nil {
		return err
	}

	err = d.deleteByContactID(contactData.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted Contact %s: %s %s", contactData.ID, contactData.FirstName, contactData.LastName))

	return nil
}

func (d *DatabaseProvider) addContact(contactData entity.Contact) (entity.Contact, error) {
	existingContact, _ := d.getByNameAndCompany(contactData.FirstName, contactData.LastName, contactData.Company)
	if (entity.Contact{}) != existingContact {
		return entity.Contact{}, fmt.Errorf("Error adding Contact %s: ID already exists", contactData.ID)
	}
	newUUID := uuid.New().String()
	newContactData := entity.Contact{
		ID:        newUUID,
		Created:   time.Now().String(),
		FirstName: contactData.FirstName,
		LastName:  contactData.LastName,
		Company:   contactData.Company,
		Email:     contactData.Email,
		Phone:     contactData.Phone,
		Timezone:  contactData.Timezone,
	}
	addContactResult, err := d.Database.Collection("contacts").Doc(newUUID).Set(context.TODO(), newContactData)
	if err != nil {
		return entity.Contact{}, fmt.Errorf("Error setting Contact %s by ID: %s", newContactData.ID, err)
	}
	rollbar.Info(fmt.Sprintf("Contact %s added at %s.", newContactData.ID, addContactResult))

	newContact, err := d.getByContactID(newContactData.ID)
	if err != nil {
		return entity.Contact{}, fmt.Errorf("Error getting newly created Contact %s by ID: %s", newContactData.ID, err)
	}

	return newContact, nil
}

func (d *DatabaseProvider) getByContactID(contactID string) (entity.Contact, error) {
	var contact entity.Contact

	contactData, err := d.Database.Collection("contacts").Doc(contactID).Get(context.TODO())
	if err != nil {
		return entity.Contact{}, fmt.Errorf("Error getting Contact %s by ID: %s", contactID, err)
	}
	contactData.DataTo(&contact)

	return contact, nil
}

func (d *DatabaseProvider) getByNameAndCompany(firstName string, lastName string, company string) (entity.Contact, error) {
	var contact entity.Contact
	contactSnapshot, err := d.Database.Collection("contacts").Where("Company", "==", company).Documents(context.TODO()).GetAll()
	if err != nil {
		return entity.Contact{}, fmt.Errorf("Error getting contact by name and company: %s", err)
	}

	for _, companyContact := range contactSnapshot {
		companyContact.DataTo(&contact)
		if contact.FirstName == firstName && contact.LastName == lastName {
			return contact, nil
		}
	}

	return entity.Contact{}, ErrContactNotFound
}

func (d *DatabaseProvider) setByContactID(ContactID string, ContactData entity.Contact) error {
	_, err := d.Database.Collection("contacts").Doc(ContactID).Set(context.TODO(), ContactData)
	if err != nil {
		return fmt.Errorf("Error setting Contact %s by ID: %s", ContactID, err)
	}

	return nil
}

func (d *DatabaseProvider) deleteByContactID(ContactID string) error {
	result, err := d.Database.Collection("contacts").Doc(ContactID).Delete(context.TODO())
	if err != nil {
		return fmt.Errorf("Error deleting Contact %s by ID: %s", ContactID, err)
	}
	log.Printf("Deleting Contact %s: %v", ContactID, result)

	return nil
}

// func (d *DatabaseProvider) getByID(ID string) (entity.Contact, error) {
// 	var Contact entity.Contact

// 	Contacts := d.Database.Collection("contacts").Where("ID", "==", ID).Documents(context.TODO())
// 	allMatchingContacts, err := Contacts.GetAll()
// 	if err != nil {
// 		return entity.Contact{}, err
// 	}
// 	for _, fbContact := range allMatchingContacts {
// 		err = fbContact.DataTo(&Contact)
// 		if err != nil {
// 			return entity.Contact{}, fmt.Errorf("ERROR: Contact error - Firestore.DataTo() error %w, for Contact %s", err, ID)
// 		}
// 		return Contact, nil
// 		// data = append(data, fbContact.Data())
// 	}

// 	return entity.Contact{}, ErrContactNotFound
// }
