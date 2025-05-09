package repositories

import (
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/segmentio/ksuid"
)

// AgendamentoRepository define os métodos para acesso aos dados de Agendamento
type AgendamentoRepository interface {
	// Métodos básicos de CRUD
	Create(agendamento *entities.Agendamento) error
	GetByID(id ksuid.KSUID) (*entities.Agendamento, error)
	Update(agendamento *entities.Agendamento) error
	UpdateStatus(id ksuid.KSUID, status entities.StatusAgendamento) error
	Delete(id ksuid.KSUID) error

	// Métodos específicos
	GetByDonoID(donoID ksuid.KSUID) ([]entities.Agendamento, error)
	GetByPetshopID(petshopID ksuid.KSUID) ([]entities.Agendamento, error)
	GetByPetID(petID ksuid.KSUID) ([]entities.Agendamento, error)
	GetAgendamentosFuturos(petshopID ksuid.KSUID) ([]entities.Agendamento, error)
}
