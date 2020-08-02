package company

import (
	"context"
	"errors"
	"fmt"
	"log"
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

// // GetByID gets a Company by ID
// func (d *DatabaseProvider) GetByID(ID string) (entity.Company, error) {
// 	return d.getByID(ID)
// }

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

// AddCompany is to update a Company record
func (d *DatabaseProvider) AddCompany(newCompanyData entity.Company) (entity.Company, error) {
	rollbar.Info(fmt.Sprintf("Adding new Company to DB %s", newCompanyData.Name))

	existingCompany, err := d.GetByName(newCompanyData.Name)
	if err == nil {
		// if (entity.Company{}) != existingCompany {
		return entity.Company{}, fmt.Errorf("Error adding Company %s: Companyname already exists. ID: %s", newCompanyData.Companyname, existingCompany.ID)
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

// UpdateCompany is to update a Company record
func (d *DatabaseProvider) UpdateCompany(newCompanyData entity.Company) (entity.Company, error) {
	currentCompanyData, err := d.getByCompanyID(newCompanyData.ID)
	// * dev code currentCompanyData, err := d.getByName(newCompanyData.Name, n, n)
	if err != nil {
		return entity.Company{}, err
	}
	rollbar.Info(fmt.Sprintf("Updating CompanyID %s. \nOld Data: %v \nNew Data: %v", currentCompanyData.ID, currentCompanyData, newCompanyData))
	updatedCompany := entity.Company{
		ID:             currentCompanyData.ID,
		Name:           newCompanyData.Name,
		PrimaryContact: newCompanyData.PrimaryContact,
		Contacts:       newCompanyData.Contacts,
		Phone:          newCompanyData.Phone,
		Email:          newCompanyData.Email,
		Created:        currentCompanyData.Created,
		Address:        newCompanyData.Address,
		City:           newCompanyData.City,
		State:          newCompanyData.State,
		Zip:            newCompanyData.Zip,
	}

	err = d.setByCompanyID(currentCompanyData.ID, updatedCompany)
	if err != nil {
		return entity.Company{}, err
	}
	updatedCompanyData, err := d.getByCompanyID(updatedCompany.ID)
	if err != nil {
		return entity.Company{}, err
	}

	return updatedCompanyData, nil
}

// DeleteCompany is to update a Company record
func (d *DatabaseProvider) DeleteCompany(company entity.Company) error {
	err := d.deleteByCompanyID(company.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted Company %s: %s", company.ID, company.Name))

	return nil
}

func (d *DatabaseProvider) addCompany(companyData entity.Company) (entity.Company, error) {
	existingCompany, _ := d.getByName(companyData.Name)
	if (existingCompany.ID) != "" {
		return entity.Company{}, fmt.Errorf("Error adding Company %s: ID already exists", companyData.ID)
	}
	newUUID := uuid.New().String()
	newCompanyData := entity.Company{
		ID:             newUUID,
		Created:        time.Now().String(),
		Name:           companyData.Name,
		PrimaryContact: companyData.PrimaryContact,
		Contacts:       companyData.Contacts,
		Phone:          companyData.Phone,
		Email:          companyData.Email,
		Address:        companyData.Address,
		City:           companyData.City,
		State:          companyData.State,
		Zip:            companyData.Zip,
	}
	addCompanyResult, err := d.Database.Collection("company").Doc(newUUID).Set(context.TODO(), newCompanyData)
	if err != nil {
		return entity.Company{}, fmt.Errorf("Error setting Company %s by ID: %s", newCompanyData.ID, err)
	}
	rollbar.Info(fmt.Sprintf("Company %s added at %s.", newCompanyData.ID, addCompanyResult))

	newCompany, err := d.getByCompanyID(newCompanyData.ID)
	if err != nil {
		return entity.Company{}, fmt.Errorf("Error getting newly created Company %s by ID: %s", newCompanyData.ID, err)
	}

	return newCompany, nil
}

func (d *DatabaseProvider) getByCompanyID(companyID string) (entity.Company, error) {
	var company entity.Company

	companyData, err := d.Database.Collection("company").Doc(companyID).Get(context.TODO())
	if err != nil {
		return entity.Company{}, fmt.Errorf("Error getting Company %s by ID: %s", companyID, err)
	}
	companyData.DataTo(&company)

	return company, nil
}

func (d *DatabaseProvider) getByName(name string) (entity.Company, error) {
	var company entity.Company
	companySnapshot, err := d.Database.Collection("company").Where("Company", "==", company).Documents(context.TODO()).GetAll()
	if err != nil {
		return entity.Company{}, fmt.Errorf("Error getting company by name and company: %s", err)
	}

	for _, companyCompany := range companySnapshot {
		companyCompany.DataTo(&company)
		return company, nil
	}

	return entity.Company{}, ErrCompanyNotFound
}

func (d *DatabaseProvider) setByCompanyID(companyID string, CompanyData entity.Company) error {
	_, err := d.Database.Collection("company").Doc(companyID).Set(context.TODO(), CompanyData)
	if err != nil {
		return fmt.Errorf("Error setting Company %s by ID: %s", companyID, err)
	}

	return nil
}

func (d *DatabaseProvider) deleteByCompanyID(companyID string) error {
	currentCompanyData, err := d.getByCompanyID(companyID)
	if err != nil {
		return ErrCompanyNotFound
	}

	result, err := d.Database.Collection("company").Doc(companyID).Delete(context.TODO())
	if err != nil {
		return fmt.Errorf("Error deleting Company %s (%s) by ID: %s", currentCompanyData.Name, companyID, err)
	}
	log.Printf("Deleting Company %s (%s): %v", currentCompanyData.Name, companyID, result)

	return nil
}

// func (d *DatabaseProvider) getByID(ID string) (entity.Company, error) {
// 	var Company entity.Company

// 	company := d.Database.Collection("company").Where("ID", "==", ID).Documents(context.TODO())
// 	allMatchingcompany, err := company.GetAll()
// 	if err != nil {
// 		return entity.Company{}, err
// 	}
// 	for _, fbCompany := range allMatchingcompany {
// 		err = fbCompany.DataTo(&Company)
// 		if err != nil {
// 			return entity.Company{}, fmt.Errorf("ERROR: Company error - Firestore.DataTo() error %w, for Company %s", err, ID)
// 		}
// 		return Company, nil
// 		// data = append(data, fbCompany.Data())
// 	}

// 	return entity.Company{}, ErrCompanyNotFound
// }
