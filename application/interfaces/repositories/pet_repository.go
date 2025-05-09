package repositories

import (
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/segmentio/ksuid"
)

// PetRepository define os métodos para acesso aos dados de Pet
type PetRepository interface {
	// Métodos básicos de CRUD
	Create(pet *entities.Pet) error
	GetByID(id ksuid.KSUID) (*entities.Pet, error)
	Update(pet *entities.Pet) error
	Delete(id ksuid.KSUID) error

	// Métodos específicos
	GetByDonoID(donoID ksuid.KSUID) ([]entities.Pet, error)
	List(page, limit int) ([]entities.Pet, error)
}
