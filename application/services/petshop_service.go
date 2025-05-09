package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/henrygoeszanin/api_petshop/application/dtos"
	"github.com/henrygoeszanin/api_petshop/application/interfaces/repositories"
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
)

// PetshopService fornece métodos para gerenciar operações de Petshops
type PetshopService struct {
	petshopRepository repositories.PetshopRepository
}

// NewPetshopService cria uma nova instância de PetshopService
func NewPetshopService(petshopRepo repositories.PetshopRepository) *PetshopService {
	return &PetshopService{
		petshopRepository: petshopRepo,
	}
}

// Create cria um novo petshop
func (s *PetshopService) Create(dto *dtos.PetshopCreateDTO) (*dtos.PetshopDetailDTO, error) {
	// Verificar duplicidade por email
	existing, err := s.petshopRepository.GetByEmail(dto.Email)
	if err != nil {
		return nil, fmt.Errorf("falha ao verificar petshop existente: %w", err)
	}
	if existing != nil {
		return nil, errors.ErrAlreadyExists
	}

	// Criar entidade Petshop
	petshop := &entities.Petshop{
		Nome:        dto.Nome,
		Email:       dto.Email,
		Telefone:    dto.Telefone,
		CEP:         dto.CEP,
		Rua:         dto.Rua,
		Bairro:      dto.Bairro,
		Cidade:      dto.Cidade,
		Estado:      dto.Estado,
		Numero:      dto.Numero,
		Complemento: dto.Complemento,
		Descricao:   dto.Descricao,
		Ativo:       true, // Por padrão, o petshop é criado como ativo
		Nota:        0,    // Inicialmente sem avaliações
	}

	// Gerar hash da senha
	if err := petshop.SetPassword(dto.Password); err != nil {
		return nil, fmt.Errorf("falha ao definir senha: %w", err)
	}

	// Salvar no repositório
	if err := s.petshopRepository.Create(petshop); err != nil {
		return nil, fmt.Errorf("falha ao criar petshop: %w", err)
	}

	// Preparar DTO de resposta
	return s.entityToDetailDTO(petshop), nil
}

// GetByID busca um petshop pelo ID
func (s *PetshopService) GetByID(id ksuid.KSUID) (*dtos.PetshopDetailDTO, error) {
	petshop, err := s.petshopRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.entityToDetailDTO(petshop), nil
}

// Update atualiza os dados básicos de um petshop
func (s *PetshopService) Update(id ksuid.KSUID, dto *dtos.PetshopUpdateDTO) (*dtos.PetshopDetailDTO, error) {
	// Buscar petshop existente
	petshop, err := s.petshopRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verificar se está tentando alterar para um email que já existe em outro registro
	if petshop.Email != dto.Email {
		existing, err := s.petshopRepository.GetByEmail(dto.Email)
		if err != nil {
			return nil, fmt.Errorf("falha ao verificar petshop existente: %w", err)
		}
		if existing != nil && existing.ID != petshop.ID {
			return nil, errors.ErrAlreadyExists
		}
	}

	// Atualizar campos
	petshop.Nome = dto.Nome
	petshop.Email = dto.Email
	petshop.Telefone = dto.Telefone
	petshop.Descricao = dto.Descricao

	// Salvar no repositório
	if err := s.petshopRepository.Update(petshop); err != nil {
		return nil, fmt.Errorf("falha ao atualizar petshop: %w", err)
	}

	return s.entityToDetailDTO(petshop), nil
}

// UpdateEndereco atualiza os dados de endereço de um petshop
func (s *PetshopService) UpdateEndereco(id ksuid.KSUID, dto *dtos.PetshopUpdateEnderecoDTO) (*dtos.PetshopDetailDTO, error) {
	// Buscar petshop existente
	petshop, err := s.petshopRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Atualizar campos de endereço
	petshop.CEP = dto.CEP
	petshop.Rua = dto.Rua
	petshop.Bairro = dto.Bairro
	petshop.Cidade = dto.Cidade
	petshop.Estado = dto.Estado
	petshop.Numero = dto.Numero
	petshop.Complemento = dto.Complemento

	// Salvar no repositório
	if err := s.petshopRepository.Update(petshop); err != nil {
		return nil, fmt.Errorf("falha ao atualizar endereço do petshop: %w", err)
	}

	return s.entityToDetailDTO(petshop), nil
}

// FindByCity busca petshops em uma determinada cidade
func (s *PetshopService) FindByCity(city string, page, limit int) ([]dtos.PetshopListItemDTO, error) {
	// Padronizar cidade para busca case-insensitive
	normalizedCity := strings.ToLower(strings.TrimSpace(city))

	// Buscar no repositório
	petshops, err := s.petshopRepository.FindByCity(normalizedCity, page, limit)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar petshops por cidade: %w", err)
	}

	// Converter para DTOs de listagem
	var petshopDTOs []dtos.PetshopListItemDTO
	for _, petshop := range petshops {
		petshopDTOs = append(petshopDTOs, dtos.PetshopListItemDTO{
			ID:        petshop.ID,
			Nome:      petshop.Nome,
			Cidade:    petshop.Cidade,
			Estado:    petshop.Estado,
			Nota:      petshop.Nota,
			Descricao: petshop.Descricao,
		})
	}

	return petshopDTOs, nil
}

// List retorna uma lista paginada de petshops
func (s *PetshopService) List(page, limit int) ([]dtos.PetshopListItemDTO, error) {
	// Buscar no repositório
	petshops, err := s.petshopRepository.List(page, limit)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar petshops: %w", err)
	}

	// Converter para DTOs de listagem
	var petshopDTOs []dtos.PetshopListItemDTO
	for _, petshop := range petshops {
		petshopDTOs = append(petshopDTOs, dtos.PetshopListItemDTO{
			ID:        petshop.ID,
			Nome:      petshop.Nome,
			Cidade:    petshop.Cidade,
			Estado:    petshop.Estado,
			Nota:      petshop.Nota,
			Descricao: petshop.Descricao,
		})
	}

	return petshopDTOs, nil
}

// Helper para converter entidade Petshop para DTO de detalhe
func (s *PetshopService) entityToDetailDTO(petshop *entities.Petshop) *dtos.PetshopDetailDTO {
	return &dtos.PetshopDetailDTO{
		ID:          petshop.ID,
		Nome:        petshop.Nome,
		Email:       petshop.Email,
		Telefone:    petshop.Telefone,
		CEP:         petshop.CEP,
		Rua:         petshop.Rua,
		Bairro:      petshop.Bairro,
		Cidade:      petshop.Cidade,
		Estado:      petshop.Estado,
		Numero:      petshop.Numero,
		Complemento: petshop.Complemento,
		Descricao:   petshop.Descricao,
		Nota:        petshop.Nota,
		Ativo:       petshop.Ativo,
		CreatedAt:   petshop.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   petshop.UpdatedAt.Format(time.RFC3339),
	}
}
