package project

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with project data
type Provider interface {
	GetByName(projectname string) (entity.Project, error)
	GetAll() ([]entity.Project, error)
	AddProject(entity.Project) (entity.Project, error)
	UpdateProject(entity.Project) (entity.Project, error)
	DeleteProject(projectname entity.Project) error
}
