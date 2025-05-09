package services

import (
	"fmt"
	"time"

	"github.com/henrygoeszanin/api_petshop/application/dtos"
	"github.com/henrygoeszanin/api_petshop/application/interfaces/repositories"
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
)

// ServicoService fornece métodos para gerenciar operações de Serviços
type ServicoService struct {
	servicoRepository repositories.ServicoRepository
	petshopRepository repositories.PetshopRepository
}

// NewServicoService cria uma nova instância de ServicoService
func NewServicoService(servicoRepo repositories.ServicoRepository, petshopRepo repositories.PetshopRepository) *ServicoService {
	return &ServicoService{
		servicoRepository: servicoRepo,
		petshopRepository: petshopRepo,
	}
}

// Create cria um novo serviço para um petshop
func (s *ServicoService) Create(petshopID ksuid.KSUID, dto *dtos.ServicoCreateDTO) (*dtos.ServicoResponseDTO, error) {
	// Verificar se o petshop existe
	_, err := s.petshopRepository.GetByID(petshopID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("falha ao verificar petshop: %w", err)
	}

	// Verificar se já existe um serviço com o mesmo nome neste petshop
	existingService, err := s.servicoRepository.GetByName(petshopID, dto.Nome)
	if err != nil {
		return nil, fmt.Errorf("falha ao verificar serviço existente: %w", err)
	}
	if existingService != nil {
		return nil, errors.ErrAlreadyExists
	}

	// Criar entidade Serviço
	servico := &entities.Servico{
		PetshopID: petshopID,
		Nome:      dto.Nome,
		Descricao: dto.Descricao,
		PrecoBase: dto.PrecoBase,
		Ativo:     true, // Por padrão, o serviço é criado como ativo
	}

	// Salvar no repositório
	if err := s.servicoRepository.Create(servico); err != nil {
		return nil, fmt.Errorf("falha ao criar serviço: %w", err)
	}

	// Preparar DTO de resposta
	return s.entityToDTO(servico), nil
}

// GetByID busca um serviço pelo ID
func (s *ServicoService) GetByID(id ksuid.KSUID) (*dtos.ServicoResponseDTO, error) {
	servico, err := s.servicoRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.entityToDTO(servico), nil
}

// Update atualiza um serviço existente
func (s *ServicoService) Update(id ksuid.KSUID, dto *dtos.ServicoUpdateDTO) (*dtos.ServicoResponseDTO, error) {
	// Buscar serviço existente
	servico, err := s.servicoRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verificar se está tentando alterar para um nome que já existe em outro serviço do mesmo petshop
	if servico.Nome != dto.Nome {
		existingService, err := s.servicoRepository.GetByName(servico.PetshopID, dto.Nome)
		if err != nil {
			return nil, fmt.Errorf("falha ao verificar serviço existente: %w", err)
		}
		if existingService != nil && existingService.ID != servico.ID {
			return nil, errors.ErrAlreadyExists
		}
	}

	// Atualizar campos
	servico.Nome = dto.Nome
	servico.Descricao = dto.Descricao
	servico.PrecoBase = dto.PrecoBase

	// Salvar no repositório
	if err := s.servicoRepository.Update(servico); err != nil {
		return nil, fmt.Errorf("falha ao atualizar serviço: %w", err)
	}

	return s.entityToDTO(servico), nil
}

// Delete desativa um serviço existente (soft delete)
func (s *ServicoService) Delete(id ksuid.KSUID) error {
	return s.servicoRepository.Delete(id)
}

// GetByPetshopID lista todos os serviços de um determinado petshop
func (s *ServicoService) GetByPetshopID(petshopID ksuid.KSUID) ([]dtos.ServicoResponseDTO, error) {
	// Verificar se o petshop existe
	_, err := s.petshopRepository.GetByID(petshopID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("falha ao verificar petshop: %w", err)
	}

	// Buscar serviços
	servicos, err := s.servicoRepository.GetByPetshopID(petshopID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar serviços do petshop: %w", err)
	}

	// Converter para DTOs
	var servicoDTOs []dtos.ServicoResponseDTO
	for _, servico := range servicos {
		servicoDTOs = append(servicoDTOs, *s.entityToDTO(&servico))
	}

	return servicoDTOs, nil
}

// Helper para converter entidade Serviço para DTO
func (s *ServicoService) entityToDTO(servico *entities.Servico) *dtos.ServicoResponseDTO {
	return &dtos.ServicoResponseDTO{
		ID:        servico.ID,
		PetshopID: servico.PetshopID,
		Nome:      servico.Nome,
		Descricao: servico.Descricao,
		PrecoBase: servico.PrecoBase,
		Ativo:     servico.Ativo,
		CreatedAt: servico.CreatedAt.Format(time.RFC3339),
		UpdatedAt: servico.UpdatedAt.Format(time.RFC3339),
	}
}
