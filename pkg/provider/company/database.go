package company

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

// DatabaseProvider is a company.Provider the uses a database
type DatabaseProvider struct {
	Database *firestore.Client
}

// ErrCompanyNotFound if no companies are found
var ErrCompanyNotFound = errors.New("Company not found")

// // GetByID gets a Company by ID
// func (d *DatabaseProvider) GetByID(ID string) (entity.Company, error) {
// 	return d.getByID(ID)
// }

// GetAll gets a Company by ID
func (d *DatabaseProvider) GetAll() ([]entity.Company, error) {
	return d.getAll()
}

func (d *DatabaseProvider) getAll() ([]entity.Company, error) {
	var company []entity.Company

	allCompanyData, err := d.Database.Collection("company").Documents(context.TODO()).GetAll()
	if err != nil {
		return []entity.Company{}, err
	}

	for _, CompanyData := range allCompanyData {
		var Company entity.Company
		err := CompanyData.DataTo(&Company)
		if err != nil {
			return []entity.Company{}, fmt.Errorf("ERROR: GetAll(): Firestore.DataTo() error %w", err)
		}
		company = append(company, Company)
	}

	return company, nil
}

// AddCompany is to update a Company record
func (d *DatabaseProvider) AddCompany(newCompanyData entity.Company) (entity.Company, error) {
	companyRef, err := d.addCompany(newCompanyData)
	if err != nil {
		return entity.Company{}, err
	}
	rollbar.Info(fmt.Sprintf("Adding new Company %s", newCompanyData.Name))
	updatedCompanyData, err := d.getByCompanyID(companyRef.ID)
	if err != nil {
		return entity.Company{}, err
	}

	return updatedCompanyData, nil
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
func (d *DatabaseProvider) DeleteCompany(Company entity.Company) error {
	CompanyData, err := d.getByCompanyID(Company.ID)
	if err != nil {
		return err
	}

	err = d.deleteByCompanyID(CompanyData.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted Company %s: %s %s", CompanyData.ID, CompanyData.Name))

	return nil
}

func (d *DatabaseProvider) addCompany(CompanyData entity.Company) (entity.Company, error) {
	existingCompany, _ := d.getByName(CompanyData.Name)
	if (entity.Company{}) != existingCompany {
		return entity.Company{}, fmt.Errorf("Error adding Company %s: ID already exists", CompanyData.ID)
	}
	newUUID := uuid.New().String()
	newCompanyData := entity.Company{
		ID:             newUUID,
		Created:        time.Now().String(),
		Name:           CompanyData.Name,
		PrimaryContact: CompanyData.PrimaryContact,
		Contacts:       CompanyData.Contacts,
		Phone:          CompanyData.Phone,
		Email:          CompanyData.Email,
		Address:        CompanyData.Address,
		City:           CompanyData.City,
		State:          CompanyData.State,
		Zip:            CompanyData.Zip,
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

func (d *DatabaseProvider) setByCompanyID(CompanyID string, CompanyData entity.Company) error {
	_, err := d.Database.Collection("company").Doc(CompanyID).Set(context.TODO(), CompanyData)
	if err != nil {
		return fmt.Errorf("Error setting Company %s by ID: %s", CompanyID, err)
	}

	return nil
}

func (d *DatabaseProvider) deleteByCompanyID(CompanyID string) error {
	result, err := d.Database.Collection("company").Doc(CompanyID).Delete(context.TODO())
	if err != nil {
		return fmt.Errorf("Error deleting Company %s by ID: %s", CompanyID, err)
	}
	log.Printf("Deleting Company %s: %v", CompanyID, result)

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
