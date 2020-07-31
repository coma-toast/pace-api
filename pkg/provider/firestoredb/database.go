package firestoredb

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
)

// DatabaseProvider is a firestore.Provider the uses a database
type DatabaseProvider struct {
	Database   *firestore.Client
	Collection string
}

// ErrFirestoreNotFound if no Firestores are found
var ErrFirestoreNotFound = errors.New("Firestore Item not found")

// GetAll gets all items in a Firestore collection
func (d *DatabaseProvider) GetAll(target interface{}) error {
	returnData := make([]interface{}, 0)
	allFirestoreData, err := d.Database.Collection(d.Collection).Documents(context.TODO()).GetAll()
	// test, err := d.Database.Collection(d.Collection).Documents(context.TODO()).GetAll()
	// for _, testData := range test {
	// 	fmt.Println(testData.Data())
	// }
	if err != nil {
		return fmt.Errorf("Error getting collection: %w", err)
	}
	for _, firestoreData := range allFirestoreData {
		data := make(map[string]interface{})
		err := firestoreData.DataTo(&data)
		returnData = append(returnData, data)
		if err != nil {
			return fmt.Errorf("ERROR: GetAll(): Firestore.DataTo() error %w", err)
		}

	}

	mapstructure.Decode(returnData, target)

	return nil
	// return ErrFirestoreNotFound
}

// GetByID gets an item by ID
func (d *DatabaseProvider) GetByID(ID string, target interface{}) error {
	firestoreData, err := d.Database.Collection(d.Collection).Doc(ID).Get(context.TODO())
	if err != nil {
		return fmt.Errorf("Error getting %s with ID %s: %w", d.Collection, ID, err)
	}

	err = firestoreData.DataTo(target)
	if err != nil {
		return fmt.Errorf("ERROR: GetAll(): Firestore.DataTo() error %w", err)

	}

	return ErrFirestoreNotFound
}

// GetFirstBy gets the first returned item by a path, operator and value
func (d *DatabaseProvider) GetFirstBy(path string, op string, value string, target interface{}) error {
	allFirestoreData, err := d.Database.Collection(d.Collection).Where(path, op, value).Documents(context.TODO()).GetAll()
	if err != nil {
		return fmt.Errorf("Error getting collection: ", err)
	}

	for _, firestoreData := range allFirestoreData {
		err := firestoreData.DataTo(target)
		if err != nil {
			return fmt.Errorf("ERROR: GetAll(): Firestore.DataTo() error %w", err)
		}

		return nil
	}

	return ErrFirestoreNotFound
}

// Set is to add a Firestore record
func (d *DatabaseProvider) Set(ID string, data interface{}) error {
	_, err := d.Database.Collection(d.Collection).Doc(ID).Set(context.TODO(), data)
	if err != nil {
		return fmt.Errorf("Error getting %s with ID %s: %w", d.Collection, ID, err)
	}

	return nil
}

// Delete is to delete a record
func (d *DatabaseProvider) Delete(ID string) error {
	_, err := d.Database.Collection(d.Collection).Doc(ID).Delete(context.TODO())
	if err != nil {
		return fmt.Errorf("Error deleting %s with ID %s: %w", d.Collection, ID, err)
	}

	return nil
}

// * Do we want to use .Update to keep existing data? Or just pull the data and then use .Set with old+new data?
// Update is to update a Firestore record
// func (d *DatabaseProvider) Update(ID string, data interface{}) error {
// 	err = d.Database.Collection(d.Collection).Doc(ID).Update(context.TODO(), []firestore.Update{data})
// 	if err != nil {
// 		return entity.Firestore{}, err
// 	}
// 	updatedFirestoreData, err := d.getByFirestoreID(updatedFirestore.ID)
// 	if err != nil {
// 		return entity.Firestore{}, err
// 	}

// 	return updatedFirestoreData, nil
// }

// // DeleteFirestore is to update a Firestore record
// func (d *DatabaseProvider) DeleteFirestore(firestore entity.Firestore) error {
// 	firestoreData, err := d.getByFirestoreID(firestore.ID)
// 	if err != nil {
// 		return err
// 	}

// 	err = d.deleteByFirestoreID(firestoreData.ID)
// 	if err != nil {
// 		return err
// 	}
// 	rollbar.Info(fmt.Sprintf("Deleted Firestore %s: %s %s", firestoreData.ID, firestoreData.FirstName, firestoreData.LastName))

// 	return nil
// }

// func (d *DatabaseProvider) addFirestore(firestoreData entity.Firestore) (entity.Firestore, error) {
// 	existingFirestore, _ := d.getByNameAndCompany(firestoreData.FirstName, firestoreData.LastName, firestoreData.Company)
// 	if (entity.Firestore{}) != existingFirestore {
// 		return entity.Firestore{}, fmt.Errorf("Error adding Firestore %s: ID already exists", firestoreData.ID)
// 	}
// 	newUUID := uuid.New().String()
// 	newFirestoreData := entity.Firestore{
// 		ID:        newUUID,
// 		Created:   time.Now().String(),
// 		FirstName: firestoreData.FirstName,
// 		LastName:  firestoreData.LastName,
// 		Company:   firestoreData.Company,
// 		Email:     firestoreData.Email,
// 		Phone:     firestoreData.Phone,
// 		Timezone:  firestoreData.Timezone,
// 	}
// 	addFirestoreResult, err := d.Database.Collection("firestores").Doc(newUUID).Set(context.TODO(), newFirestoreData)
// 	if err != nil {
// 		return entity.Firestore{}, fmt.Errorf("Error setting Firestore %s by ID: %s", newFirestoreData.ID, err)
// 	}
// 	rollbar.Info(fmt.Sprintf("Firestore %s added at %s.", newFirestoreData.ID, addFirestoreResult))

// 	newFirestore, err := d.getByFirestoreID(newFirestoreData.ID)
// 	if err != nil {
// 		return entity.Firestore{}, fmt.Errorf("Error getting newly created Firestore %s by ID: %s", newFirestoreData.ID, err)
// 	}

// 	return newFirestore, nil
// }

// func (d *DatabaseProvider) getByFirestoreID(firestoreID string) (entity.Firestore, error) {
// 	var firestore entity.Firestore

// 	firestoreData, err := d.Database.Collection("firestores").Doc(firestoreID).Get(context.TODO())
// 	if err != nil {
// 		return entity.Firestore{}, fmt.Errorf("Error getting Firestore %s by ID: %s", firestoreID, err)
// 	}
// 	firestoreData.DataTo(&firestore)

// 	return firestore, nil
// }

// func (d *DatabaseProvider) getByNameAndCompany(firstName string, lastName string, company string) (entity.Firestore, error) {
// 	var firestore entity.Firestore
// 	firestoreSnapshot, err := d.Database.Collection("firestores").Where("Company", "==", company).Documents(context.TODO()).GetAll()
// 	if err != nil {
// 		return entity.Firestore{}, fmt.Errorf("Error getting firestore by name and company: %s", err)
// 	}

// 	for _, companyFirestore := range firestoreSnapshot {
// 		companyFirestore.DataTo(&firestore)
// 		if firestore.FirstName == firstName && firestore.LastName == lastName {
// 			return firestore, nil
// 		}
// 	}

// 	return entity.Firestore{}, ErrFirestoreNotFound
// }

// func (d *DatabaseProvider) setByFirestoreID(FirestoreID string, FirestoreData entity.Firestore) error {
// 	_, err := d.Database.Collection("firestores").Doc(FirestoreID).Set(context.TODO(), FirestoreData)
// 	if err != nil {
// 		return fmt.Errorf("Error setting Firestore %s by ID: %s", FirestoreID, err)
// 	}

// 	return nil
// }

// func (d *DatabaseProvider) deleteByFirestoreID(FirestoreID string) error {
// 	result, err := d.Database.Collection("firestores").Doc(FirestoreID).Delete(context.TODO())
// 	if err != nil {
// 		return fmt.Errorf("Error deleting Firestore %s by ID: %s", FirestoreID, err)
// 	}
// 	log.Printf("Deleting Firestore %s: %v", FirestoreID, result)

// 	return nil
// }

// // func (d *DatabaseProvider) getByID(ID string) (entity.Firestore, error) {
// // 	var Firestore entity.Firestore

// // 	Firestores := d.Database.Collection("firestores").Where("ID", "==", ID).Documents(context.TODO())
// // 	allMatchingFirestores, err := Firestores.GetAll()
// // 	if err != nil {
// // 		return entity.Firestore{}, err
// // 	}
// // 	for _, fbFirestore := range allMatchingFirestores {
// // 		err = fbFirestore.DataTo(&Firestore)
// // 		if err != nil {
// // 			return entity.Firestore{}, fmt.Errorf("ERROR: Firestore error - Firestore.DataTo() error %w, for Firestore %s", err, ID)
// // 		}
// // 		return Firestore, nil
// // 		// data = append(data, fbFirestore.Data())
// // 	}

// // 	return entity.Firestore{}, ErrFirestoreNotFound
// // }
