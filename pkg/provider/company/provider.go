package company

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with company data
type Provider interface {
	GetBy(companyname string) (entity.Company, error)
	GetAll() ([]entity.Company, error)
	AddCompany(entity.Company) (entity.Company, error)
	UpdateCompany(entity.Company) (entity.Company, error)
	DeleteCompany(companyname entity.Company) error
}
