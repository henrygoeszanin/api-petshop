package repositories

import (
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// ProcedimentoRepositoryImpl implementa o repositório de Procedimento usando o GORM
type ProcedimentoRepositoryImpl struct {
	db *gorm.DB
}

// NewProcedimentoRepository cria uma nova instância do repositório de Procedimento
func NewProcedimentoRepository(db *gorm.DB) *ProcedimentoRepositoryImpl {
	return &ProcedimentoRepositoryImpl{db: db}
}

// Create insere um novo procedimento no banco de dados
func (r *ProcedimentoRepositoryImpl) Create(procedimento *entities.Procedimento) error {
	result := r.db.Create(procedimento)
	if result.Error != nil {
		return errors.ErrInvalidData
	}
	return nil
}

// GetByID busca um procedimento pelo ID
func (r *ProcedimentoRepositoryImpl) GetByID(id ksuid.KSUID) (*entities.Procedimento, error) {
	var procedimento entities.Procedimento
	result := r.db.Preload("Itens").First(&procedimento, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInvalidData
	}
	return &procedimento, nil
}

// GetByPetID busca todos os procedimentos de um determinado pet
func (r *ProcedimentoRepositoryImpl) GetByPetID(petID ksuid.KSUID) ([]entities.Procedimento, error) {
	var procedimentos []entities.Procedimento
	result := r.db.Preload("Itens").Where("pet_id = ?", petID).Order("data_realizacao DESC").Find(&procedimentos)
	if result.Error != nil {
		return nil, errors.ErrInvalidData
	}
	return procedimentos, nil
}

// GetByPetshopID busca todos os procedimentos realizados por um determinado petshop
func (r *ProcedimentoRepositoryImpl) GetByPetshopID(petshopID ksuid.KSUID) ([]entities.Procedimento, error) {
	var procedimentos []entities.Procedimento
	result := r.db.Preload("Itens").Where("petshop_id = ?", petshopID).Order("data_realizacao DESC").Find(&procedimentos)
	if result.Error != nil {
		return nil, errors.ErrInvalidData
	}
	return procedimentos, nil
}
