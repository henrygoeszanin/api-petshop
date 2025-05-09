package repositories

import (
	"strings"

	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// PetshopRepositoryImpl implementa o repositório de Petshop usando o GORM
type PetshopRepositoryImpl struct {
	db *gorm.DB
}

// NewPetshopRepository cria uma nova instância do repositório de Petshop
func NewPetshopRepository(db *gorm.DB) *PetshopRepositoryImpl {
	return &PetshopRepositoryImpl{db: db}
}

// Create insere um novo petshop no banco de dados
func (r *PetshopRepositoryImpl) Create(petshop *entities.Petshop) error {
	result := r.db.Create(petshop)
	if result.Error != nil {
		return errors.ErrInvalidData
	}
	return nil
}

// GetByID busca um petshop pelo ID
func (r *PetshopRepositoryImpl) GetByID(id ksuid.KSUID) (*entities.Petshop, error) {
	var petshop entities.Petshop
	result := r.db.First(&petshop, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInvalidData
	}
	return &petshop, nil
}

// Update atualiza os dados de um petshop
func (r *PetshopRepositoryImpl) Update(petshop *entities.Petshop) error {
	result := r.db.Save(petshop)
	if result.Error != nil {
		return errors.ErrInvalidData
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// Delete exclui um petshop do banco de dados (soft delete)
func (r *PetshopRepositoryImpl) Delete(id ksuid.KSUID) error {
	result := r.db.Delete(&entities.Petshop{}, "id = ?", id)
	if result.Error != nil {
		return errors.ErrInvalidData
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// List retorna uma lista paginada de petshops
func (r *PetshopRepositoryImpl) List(page, limit int) ([]entities.Petshop, error) {
	var petshops []entities.Petshop

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	result := r.db.Offset(offset).Limit(limit).Find(&petshops)
	if result.Error != nil {
		return nil, errors.ErrInvalidData
	}

	return petshops, nil
}

// GetByEmail busca um petshop pelo email
func (r *PetshopRepositoryImpl) GetByEmail(email string) (*entities.Petshop, error) {
	var petshop entities.Petshop
	result := r.db.Where("email = ?", email).First(&petshop)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Retorna nil quando não encontrado para tratamento no serviço
		}
		return nil, errors.ErrInvalidData
	}

	return &petshop, nil
}

// FindByCity busca petshops por cidade (case-insensitive)
func (r *PetshopRepositoryImpl) FindByCity(city string, page, limit int) ([]entities.Petshop, error) {
	var petshops []entities.Petshop

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Usando LOWER para busca case-insensitive
	result := r.db.Where("LOWER(cidade) LIKE ?", "%"+strings.ToLower(city)+"%").
		Offset(offset).
		Limit(limit).
		Find(&petshops)

	if result.Error != nil {
		return nil, errors.ErrInvalidData
	}

	return petshops, nil
}
