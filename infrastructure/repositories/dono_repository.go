package repositories

import (
	"fmt"

	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// DonoRepositoryImpl implementa o repositório de Dono usando o GORM
type DonoRepositoryImpl struct {
	db *gorm.DB
}

// NewDonoRepository cria uma nova instância do repositório de Dono
func NewDonoRepository(db *gorm.DB) *DonoRepositoryImpl {
	return &DonoRepositoryImpl{db: db}
}

// Create insere um novo dono no banco de dados
func (r *DonoRepositoryImpl) Create(dono *entities.Dono) error {
	result := r.db.Create(dono)
	if result.Error != nil {
		return fmt.Errorf("erro ao criar dono: %w", result.Error)
	}
	return nil
}

// GetByID busca um dono pelo ID
func (r *DonoRepositoryImpl) GetByID(id ksuid.KSUID) (*entities.Dono, error) {
	var dono entities.Dono
	result := r.db.First(&dono, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("erro ao buscar dono: %w", result.Error)
	}
	return &dono, nil
}

// Update atualiza os dados de um dono
func (r *DonoRepositoryImpl) Update(dono *entities.Dono) error {
	result := r.db.Save(dono)
	if result.Error != nil {
		return fmt.Errorf("erro ao atualizar dono: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// Delete exclui um dono do banco de dados (soft delete)
func (r *DonoRepositoryImpl) Delete(id ksuid.KSUID) error {
	result := r.db.Delete(&entities.Dono{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("erro ao excluir dono: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// List retorna uma lista paginada de donos
func (r *DonoRepositoryImpl) List(page, limit int) ([]entities.Dono, error) {
	var donos []entities.Dono

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	result := r.db.Offset(offset).Limit(limit).Find(&donos)
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao listar donos: %w", result.Error)
	}

	return donos, nil
}

// GetByEmail busca um dono pelo email
func (r *DonoRepositoryImpl) GetByEmail(email string) (*entities.Dono, error) {
	var dono entities.Dono
	result := r.db.Where("email = ?", email).First(&dono)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Retorna nil quando não encontrado para tratamento no serviço
		}
		return nil, fmt.Errorf("erro ao buscar dono por email: %w", result.Error)
	}

	return &dono, nil
}
