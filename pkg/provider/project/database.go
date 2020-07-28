package project

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

// DatabaseProvider is a project.Provider the uses a database
type DatabaseProvider struct {
	Database *firestore.Client
}

// ErrProjectNotFound if no Projects are found
var ErrProjectNotFound = errors.New("Project not found")

// // GetByID gets a Project by ID
// func (d *DatabaseProvider) GetByID(ID string) (entity.Project, error) {
// 	return d.getByID(ID)
// }

// GetAll gets a Project by ID
func (d *DatabaseProvider) GetAll() ([]entity.Project, error) {
	return d.getAll()
}

func (d *DatabaseProvider) getAll() ([]entity.Project, error) {
	var projects []entity.Project

	allProjectData, err := d.Database.Collection("projects").Documents(context.TODO()).GetAll()
	if err != nil {
		return []entity.Project{}, err
	}

	for _, ProjectData := range allProjectData {
		var Project entity.Project
		err := ProjectData.DataTo(&Project)
		if err != nil {
			return []entity.Project{}, fmt.Errorf("ERROR: GetAll(): Firestore.DataTo() error %w", err)
		}
		projects = append(projects, Project)
	}

	return projects, nil
}

// GetByName gets a project by the name
func (d *DatabaseProvider) GetByName(name string) (entity.Project, error) {
	return d.getByName(name)
}

func (d *DatabaseProvider) getByName(name string) (entity.Project, error) {
	var project entity.Project
	projectData, err := d.Database.Collection("projects").Where("Name", "==", name).Documents(context.TODO()).GetAll()
	if err != nil {
		return entity.Project{}, fmt.Errorf("Error getting project by name and project: %s", err)
	}

	for _, projectItem := range projectData {
		projectItem.DataTo(&project)
		return project, nil
	}

	return entity.Project{}, ErrProjectNotFound
}

// AddProject is to update a Project record
func (d *DatabaseProvider) AddProject(newProjectData entity.Project) (entity.Project, error) {
	projectRef, err := d.addProject(newProjectData)
	if err != nil {
		return entity.Project{}, err
	}
	rollbar.Info(fmt.Sprintf("Adding new Project %s %s", newProjectData.Name, newProjectData.ID))
	updatedProjectData, err := d.getByProjectID(projectRef.ID)
	if err != nil {
		return entity.Project{}, err
	}

	return updatedProjectData, nil
}

// UpdateProject is to update a Project record
func (d *DatabaseProvider) UpdateProject(newProjectData entity.Project) (entity.Project, error) {
	currentProjectData, err := d.getByProjectID(newProjectData.ID)
	// * dev code currentProjectData, err := d.getByNameAndProject(newProjectData.FirstName, newProjectData.LastName, newProjectData.Project)
	if err != nil {
		return entity.Project{}, err
	}
	rollbar.Info(fmt.Sprintf("Updating ProjectID %s. \nOld Data: %v \nNew Data: %v", currentProjectData.ID, currentProjectData, newProjectData))
	updatedProject := entity.Project{
		ID:                  currentProjectData.ID,
		Created:             currentProjectData.Created,
		Name:                newProjectData.Name,
		StartDate:           newProjectData.StartDate,
		DueDate:             newProjectData.DueDate,
		Address:             newProjectData.Address,
		City:                newProjectData.City,
		State:               newProjectData.State,
		Zip:                 newProjectData.Zip,
		ProjectManager:      newProjectData.ProjectManager,
		ClientName:          newProjectData.ClientName,
		EORName:             newProjectData.EORName,
		DetailerName:        newProjectData.DetailerName,
		InspectionLab:       newProjectData.InspectionLab,
		SteelErectorName:    newProjectData.SteelErectorName,
		SteelFabricatorName: newProjectData.SteelFabricatorName,
		GeneralContractor:   newProjectData.GeneralContractor,
		PrimaryContactName:  newProjectData.PrimaryContactName,
		PrimaryContactPhone: newProjectData.PrimaryContactPhone,
		PrimaryContactEmail: newProjectData.PrimaryContactEmail,
		SquareFootage:       newProjectData.SquareFootage,
		WeightInTons:        newProjectData.WeightInTons,
	}

	err = d.setByProjectID(currentProjectData.ID, updatedProject)
	if err != nil {
		return entity.Project{}, err
	}
	updatedProjectData, err := d.getByProjectID(updatedProject.ID)
	if err != nil {
		return entity.Project{}, err
	}

	return updatedProjectData, nil
}

// DeleteProject is to update a Project record
func (d *DatabaseProvider) DeleteProject(project entity.Project) error {
	projectData, err := d.getByProjectID(project.ID)
	if err != nil {
		return err
	}

	err = d.deleteByProjectID(projectData.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted Project %s: %s %s", projectData.ID, projectData.Name, projectData.ID))

	return nil
}

func (d *DatabaseProvider) addProject(projectData entity.Project) (entity.Project, error) {
	existingProject, _ := d.getByNameAndClient(projectData.Name, projectData.ClientName)
	if (entity.Project{}) != existingProject {
		return entity.Project{}, fmt.Errorf("Error adding Project %s: ID already exists", projectData.ID)
	}
	newUUID := uuid.New().String()
	newProjectData := entity.Project{
		ID:                  newUUID,
		Created:             time.Now().String(),
		Name:                projectData.Name,
		StartDate:           projectData.StartDate,
		DueDate:             projectData.DueDate,
		Address:             projectData.Address,
		City:                projectData.City,
		State:               projectData.State,
		Zip:                 projectData.Zip,
		ProjectManager:      projectData.ProjectManager,
		ClientName:          projectData.ClientName,
		EORName:             projectData.EORName,
		DetailerName:        projectData.DetailerName,
		InspectionLab:       projectData.InspectionLab,
		SteelErectorName:    projectData.SteelErectorName,
		SteelFabricatorName: projectData.SteelFabricatorName,
		GeneralContractor:   projectData.GeneralContractor,
		PrimaryContactName:  projectData.PrimaryContactName,
		PrimaryContactPhone: projectData.PrimaryContactPhone,
		PrimaryContactEmail: projectData.PrimaryContactEmail,
		SquareFootage:       projectData.SquareFootage,
		WeightInTons:        projectData.WeightInTons,
	}
	addProjectResult, err := d.Database.Collection("projects").Doc(newUUID).Set(context.TODO(), newProjectData)
	if err != nil {
		return entity.Project{}, fmt.Errorf("Error setting Project %s by ID: %s", newProjectData.ID, err)
	}
	rollbar.Info(fmt.Sprintf("Project %s added at %s.", newProjectData.ID, addProjectResult))

	newProject, err := d.getByProjectID(newProjectData.ID)
	if err != nil {
		return entity.Project{}, fmt.Errorf("Error getting newly created Project %s by ID: %s", newProjectData.ID, err)
	}

	return newProject, nil
}

func (d *DatabaseProvider) getByProjectID(projectID string) (entity.Project, error) {
	var project entity.Project

	projectData, err := d.Database.Collection("projects").Doc(projectID).Get(context.TODO())
	if err != nil {
		return entity.Project{}, fmt.Errorf("Error getting Project %s by ID: %s", projectID, err)
	}
	projectData.DataTo(&project)

	return project, nil
}

func (d *DatabaseProvider) getByNameAndClient(name string, client string) (entity.Project, error) {
	var project entity.Project
	projectSnapshot, err := d.Database.Collection("projects").Where("ClientName", "==", client).Documents(context.TODO()).GetAll()
	if err != nil {
		return entity.Project{}, fmt.Errorf("Error getting project by name and client: %s", err)
	}

	for _, clientProject := range projectSnapshot {
		clientProject.DataTo(&project)
		if project.Name == name {
			return project, nil
		}
	}

	return entity.Project{}, ErrProjectNotFound
}

func (d *DatabaseProvider) setByProjectID(ProjectID string, ProjectData entity.Project) error {
	_, err := d.Database.Collection("projects").Doc(ProjectID).Set(context.TODO(), ProjectData)
	if err != nil {
		return fmt.Errorf("Error setting Project %s by ID: %s", ProjectID, err)
	}

	return nil
}

func (d *DatabaseProvider) deleteByProjectID(ProjectID string) error {
	result, err := d.Database.Collection("projects").Doc(ProjectID).Delete(context.TODO())
	if err != nil {
		return fmt.Errorf("Error deleting Project %s by ID: %s", ProjectID, err)
	}
	log.Printf("Deleting Project %s: %v", ProjectID, result)

	return nil
}

// func (d *DatabaseProvider) getByID(ID string) (entity.Project, error) {
// 	var Project entity.Project

// 	Projects := d.Database.Collection("projects").Where("ID", "==", ID).Documents(context.TODO())
// 	allMatchingProjects, err := Projects.GetAll()
// 	if err != nil {
// 		return entity.Project{}, err
// 	}
// 	for _, fbProject := range allMatchingProjects {
// 		err = fbProject.DataTo(&Project)
// 		if err != nil {
// 			return entity.Project{}, fmt.Errorf("ERROR: Project error - Firestore.DataTo() error %w, for Project %s", err, ID)
// 		}
// 		return Project, nil
// 		// data = append(data, fbProject.Data())
// 	}

// 	return entity.Project{}, ErrProjectNotFound
// }
