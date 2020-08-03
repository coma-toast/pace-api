package inventory

import "github.com/coma-toast/pace-api/pkg/entity"

// Provider is for working with inventory data
type Provider interface {
	GetByName(inventoryname string) (entity.Inventory, error)
	GetAll() ([]entity.Inventory, error)
	Add(entity.Inventory) (entity.Inventory, error)
	Update(entity.UpdateInventoryRequest) (entity.Inventory, error)
	Delete(inventoryname entity.Inventory) error
}
