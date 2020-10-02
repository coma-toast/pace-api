package project

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with project data
type Provider interface {
	GetByName(projectname string) (entity.Project, error)
	GetAll() ([]entity.Project, error)
	Add(entity.Project) (entity.Project, error)
	Update(entity.UpdateProjectRequest) (entity.Project, error)
	Delete(projectname entity.Project) error
}
