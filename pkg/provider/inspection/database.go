package inspection

import (
	"errors"
	"fmt"
	"time"

	"github.com/coma-toast/pace-api/pkg/entity"
	"github.com/coma-toast/pace-api/pkg/provider/firestoredb"
	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

// DatabaseProvider is a inspection.Provider the uses a database
type DatabaseProvider struct {
	SharedProvider *firestoredb.DatabaseProvider
}

// ErrInspectionNotFound if no Inspections are found
var ErrInspectionNotFound = errors.New("Inspection not found")

// GetAll gets a Inspection by inspectionname
func (d *DatabaseProvider) GetAll() ([]entity.Inspection, error) {
	var inspections []entity.Inspection
	err := d.SharedProvider.GetAll(&inspections)
	if err != nil {
		return nil, err
	}

	return inspections, nil
}

// GetByID gets a Inspection by inspectionname
func (d *DatabaseProvider) GetByID(ID string) (entity.Inspection, error) {
	var inspection entity.Inspection
	err := d.SharedProvider.GetFirstBy("ID", "==", ID, &inspection)
	if err != nil {
		return entity.Inspection{}, fmt.Errorf("%s: %w", err, ErrInspectionNotFound)
	}

	return inspection, nil
}

// Add is to update a inspection record
func (d *DatabaseProvider) Add(inspectionData entity.UpdateInspectionRequest) (entity.Inspection, error) {
	rollbar.Info(fmt.Sprintf("Adding new Inspection to DB %s", inspectionData.ID))

	var existingInspection entity.Inspection
	err := d.SharedProvider.GetFirstBy("ID", "==", inspectionData.ID, &existingInspection)
	if (entity.Inspection{}) != existingInspection {
		return entity.Inspection{}, fmt.Errorf("Error adding inspection %s: ID already exists. ID: %s", inspectionData.ID, existingInspection.ID)
	}

	newUUID := uuid.New().String()
	newInspectionData := entity.Inspection{
		ID:             newUUID,
		Created:        time.Now().String(),
		ProjectID:      inspectionData.ProjectID,
		Username:       inspectionData.Username,
		StartTime:      inspectionData.StartTime,
		EndTime:        inspectionData.EndTime,
		InspectedParts: inspectionData.InspectedParts,
	}
	err = d.SharedProvider.Set(newInspectionData.ID, newInspectionData)
	if err != nil {
		return entity.Inspection{}, fmt.Errorf("Error setting inspection %s by ID: %s", newInspectionData.ID, err)
	}

	var newInspection = entity.Inspection{}
	err = d.SharedProvider.GetByID(newInspectionData.ID, &newInspection)
	if err != nil {
		return entity.Inspection{}, fmt.Errorf("Error getting newly created inspection %s by ID: %s", newInspectionData.ID, err)
	}

	rollbar.Info(fmt.Sprintf("Inspection %s added.", newInspectionData.ID))
	return newInspection, nil
}

// Update is to update a inspection record
func (d *DatabaseProvider) Update(newInspectionData entity.UpdateInspectionRequest) (entity.Inspection, error) {
	var currentInspectionData entity.Inspection
	err := d.SharedProvider.GetFirstBy("ID", "==", newInspectionData.ID, &currentInspectionData)
	if err != nil {
		return entity.Inspection{}, err
	}

	rollbar.Info(fmt.Sprintf("Updating inspectionID %s. \nOld Data: %v \nNew Data: %v", currentInspectionData.ID, currentInspectionData, newInspectionData))
	updatedInspection := entity.Inspection{
		ID:             currentInspectionData.ID,
		Created:        currentInspectionData.Created,
		ProjectID:      newInspectionData.ProjectID,
		Username:       newInspectionData.Username,
		StartTime:      newInspectionData.StartTime,
		EndTime:        newInspectionData.EndTime,
		InspectedParts: newInspectionData.InspectedParts,
	}

	err = d.SharedProvider.Set(currentInspectionData.ID, updatedInspection)
	if err != nil {
		return entity.Inspection{}, err
	}

	var updatedInspectionData = entity.Inspection{}
	err = d.SharedProvider.GetByID(currentInspectionData.ID, &updatedInspectionData)
	if err != nil {
		return entity.Inspection{}, err
	}

	rollbar.Info(fmt.Sprintf("Inspection %s updated.", updatedInspectionData.ID))

	return updatedInspectionData, nil
}

// Delete deletes an inspection
func (d *DatabaseProvider) Delete(inspection entity.Inspection) error {
	var currentInspection entity.Inspection

	err := d.SharedProvider.GetByID(inspection.ID, &currentInspection)
	if (entity.Inspection{}) == currentInspection {
		return fmt.Errorf("Inspection not found")
	}

	rollbar.Info(fmt.Sprintf("Deleting Inspection from DB: %s", inspection.ID))

	err = d.SharedProvider.Delete(inspection.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted inspection %s", inspection.ID))

	return nil
}
