package project

import (
	"errors"
	"fmt"
	"time"

	"github.com/coma-toast/pace-api/pkg/entity"
	"github.com/coma-toast/pace-api/pkg/provider/firestoredb"
	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

// DatabaseProvider is a project.Provider the uses a database
type DatabaseProvider struct {
	SharedProvider *firestoredb.DatabaseProvider
}

// ErrProjectNotFound if no Projects are found
var ErrProjectNotFound = errors.New("Project not found")

// GetAll gets a Project by projectname
func (d *DatabaseProvider) GetAll() ([]entity.Project, error) {
	var projects []entity.Project
	err := d.SharedProvider.GetAll(&projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// GetByName gets a Project by projectname
func (d *DatabaseProvider) GetByName(projectname string) (entity.Project, error) {
	var project entity.Project
	err := d.SharedProvider.GetFirstBy("Name", "==", projectname, &project)
	if err != nil {
		return entity.Project{}, fmt.Errorf("%s: %w", err, ErrProjectNotFound)
	}

	return project, nil
}

// Add is to update a project record
func (d *DatabaseProvider) Add(newProjectData entity.Project) (entity.Project, error) {
	rollbar.Info(fmt.Sprintf("Adding new Project to DB %s - %s", newProjectData.Name, newProjectData.Name))

	var existingProject entity.Project
	err := d.SharedProvider.GetFirstBy("Name", "==", newProjectData.Name, &existingProject)
	if (entity.Project{}) != existingProject {
		return entity.Project{}, fmt.Errorf("Error adding project %s: Projectname already exists. ID: %s", newProjectData.Name, existingProject.ID)
	}

	newUUID := uuid.New().String()
	newProjectData = entity.Project{
		ID:                    newUUID,
		Created:               time.Now().Format(time.RFC3339),
		Name:                  newProjectData.Name,
		StartDate:             newProjectData.StartDate,
		DueDate:               newProjectData.DueDate,
		Address:               newProjectData.Address,
		City:                  newProjectData.City,
		State:                 newProjectData.State,
		Zip:                   newProjectData.Zip,
		ClientID:              newProjectData.ClientID,
		EORNameID:             newProjectData.EORNameID,
		DetailerNameID:        newProjectData.DetailerNameID,
		InspectionLabID:       newProjectData.InspectionLabID,
		SteelErectorNameID:    newProjectData.SteelErectorNameID,
		SteelFabricatorNameID: newProjectData.SteelFabricatorNameID,
		GeneralContractorID:   newProjectData.GeneralContractorID,
		PrimaryContactNameID:  newProjectData.PrimaryContactNameID,
		PrimaryContactPhone:   newProjectData.PrimaryContactPhone,
		PrimaryContactEmail:   newProjectData.PrimaryContactEmail,
		SquareFootage:         newProjectData.SquareFootage,
		WeightInTons:          newProjectData.WeightInTons,
	}
	err = d.SharedProvider.Set(newProjectData.ID, newProjectData)
	if err != nil {
		return entity.Project{}, fmt.Errorf("Error setting project %s by ID: %s", newProjectData.Name, err)
	}

	var newProject = entity.Project{}
	err = d.SharedProvider.GetByID(newProjectData.ID, &newProject)
	if err != nil {
		return entity.Project{}, fmt.Errorf("Error getting newly created project %s by ID: %s", newProjectData.Name, err)
	}

	rollbar.Info(fmt.Sprintf("Project %s added.", newProjectData.Name))
	return newProject, nil
}

// Update is to update a project record
func (d *DatabaseProvider) Update(newProjectData entity.UpdateProjectRequest) (entity.Project, error) {
	var currentProjectData entity.Project
	err := d.SharedProvider.GetFirstBy("Name", "==", newProjectData.Name, &currentProjectData)
	if err != nil {
		return entity.Project{}, err
	}

	rollbar.Info(fmt.Sprintf("Updating projectID %s. \nOld Data: %v \nNew Data: %v", currentProjectData.ID, currentProjectData, newProjectData))
	updatedProject := entity.Project{
		ID:                    currentProjectData.ID,
		Created:               currentProjectData.Created,
		Name:                  newProjectData.Name,
		StartDate:             newProjectData.StartDate,
		DueDate:               newProjectData.DueDate,
		Address:               newProjectData.Address,
		City:                  newProjectData.City,
		State:                 newProjectData.State,
		Zip:                   newProjectData.Zip,
		ProjectManager:        newProjectData.ProjectManager,
		ClientID:              newProjectData.ClientID,
		EORNameID:             newProjectData.EORNameID,
		DetailerNameID:        newProjectData.DetailerNameID,
		InspectionLabID:       newProjectData.InspectionLabID,
		SteelErectorNameID:    newProjectData.SteelErectorNameID,
		SteelFabricatorNameID: newProjectData.SteelFabricatorNameID,
		GeneralContractorID:   newProjectData.GeneralContractorID,
		PrimaryContactNameID:  newProjectData.PrimaryContactNameID,
		PrimaryContactPhone:   newProjectData.PrimaryContactPhone,
		PrimaryContactEmail:   newProjectData.PrimaryContactEmail,
		SquareFootage:         newProjectData.SquareFootage,
		WeightInTons:          newProjectData.WeightInTons,
	}

	err = d.SharedProvider.Set(currentProjectData.ID, updatedProject)
	if err != nil {
		return entity.Project{}, err
	}

	var updatedProjectData = entity.Project{}
	err = d.SharedProvider.GetByID(currentProjectData.ID, &updatedProjectData)
	if err != nil {
		return entity.Project{}, err
	}

	rollbar.Info(fmt.Sprintf("Project %s updated.", updatedProjectData.Name))

	return updatedProjectData, nil
}

// Delete is to update a project record
func (d *DatabaseProvider) Delete(project entity.Project) error {
	rollbar.Info(fmt.Sprintf("Deleting Project from DB: %s", project.Name))
	var currentProject entity.Project

	err := d.SharedProvider.GetByID(project.ID, &currentProject)
	if (entity.Project{}) == currentProject {
		return fmt.Errorf("Project not found")
	}

	err = d.SharedProvider.Delete(project.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted project %s", project.ID))

	return nil
}
