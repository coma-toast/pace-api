package firestore

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

// DatabaseProvider is a firestore.Provider the uses a database
type DatabaseProvider struct {
	Database *firestore.Client
}

// ErrFirestoreNotFound if no Firestores are found
var ErrFirestoreNotFound = errors.New("Firestore not found")

// // GetByID gets a Firestore by ID
// func (d *DatabaseProvider) GetByID(ID string) (entity.Firestore, error) {
// 	return d.getByID(ID)
// }

// GetAll gets a Firestore by ID
func (d *DatabaseProvider) GetAll() ([]entity.Firestore, error) {
	return d.getAll()
}

func (d *DatabaseProvider) getAll() ([]entity.Firestore, error) {
	var firestores []entity.Firestore

	allFirestoreData, err := d.Database.Collection("firestores").Documents(context.TODO()).GetAll()
	if err != nil {
		return []entity.Firestore{}, err
	}

	for _, FirestoreData := range allFirestoreData {
		var Firestore entity.Firestore
		err := FirestoreData.DataTo(&Firestore)
		if err != nil {
			return []entity.Firestore{}, fmt.Errorf("ERROR: GetAll(): Firestore.DataTo() error %w", err)
		}
		firestores = append(firestores, Firestore)
	}

	return firestores, nil
}

// AddFirestore is to update a Firestore record
func (d *DatabaseProvider) AddFirestore(newFirestoreData entity.Firestore) (entity.Firestore, error) {
	firestoreRef, err := d.addFirestore(newFirestoreData)
	if err != nil {
		return entity.Firestore{}, err
	}
	rollbar.Info(fmt.Sprintf("Adding new Firestore %s %s", newFirestoreData.FirstName, newFirestoreData.LastName))
	updatedFirestoreData, err := d.getByFirestoreID(firestoreRef.ID)
	if err != nil {
		return entity.Firestore{}, err
	}

	return updatedFirestoreData, nil
}

// UpdateFirestore is to update a Firestore record
func (d *DatabaseProvider) UpdateFirestore(newFirestoreData entity.Firestore) (entity.Firestore, error) {
	currentFirestoreData, err := d.getByFirestoreID(newFirestoreData.ID)
	// * dev code currentFirestoreData, err := d.getByNameAndCompany(newFirestoreData.FirstName, newFirestoreData.LastName, newFirestoreData.Company)
	if err != nil {
		return entity.Firestore{}, err
	}
	rollbar.Info(fmt.Sprintf("Updating FirestoreID %s. \nOld Data: %v \nNew Data: %v", currentFirestoreData.ID, currentFirestoreData, newFirestoreData))
	updatedFirestore := entity.Firestore{
		ID:        currentFirestoreData.ID,
		Created:   currentFirestoreData.Created,
		FirstName: newFirestoreData.FirstName,
		LastName:  newFirestoreData.LastName,
		Company:   newFirestoreData.Company,
		Email:     newFirestoreData.Email,
		Phone:     newFirestoreData.Phone,
		Timezone:  newFirestoreData.Timezone,
	}

	err = d.setByFirestoreID(currentFirestoreData.ID, updatedFirestore)
	if err != nil {
		return entity.Firestore{}, err
	}
	updatedFirestoreData, err := d.getByFirestoreID(updatedFirestore.ID)
	if err != nil {
		return entity.Firestore{}, err
	}

	return updatedFirestoreData, nil
}

// DeleteFirestore is to update a Firestore record
func (d *DatabaseProvider) DeleteFirestore(firestore entity.Firestore) error {
	firestoreData, err := d.getByFirestoreID(firestore.ID)
	if err != nil {
		return err
	}

	err = d.deleteByFirestoreID(firestoreData.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted Firestore %s: %s %s", firestoreData.ID, firestoreData.FirstName, firestoreData.LastName))

	return nil
}

func (d *DatabaseProvider) addFirestore(firestoreData entity.Firestore) (entity.Firestore, error) {
	existingFirestore, _ := d.getByNameAndCompany(firestoreData.FirstName, firestoreData.LastName, firestoreData.Company)
	if (entity.Firestore{}) != existingFirestore {
		return entity.Firestore{}, fmt.Errorf("Error adding Firestore %s: ID already exists", firestoreData.ID)
	}
	newUUID := uuid.New().String()
	newFirestoreData := entity.Firestore{
		ID:        newUUID,
		Created:   time.Now().String(),
		FirstName: firestoreData.FirstName,
		LastName:  firestoreData.LastName,
		Company:   firestoreData.Company,
		Email:     firestoreData.Email,
		Phone:     firestoreData.Phone,
		Timezone:  firestoreData.Timezone,
	}
	addFirestoreResult, err := d.Database.Collection("firestores").Doc(newUUID).Set(context.TODO(), newFirestoreData)
	if err != nil {
		return entity.Firestore{}, fmt.Errorf("Error setting Firestore %s by ID: %s", newFirestoreData.ID, err)
	}
	rollbar.Info(fmt.Sprintf("Firestore %s added at %s.", newFirestoreData.ID, addFirestoreResult))

	newFirestore, err := d.getByFirestoreID(newFirestoreData.ID)
	if err != nil {
		return entity.Firestore{}, fmt.Errorf("Error getting newly created Firestore %s by ID: %s", newFirestoreData.ID, err)
	}

	return newFirestore, nil
}

func (d *DatabaseProvider) getByFirestoreID(firestoreID string) (entity.Firestore, error) {
	var firestore entity.Firestore

	firestoreData, err := d.Database.Collection("firestores").Doc(firestoreID).Get(context.TODO())
	if err != nil {
		return entity.Firestore{}, fmt.Errorf("Error getting Firestore %s by ID: %s", firestoreID, err)
	}
	firestoreData.DataTo(&firestore)

	return firestore, nil
}

func (d *DatabaseProvider) getByNameAndCompany(firstName string, lastName string, company string) (entity.Firestore, error) {
	var firestore entity.Firestore
	firestoreSnapshot, err := d.Database.Collection("firestores").Where("Company", "==", company).Documents(context.TODO()).GetAll()
	if err != nil {
		return entity.Firestore{}, fmt.Errorf("Error getting firestore by name and company: %s", err)
	}

	for _, companyFirestore := range firestoreSnapshot {
		companyFirestore.DataTo(&firestore)
		if firestore.FirstName == firstName && firestore.LastName == lastName {
			return firestore, nil
		}
	}

	return entity.Firestore{}, ErrFirestoreNotFound
}

func (d *DatabaseProvider) setByFirestoreID(FirestoreID string, FirestoreData entity.Firestore) error {
	_, err := d.Database.Collection("firestores").Doc(FirestoreID).Set(context.TODO(), FirestoreData)
	if err != nil {
		return fmt.Errorf("Error setting Firestore %s by ID: %s", FirestoreID, err)
	}

	return nil
}

func (d *DatabaseProvider) deleteByFirestoreID(FirestoreID string) error {
	result, err := d.Database.Collection("firestores").Doc(FirestoreID).Delete(context.TODO())
	if err != nil {
		return fmt.Errorf("Error deleting Firestore %s by ID: %s", FirestoreID, err)
	}
	log.Printf("Deleting Firestore %s: %v", FirestoreID, result)

	return nil
}

// func (d *DatabaseProvider) getByID(ID string) (entity.Firestore, error) {
// 	var Firestore entity.Firestore

// 	Firestores := d.Database.Collection("firestores").Where("ID", "==", ID).Documents(context.TODO())
// 	allMatchingFirestores, err := Firestores.GetAll()
// 	if err != nil {
// 		return entity.Firestore{}, err
// 	}
// 	for _, fbFirestore := range allMatchingFirestores {
// 		err = fbFirestore.DataTo(&Firestore)
// 		if err != nil {
// 			return entity.Firestore{}, fmt.Errorf("ERROR: Firestore error - Firestore.DataTo() error %w, for Firestore %s", err, ID)
// 		}
// 		return Firestore, nil
// 		// data = append(data, fbFirestore.Data())
// 	}

// 	return entity.Firestore{}, ErrFirestoreNotFound
// }
