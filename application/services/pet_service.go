package services

import (
	"time"

	"github.com/henrygoeszanin/api_petshop/application/dtos"
	"github.com/henrygoeszanin/api_petshop/application/interfaces/repositories"
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
)

// PetService fornece métodos para gerenciar operações de Pets
type PetService struct {
	petRepository  repositories.PetRepository
	donoRepository repositories.DonoRepository
}

// NewPetService cria uma nova instância de PetService
func NewPetService(petRepo repositories.PetRepository, donoRepo repositories.DonoRepository) *PetService {
	return &PetService{
		petRepository:  petRepo,
		donoRepository: donoRepo,
	}
}

// Create cria um novo pet
func (s *PetService) Create(dto *dtos.PetCreateDTO) (*dtos.PetResponseDTO, error) { // Converter DonoID string para ksuid.KSUID
	donoID, err := ksuid.Parse(dto.DonoID)
	if err != nil {
		return nil, errors.ErrInvalidID
	}

	// Verificar se o dono existe
	dono, err := s.donoRepository.GetByID(donoID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrDonoNotFound
		}
		return nil, errors.ErrFailedToCheckDono
	}

	if dono == nil {
		return nil, errors.ErrDonoNotFound
	}

	// Criar entidade Pet
	pet := &entities.Pet{
		Nome:       dto.Nome,
		Especie:    dto.Especie,
		Raca:       dto.Raca,
		Nascimento: dto.Nascimento,
		DonoID:     donoID,
	}
	// Salvar no repositório
	if err := s.petRepository.Create(pet); err != nil {
		return nil, errors.ErrFailedToCreatePet
	}

	// Preparar DTO de resposta
	return s.entityToResponseDTO(pet), nil
}

// GetByID busca um pet pelo ID
func (s *PetService) GetByID(id ksuid.KSUID) (*dtos.PetResponseDTO, error) {
	pet, err := s.petRepository.GetByID(id)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return s.entityToResponseDTO(pet), nil
}

// GetByDonoID lista todos os pets de um determinado dono
func (s *PetService) GetByDonoID(donoID ksuid.KSUID) ([]dtos.PetResponseDTO, error) { // Verificar se o dono existe
	_, err := s.donoRepository.GetByID(donoID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrDonoNotFound
		}
		return nil, errors.ErrFailedToCheckDono
	}

	// Buscar pets do dono
	pets, err := s.petRepository.GetByDonoID(donoID)
	if err != nil {
		return nil, errors.ErrFailedToCheckPet
	}

	// Converter para DTOs
	var petDTOs []dtos.PetResponseDTO
	for _, pet := range pets {
		petDTOs = append(petDTOs, *s.entityToResponseDTO(&pet))
	}

	return petDTOs, nil
}

// Helper para converter entidade Pet para DTO de resposta
func (s *PetService) entityToResponseDTO(pet *entities.Pet) *dtos.PetResponseDTO {
	return &dtos.PetResponseDTO{
		ID:         pet.ID,
		Nome:       pet.Nome,
		Especie:    pet.Especie,
		Raca:       pet.Raca,
		Nascimento: pet.Nascimento,
		DonoID:     pet.DonoID,
		CreatedAt:  pet.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  pet.UpdatedAt.Format(time.RFC3339),
	}
}
