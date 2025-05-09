package repositories

import (
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/segmentio/ksuid"
)

// DonoRepository define os métodos para acesso aos dados de Dono
type DonoRepository interface {
	// Métodos básicos de CRUD
	Create(dono *entities.Dono) error
	GetByID(id ksuid.KSUID) (*entities.Dono, error)
	Update(dono *entities.Dono) error
	Delete(id ksuid.KSUID) error
	List(page, limit int) ([]entities.Dono, error)

	// Métodos específicos para autenticação
	GetByEmail(email string) (*entities.Dono, error)
}
