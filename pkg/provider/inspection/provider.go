package inspection

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with inspection data
type Provider interface {
	GetByName(inspectionname string) (entity.Inspection, error)
	GetAll() ([]entity.Inspection, error)
	Add(entity.Inspection) (entity.Inspection, error)
	Update(entity.UpdateInspectionRequest) (entity.Inspection, error)
	Delete(inspectionname entity.Inspection) error
}
