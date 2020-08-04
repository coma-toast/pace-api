package inventory

import (
	"errors"
	"fmt"
	"time"

	"github.com/coma-toast/pace-api/pkg/entity"
	"github.com/coma-toast/pace-api/pkg/provider/firestoredb"
	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

// DatabaseProvider is a inventory.Provider the uses a database
type DatabaseProvider struct {
	SharedProvider *firestoredb.DatabaseProvider
}

// ErrInventoryNotFound if no Inventor is found
var ErrInventoryNotFound = errors.New("Inventory not found")

// GetAll gets all inventory
func (d *DatabaseProvider) GetAll() ([]entity.Inventory, error) {
	var inventory []entity.Inventory
	err := d.SharedProvider.GetAll(&inventory)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

// GetByID gets a Inventory by ID
func (d *DatabaseProvider) GetByID(ID string) (entity.Inventory, error) {
	var inventory entity.Inventory
	err := d.SharedProvider.GetFirstBy("ID", "==", ID, &inventory)
	if err != nil {
		return entity.Inventory{}, fmt.Errorf("%s: %w", err, ErrInventoryNotFound)
	}

	return inventory, nil
}

// Add is to update a inventory record
func (d *DatabaseProvider) Add(newInventoryData entity.Inventory) (entity.Inventory, error) {
	rollbar.Info(fmt.Sprintf("Adding new Inventory to DB %s", newInventoryData.ID))

	var existingInventory entity.Inventory
	err := d.SharedProvider.GetFirstBy("ID", "==", newInventoryData.ID, &existingInventory)
	if (entity.Inventory{}) != existingInventory {
		return entity.Inventory{}, fmt.Errorf("Error adding inventory %s: Inventoryname already exists. ID: %s", newInventoryData.ID, existingInventory.ID)
	}

	newUUID := uuid.New().String()
	newInventoryData = entity.Inventory{
		ID:        newUUID,
		Created:   time.Now().String(),
		ProjectID: newInventoryData.ProjectID,
		Stage:     newInventoryData.Stage,
		Size:      newInventoryData.Size,
		Length:    newInventoryData.Length,
		Grade:     newInventoryData.Grade,
		Shape:     newInventoryData.Shape,
		Passed:    newInventoryData.Passed,
		Sequence:  newInventoryData.Sequence,
		Priority:  newInventoryData.Priority,
	}
	err = d.SharedProvider.Set(newInventoryData.ID, newInventoryData)
	if err != nil {
		return entity.Inventory{}, fmt.Errorf("Error setting inventory %s by ID: %s", newInventoryData.ID, err)
	}

	var newInventory = entity.Inventory{}
	err = d.SharedProvider.GetByID(newInventoryData.ID, &newInventory)
	if err != nil {
		return entity.Inventory{}, fmt.Errorf("Error getting newly created inventory %s by ID: %s", newInventoryData.ID, err)
	}

	rollbar.Info(fmt.Sprintf("Inventory %s added.", newInventoryData.ID))
	return newInventory, nil
}

// Update is to update a inventory record
func (d *DatabaseProvider) Update(newInventoryData entity.UpdateInventoryRequest) (entity.Inventory, error) {
	var currentInventoryData entity.Inventory
	err := d.SharedProvider.GetFirstBy("ID", "==", newInventoryData.ID, &currentInventoryData)
	if err != nil {
		return entity.Inventory{}, err
	}

	rollbar.Info(fmt.Sprintf("Updating inventoryID %s. \nOld Data: %v \nNew Data: %v", currentInventoryData.ID, currentInventoryData, newInventoryData))
	updatedInventory := entity.Inventory{
		ID:        currentInventoryData.ID,
		Created:   currentInventoryData.Created,
		ProjectID: newInventoryData.ProjectID,
		Stage:     newInventoryData.Stage,
		Size:      newInventoryData.Size,
		Length:    newInventoryData.Length,
		Grade:     newInventoryData.Grade,
		Shape:     newInventoryData.Shape,
		Passed:    newInventoryData.Passed,
		Sequence:  newInventoryData.Sequence,
		Priority:  newInventoryData.Priority,
	}

	err = d.SharedProvider.Set(currentInventoryData.ID, updatedInventory)
	if err != nil {
		return entity.Inventory{}, err
	}

	var updatedInventoryData = entity.Inventory{}
	err = d.SharedProvider.GetByID(currentInventoryData.ID, &updatedInventoryData)
	if err != nil {
		return entity.Inventory{}, err
	}

	rollbar.Info(fmt.Sprintf("Inventory %s updated.", updatedInventoryData.ID))

	return updatedInventoryData, nil
}

// Delete deletes an inventory item
func (d *DatabaseProvider) Delete(inventory entity.Inventory) error {
	rollbar.Info(fmt.Sprintf("Deleting Inventory from DB: %s", inventory.ID))

	err := d.SharedProvider.Delete(inventory.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted inventory %s", inventory.ID))

	return nil
}
