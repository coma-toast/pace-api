package company

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with company data
type Provider interface {
	GetByName(companyname string) (entity.Company, error)
	GetAll() ([]entity.Company, error)
	Add(entity.Company) (entity.Company, error)
	Update(entity.Company) (entity.Company, error)
	Delete(companyname entity.Company) error
}
