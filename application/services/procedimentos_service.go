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

// ProcedimentoService fornece métodos para gerenciar operações de Procedimentos
type ProcedimentoService struct {
	procedimentoRepository repositories.ProcedimentoRepository
	petRepository          repositories.PetRepository
	petshopRepository      repositories.PetshopRepository
	servicoRepository      repositories.ServicoRepository
}

// NewProcedimentoService cria uma nova instância de ProcedimentoService
func NewProcedimentoService(
	procedimentoRepo repositories.ProcedimentoRepository,
	petRepo repositories.PetRepository,
	petshopRepo repositories.PetshopRepository,
	servicoRepo repositories.ServicoRepository,
) *ProcedimentoService {
	return &ProcedimentoService{
		procedimentoRepository: procedimentoRepo,
		petRepository:          petRepo,
		petshopRepository:      petshopRepo,
		servicoRepository:      servicoRepo,
	}
}

// Create cria um novo registro de procedimento
func (s *ProcedimentoService) Create(dto *dtos.ProcedimentoCreateDTO) (*dtos.ProcedimentoResponseDTO, error) {
	// Converter IDs de string para ksuid.KSUID
	petID, err := ksuid.Parse(dto.PetID)
	if err != nil {
		return nil, fmt.Errorf("ID do pet inválido: %w", err)
	}

	petshopID, err := ksuid.Parse(dto.PetshopID)
	if err != nil {
		return nil, fmt.Errorf("ID do petshop inválido: %w", err)
	}

	// Verificar se o pet existe
	pet, err := s.petRepository.GetByID(petID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, fmt.Errorf("pet não encontrado")
		}
		return nil, fmt.Errorf("erro ao verificar pet: %w", err)
	}

	// Verificar se o petshop existe
	petshop, err := s.petshopRepository.GetByID(petshopID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, fmt.Errorf("petshop não encontrado")
		}
		return nil, fmt.Errorf("erro ao verificar petshop: %w", err)
	}

	// Converter data de string para time.Time
	dataRealizacao, err := time.Parse("2006-01-02T15:04:05Z07:00", dto.DataRealizacao)
	if err != nil {
		return nil, fmt.Errorf("formato de data inválido, use ISO 8601: %w", err)
	}

	// Validar que a data não é futura
	if dataRealizacao.After(time.Now()) {
		return nil, fmt.Errorf("a data de realização não pode ser futura")
	}

	// Criar entidade Procedimento
	procedimento := &entities.Procedimento{
		PetID:          petID,
		PetshopID:      petshopID,
		NomePetshop:    petshop.Nome,
		DataRealizacao: dataRealizacao,
		Observacoes:    dto.Observacoes,
		Total:          dto.Total,
		Itens:          []entities.ItemProcedimento{},
	}

	// Processar itens do procedimento
	var totalCalculado float64
	for _, itemDTO := range dto.Itens {
		servicoID, err := ksuid.Parse(itemDTO.ServicoID)
		if err != nil {
			return nil, fmt.Errorf("ID de serviço inválido: %w", err)
		}

		// Verificar se o serviço existe e pertence ao petshop
		servico, err := s.servicoRepository.GetByID(servicoID)
		if err != nil {
			if err == errors.ErrNotFound {
				return nil, fmt.Errorf("serviço não encontrado: %s", itemDTO.ServicoID)
			}
			return nil, fmt.Errorf("erro ao verificar serviço: %w", err)
		}

		if servico.PetshopID != petshopID {
			return nil, fmt.Errorf("o serviço %s não pertence ao petshop %s", itemDTO.ServicoID, dto.PetshopID)
		}

		// Adicionar item ao procedimento
		item := entities.ItemProcedimento{
			ServicoID:   servicoID,
			NomeServico: servico.Nome,
			PrecoFinal:  itemDTO.PrecoFinal,
		}
		procedimento.Itens = append(procedimento.Itens, item)
		totalCalculado += itemDTO.PrecoFinal
	}

	// Verificar se o total informado bate com a soma dos preços finais
	// Pequena margem de erro para lidar com arredondamentos
	const epsilon = 0.01
	if dto.Total < totalCalculado-epsilon || dto.Total > totalCalculado+epsilon {
		return nil, fmt.Errorf("o total informado (%.2f) não corresponde à soma dos preços finais (%.2f)", dto.Total, totalCalculado)
	}

	// Salvar no repositório
	if err := s.procedimentoRepository.Create(procedimento); err != nil {
		return nil, fmt.Errorf("falha ao criar procedimento: %w", err)
	}

	// Preparar DTO de resposta
	return s.entityToResponseDTO(procedimento, pet.Nome), nil
}

// GetByPetID lista todos os procedimentos de um determinado pet
func (s *ProcedimentoService) GetByPetID(petID ksuid.KSUID) ([]dtos.ProcedimentoResponseDTO, error) {
	// Verificar se o pet existe
	pet, err := s.petRepository.GetByID(petID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, fmt.Errorf("pet não encontrado")
		}
		return nil, fmt.Errorf("erro ao verificar pet: %w", err)
	}

	// Buscar procedimentos
	procedimentos, err := s.procedimentoRepository.GetByPetID(petID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar procedimentos do pet: %w", err)
	}

	// Converter para DTOs
	var procedimentoDTOs []dtos.ProcedimentoResponseDTO
	for _, procedimento := range procedimentos {
		procedimentoDTOs = append(procedimentoDTOs, *s.entityToResponseDTO(&procedimento, pet.Nome))
	}

	return procedimentoDTOs, nil
}

// Helper para converter entidade Procedimento para DTO
func (s *ProcedimentoService) entityToResponseDTO(procedimento *entities.Procedimento, nomePet string) *dtos.ProcedimentoResponseDTO {
	// Converter itens
	var itensDTO []dtos.ItemProcedimentoResponseDTO
	for _, item := range procedimento.Itens {
		itensDTO = append(itensDTO, dtos.ItemProcedimentoResponseDTO{
			ID:          item.ID.String(),
			ServicoID:   item.ServicoID.String(),
			NomeServico: item.NomeServico,
			PrecoFinal:  item.PrecoFinal,
		})
	}

	return &dtos.ProcedimentoResponseDTO{
		ID:             procedimento.ID.String(),
		PetID:          procedimento.PetID.String(),
		NomePet:        nomePet,
		PetshopID:      procedimento.PetshopID.String(),
		NomePetshop:    procedimento.NomePetshop,
		DataRealizacao: procedimento.DataRealizacao.Format(time.RFC3339),
		Observacoes:    procedimento.Observacoes,
		Total:          procedimento.Total,
		Itens:          itensDTO,
		CreatedAt:      procedimento.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      procedimento.UpdatedAt.Format(time.RFC3339),
	}
}
