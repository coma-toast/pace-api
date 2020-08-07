package inspection

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with inspection data
type Provider interface {
	GetByID(ID string) (entity.Inspection, error)
	GetAll() ([]entity.Inspection, error)
	Add(entity.UpdateInspectionRequest) (entity.Inspection, error)
	Update(entity.UpdateInspectionRequest) (entity.Inspection, error)
	Delete(inspection entity.Inspection) error
}
