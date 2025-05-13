package services

import (
	"time"

	"github.com/henrygoeszanin/api_petshop/application/dtos"
	"github.com/henrygoeszanin/api_petshop/application/interfaces/repositories"
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
)

// DonoService fornece métodos para gerenciar operações de Donos
type DonoService struct {
	donoRepository repositories.DonoRepository
}

// NewDonoService cria uma nova instância de DonoService
func NewDonoService(donoRepo repositories.DonoRepository) *DonoService {
	return &DonoService{
		donoRepository: donoRepo,
	}
}

// Create cria um novo dono
func (s *DonoService) Create(dto *dtos.DonoCreateDTO) (*dtos.DonoDetailDTO, error) {
	// Verificar duplicidade por email
	existing, err := s.donoRepository.GetByEmail(dto.Email)
	if err != errors.ErrNotFound {
		return nil, errors.ErrCheckExistingOwner
	}
	if existing != nil {
		return nil, errors.ErrAlreadyExists
	}

	// Criar entidade Dono
	dono := &entities.Dono{
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
	}

	// Gerar hash da senha
	if err := dono.SetPassword(dto.Password); err != nil {
		return nil, errors.ErrSetPassword
	}

	// Salvar no repositório
	if err := s.donoRepository.Create(dono); err != nil {
		return nil, errors.ErrCreateOwner
	}

	// Preparar DTO de resposta
	return s.entityToDetailDTO(dono), nil
}

// GetByID busca um dono pelo ID
func (s *DonoService) GetByID(id ksuid.KSUID) (*dtos.DonoDetailDTO, error) {
	dono, err := s.donoRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.entityToDetailDTO(dono), nil
}

// Update atualiza os dados básicos de um dono
func (s *DonoService) Update(id ksuid.KSUID, dto *dtos.DonoUpdateDTO) (*dtos.DonoDetailDTO, error) {
	// Buscar dono existente
	dono, err := s.donoRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verificar se está tentando alterar para um email que já existe em outro registro
	if dono.Email != dto.Email {
		existing, err := s.donoRepository.GetByEmail(dto.Email)
		if err != nil {
			return nil, errors.ErrCheckExistingOwner
		}
		if existing != nil && existing.ID != dono.ID {
			return nil, errors.ErrAlreadyExists
		}
	}

	// Atualizar campos
	dono.Nome = dto.Nome
	dono.Email = dto.Email
	dono.Telefone = dto.Telefone

	// Salvar no repositório
	if err := s.donoRepository.Update(dono); err != nil {
		return nil, errors.ErrUpdateOwner
	}

	return s.entityToDetailDTO(dono), nil
}

// UpdateLocalizacao atualiza os dados de localização de um dono
func (s *DonoService) UpdateLocalizacao(id ksuid.KSUID, dto *dtos.DonoUpdateLocalizacaoDTO) (*dtos.DonoDetailDTO, error) {
	// Buscar dono existente
	dono, err := s.donoRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Atualizar campos de localização
	dono.CEP = dto.CEP
	dono.Rua = dto.Rua
	dono.Bairro = dto.Bairro
	dono.Cidade = dto.Cidade
	dono.Estado = dto.Estado
	dono.Numero = dto.Numero
	dono.Complemento = dto.Complemento

	// Salvar no repositório
	if err := s.donoRepository.Update(dono); err != nil {
		return nil, errors.ErrUpdateOwnerLocation
	}

	return s.entityToDetailDTO(dono), nil
}

// Helper para converter entidade Dono para DTO de detalhe
func (s *DonoService) entityToDetailDTO(dono *entities.Dono) *dtos.DonoDetailDTO {
	return &dtos.DonoDetailDTO{
		ID:          dono.ID,
		Nome:        dono.Nome,
		Email:       dono.Email,
		Telefone:    dono.Telefone,
		CEP:         dono.CEP,
		Rua:         dono.Rua,
		Bairro:      dono.Bairro,
		Cidade:      dono.Cidade,
		Estado:      dono.Estado,
		Numero:      dono.Numero,
		Complemento: dono.Complemento,
		CreatedAt:   dono.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   dono.UpdatedAt.Format(time.RFC3339),
	}
}
