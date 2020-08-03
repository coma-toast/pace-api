package company

import (
	"errors"
	"fmt"
	"time"

	"github.com/coma-toast/pace-api/pkg/entity"
	"github.com/coma-toast/pace-api/pkg/provider/firestoredb"
	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

// DatabaseProvider is a company.Provider the uses a database
type DatabaseProvider struct {
	SharedProvider *firestoredb.DatabaseProvider
}

// ErrCompanyNotFound if no companies are found
var ErrCompanyNotFound = errors.New("Company not found")

// GetAll gets a Company by ID
func (d *DatabaseProvider) GetAll() ([]entity.Company, error) {
	var allCompanyData []entity.Company

	err := d.SharedProvider.GetAll(&allCompanyData)
	if err != nil {
		return []entity.Company{}, err
	}

	return allCompanyData, nil
}

// GetByName gets a Company by name
func (d *DatabaseProvider) GetByName(companyName string) (entity.Company, error) {
	var company entity.Company
	err := d.SharedProvider.GetFirstBy("Name", "==", companyName, &company)
	if err != nil {
		return entity.Company{}, err
	}

	return company, nil
}

// Add is to update a Company record
func (d *DatabaseProvider) Add(newCompanyData entity.Company) (entity.Company, error) {
	rollbar.Info(fmt.Sprintf("Adding new Company to DB %s", newCompanyData.Name))

	var existingCompany entity.Company
	err := d.SharedProvider.GetFirstBy("Name", "==", newCompanyData.Name, &existingCompany)
	if (entity.Company{}) != existingCompany {
		return entity.Company{}, fmt.Errorf("Error adding Company %s: Companyname already exists. ID: %s", newCompanyData.Name, existingCompany.ID)
	}

	newUUID := uuid.New().String()
	newCompanyData = entity.Company{
		ID:             newUUID,
		Created:        time.Now().String(),
		Name:           newCompanyData.Name,
		PrimaryContact: newCompanyData.PrimaryContact,
		Contacts:       newCompanyData.Contacts,
		Phone:          newCompanyData.Phone,
		Email:          newCompanyData.Email,
		Address:        newCompanyData.Address,
		City:           newCompanyData.City,
		State:          newCompanyData.State,
		Zip:            newCompanyData.Zip,
	}
	err = d.SharedProvider.Set(newCompanyData.ID, newCompanyData)
	if err != nil {
		return entity.Company{}, fmt.Errorf("Error setting Company %s by ID: %s", newCompanyData.Name, err)
	}

	var newCompany = entity.Company{}
	err = d.SharedProvider.GetByID(newCompanyData.ID, &newCompany)
	if err != nil {
		return entity.Company{}, fmt.Errorf("Error getting newly created Company %s by ID: %s", newCompanyData.Name, err)
	}

	rollbar.Info(fmt.Sprintf("Company %s added.", newCompanyData.Name))
	return newCompany, nil
}

// Update is to update a Company record
func (d *DatabaseProvider) Update(newCompanyData entity.Company) (entity.Company, error) {
	var currentCompanyData = entity.Company{}
	err := d.SharedProvider.GetByID(newCompanyData.ID, &currentCompanyData)
	if err != nil {
		return entity.Company{}, err
	}
	rollbar.Info(fmt.Sprintf("Updating CompanyID %s. \nOld Data: %v \nNew Data: %v", currentCompanyData.ID, currentCompanyData, newCompanyData))
	updatedCompany := entity.Company{
		ID:             currentCompanyData.ID,
		Created:        currentCompanyData.Created,
		Name:           newCompanyData.Name,
		PrimaryContact: newCompanyData.PrimaryContact,
		Contacts:       newCompanyData.Contacts,
		Phone:          newCompanyData.Phone,
		Email:          newCompanyData.Email,
		Address:        newCompanyData.Address,
		City:           newCompanyData.City,
		State:          newCompanyData.State,
		Zip:            newCompanyData.Zip,
	}

	err = d.SharedProvider.Set(currentCompanyData.ID, updatedCompany)
	if err != nil {
		return entity.Company{}, err
	}

	var updatedCompanyData = entity.Company{}
	err = d.SharedProvider.GetByID(updatedCompany.ID, &updatedCompanyData)
	if err != nil {
		return entity.Company{}, err
	}

	return updatedCompanyData, nil
}

// Delete is to update a Company record
func (d *DatabaseProvider) Delete(company entity.Company) error {
	var currentCompany entity.Company

	err := d.SharedProvider.GetByID(company.ID, &currentCompany)
	if (entity.Company{}) == currentCompany {
		return fmt.Errorf("Company not found")
	}

	err = d.SharedProvider.Delete(company.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted Company %s: %s", company.ID, company.Name))

	return nil
}
