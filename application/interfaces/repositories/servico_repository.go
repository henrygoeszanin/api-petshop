package repositories

import (
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/segmentio/ksuid"
)

// ServicoRepository define os métodos para acesso aos dados de Serviço
type ServicoRepository interface {
	// Métodos básicos de CRUD
	Create(servico *entities.Servico) error
	GetByID(id ksuid.KSUID) (*entities.Servico, error)
	Update(servico *entities.Servico) error
	Delete(id ksuid.KSUID) error

	// Métodos específicos
	GetByPetshopID(petshopID ksuid.KSUID) ([]entities.Servico, error)
	GetByName(petshopID ksuid.KSUID, nome string) (*entities.Servico, error)
}
