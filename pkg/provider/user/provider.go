package user

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with User data
type Provider interface {
	GetByUsername(username string) (entity.User, error)
	GetAll() ([]entity.User, error)
}
