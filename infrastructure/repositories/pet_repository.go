package repositories

import (
	"fmt"

	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// PetRepositoryImpl implementa o repositório de Pet usando o GORM
type PetRepositoryImpl struct {
	db *gorm.DB
}

// NewPetRepository cria uma nova instância do repositório de Pet
func NewPetRepository(db *gorm.DB) *PetRepositoryImpl {
	return &PetRepositoryImpl{db: db}
}

// Create insere um novo pet no banco de dados
func (r *PetRepositoryImpl) Create(pet *entities.Pet) error {
	result := r.db.Create(pet)
	if result.Error != nil {
		return fmt.Errorf("erro ao criar pet: %w", result.Error)
	}
	return nil
}

// GetByID busca um pet pelo ID
func (r *PetRepositoryImpl) GetByID(id ksuid.KSUID) (*entities.Pet, error) {
	var pet entities.Pet
	result := r.db.First(&pet, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("erro ao buscar pet: %w", result.Error)
	}
	return &pet, nil
}

// Update atualiza os dados de um pet
func (r *PetRepositoryImpl) Update(pet *entities.Pet) error {
	result := r.db.Save(pet)
	if result.Error != nil {
		return fmt.Errorf("erro ao atualizar pet: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// Delete exclui um pet do banco de dados (soft delete)
func (r *PetRepositoryImpl) Delete(id ksuid.KSUID) error {
	result := r.db.Delete(&entities.Pet{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("erro ao excluir pet: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// GetByDonoID lista todos os pets de um determinado dono
func (r *PetRepositoryImpl) GetByDonoID(donoID ksuid.KSUID) ([]entities.Pet, error) {
	var pets []entities.Pet
	result := r.db.Where("dono_id = ?", donoID).Find(&pets)
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao buscar pets do dono: %w", result.Error)
	}
	return pets, nil
}

// List retorna uma lista paginada de pets
func (r *PetRepositoryImpl) List(page, limit int) ([]entities.Pet, error) {
	var pets []entities.Pet
	offset := (page - 1) * limit

	result := r.db.Offset(offset).Limit(limit).Find(&pets)
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao listar pets: %w", result.Error)
	}
	return pets, nil
}
