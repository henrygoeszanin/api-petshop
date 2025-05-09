package repositories

import (
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// ServicoRepositoryImpl implementa o repositório de Serviço usando o GORM
type ServicoRepositoryImpl struct {
	db *gorm.DB
}

// NewServicoRepository cria uma nova instância do repositório de Serviço
func NewServicoRepository(db *gorm.DB) *ServicoRepositoryImpl {
	return &ServicoRepositoryImpl{db: db}
}

// Create insere um novo serviço no banco de dados
func (r *ServicoRepositoryImpl) Create(servico *entities.Servico) error {
	result := r.db.Create(servico)
	if result.Error != nil {
		return errors.ErrInvalidData
	}
	return nil
}

// GetByID busca um serviço pelo ID
func (r *ServicoRepositoryImpl) GetByID(id ksuid.KSUID) (*entities.Servico, error) {
	var servico entities.Servico
	result := r.db.First(&servico, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInvalidData
	}
	return &servico, nil
}

// Update atualiza os dados de um serviço
func (r *ServicoRepositoryImpl) Update(servico *entities.Servico) error {
	result := r.db.Save(servico)
	if result.Error != nil {
		return errors.ErrInvalidData
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// Delete exclui um serviço do banco de dados (soft delete)
func (r *ServicoRepositoryImpl) Delete(id ksuid.KSUID) error {
	result := r.db.Delete(&entities.Servico{}, "id = ?", id)
	if result.Error != nil {
		return errors.ErrInvalidData
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// GetByPetshopID busca todos os serviços de um determinado petshop
func (r *ServicoRepositoryImpl) GetByPetshopID(petshopID ksuid.KSUID) ([]entities.Servico, error) {
	var servicos []entities.Servico
	result := r.db.Where("petshop_id = ?", petshopID).Find(&servicos)
	if result.Error != nil {
		return nil, errors.ErrInvalidData
	}
	return servicos, nil
}

// GetByName busca um serviço pelo nome em um determinado petshop
func (r *ServicoRepositoryImpl) GetByName(petshopID ksuid.KSUID, nome string) (*entities.Servico, error) {
	var servico entities.Servico
	result := r.db.Where("petshop_id = ? AND nome = ?", petshopID, nome).First(&servico)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Retorna nil quando não encontrado para tratamento no serviço
		}
		return nil, errors.ErrInvalidData
	}
	return &servico, nil
}
