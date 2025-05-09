package repositories

import (
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/segmentio/ksuid"
)

// PetshopRepository define os métodos para acesso aos dados de Petshop
type PetshopRepository interface {
	// Métodos básicos de CRUD
	Create(petshop *entities.Petshop) error
	GetByID(id ksuid.KSUID) (*entities.Petshop, error)
	Update(petshop *entities.Petshop) error
	Delete(id ksuid.KSUID) error
	List(page, limit int) ([]entities.Petshop, error)
	FindByCity(city string, page int, limit int) ([]entities.Petshop, error)

	// Métodos específicos para autenticação
	GetByEmail(email string) (*entities.Petshop, error)
}
