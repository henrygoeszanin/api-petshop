package repositories

import (
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/segmentio/ksuid"
)

// ProcedimentoRepository define os métodos para acesso aos dados de Procedimento
type ProcedimentoRepository interface {
	// Métodos básicos de CRUD
	Create(procedimento *entities.Procedimento) error
	GetByID(id ksuid.KSUID) (*entities.Procedimento, error)

	// Métodos específicos
	GetByPetID(petID ksuid.KSUID) ([]entities.Procedimento, error)
	GetByPetshopID(petshopID ksuid.KSUID) ([]entities.Procedimento, error)
}
