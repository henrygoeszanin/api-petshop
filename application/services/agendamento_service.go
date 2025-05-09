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

// AgendamentoService fornece métodos para gerenciar operações de Agendamentos
type AgendamentoService struct {
	agendamentoRepository repositories.AgendamentoRepository
	donoRepository        repositories.DonoRepository
	petRepository         repositories.PetRepository
	petshopRepository     repositories.PetshopRepository
	servicoRepository     repositories.ServicoRepository
}

// NewAgendamentoService cria uma nova instância de AgendamentoService
func NewAgendamentoService(
	agendamentoRepo repositories.AgendamentoRepository,
	donoRepo repositories.DonoRepository,
	petRepo repositories.PetRepository,
	petshopRepo repositories.PetshopRepository,
	servicoRepo repositories.ServicoRepository,
) *AgendamentoService {
	return &AgendamentoService{
		agendamentoRepository: agendamentoRepo,
		donoRepository:        donoRepo,
		petRepository:         petRepo,
		petshopRepository:     petshopRepo,
		servicoRepository:     servicoRepo,
	}
}

// Create cria um novo agendamento
func (s *AgendamentoService) Create(dto *dtos.AgendamentoCreateDTO) (*dtos.AgendamentoResponseDTO, error) {
	// Converter IDs de string para KSUID
	donoID, err := ksuid.Parse(dto.DonoID)
	if err != nil {
		return nil, fmt.Errorf("ID do dono inválido: %w", err)
	}

	petID, err := ksuid.Parse(dto.PetID)
	if err != nil {
		return nil, fmt.Errorf("ID do pet inválido: %w", err)
	}

	petshopID, err := ksuid.Parse(dto.PetshopID)
	if err != nil {
		return nil, fmt.Errorf("ID do petshop inválido: %w", err)
	}

	// Verificar se o dono existe
	dono, err := s.donoRepository.GetByID(donoID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, fmt.Errorf("dono não encontrado")
		}
		return nil, fmt.Errorf("erro ao verificar dono: %w", err)
	}

	// Verificar se o pet existe e se pertence ao dono
	pet, err := s.petRepository.GetByID(petID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, fmt.Errorf("pet não encontrado")
		}
		return nil, fmt.Errorf("erro ao verificar pet: %w", err)
	}

	if pet.DonoID != donoID {
		return nil, fmt.Errorf("o pet não pertence ao dono especificado")
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
	dataAgendada, err := time.Parse("2006-01-02T15:04:05Z07:00", dto.DataAgendada)
	if err != nil {
		return nil, fmt.Errorf("formato de data inválido, use ISO 8601: %w", err)
	}

	// Validar que a data não é passada
	if dataAgendada.Before(time.Now()) {
		return nil, fmt.Errorf("a data agendada não pode ser no passado")
	}

	// Criar entidade Agendamento
	agendamento := &entities.Agendamento{
		DonoID:        donoID,
		PetID:         petID,
		PetshopID:     petshopID,
		DataAgendada:  dataAgendada,
		Status:        entities.StatusPendente,
		Observacoes:   dto.Observacoes,
		TotalPrevisto: dto.TotalPrevisto,
		Itens:         []entities.ItemAgendamento{},
	}

	// Processar itens do agendamento
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
				return nil, fmt.Errorf("serviço não encontrado")
			}
			return nil, fmt.Errorf("erro ao verificar serviço: %w", err)
		}

		if servico.PetshopID != petshopID {
			return nil, fmt.Errorf("o serviço não pertence ao petshop especificado")
		}

		if !servico.Ativo {
			return nil, fmt.Errorf("o serviço '%s' não está ativo", servico.Nome)
		}

		// Adicionar item ao agendamento
		agendamento.Itens = append(agendamento.Itens, entities.ItemAgendamento{
			ServicoID:     servicoID,
			NomeServico:   servico.Nome,
			PrecoPrevisto: itemDTO.PrecoPrevisto,
		})

		totalCalculado += itemDTO.PrecoPrevisto
	}

	// Validar o total previsto
	if totalCalculado != dto.TotalPrevisto {
		return nil, fmt.Errorf("o total previsto (%f) não corresponde à soma dos preços dos itens (%f)", dto.TotalPrevisto, totalCalculado)
	}

	// Salvar no repositório
	if err := s.agendamentoRepository.Create(agendamento); err != nil {
		return nil, fmt.Errorf("falha ao criar agendamento: %w", err)
	}

	// Converter para DTO de resposta
	return s.entityToResponseDTO(agendamento, pet.Nome, dono.Nome, petshop.Nome), nil
}

// GetByID busca um agendamento pelo ID
func (s *AgendamentoService) GetByID(id ksuid.KSUID) (*dtos.AgendamentoResponseDTO, error) {
	// Buscar agendamento
	agendamento, err := s.agendamentoRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Buscar informações adicionais para o DTO
	pet, err := s.petRepository.GetByID(agendamento.PetID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar informações do pet: %w", err)
	}

	dono, err := s.donoRepository.GetByID(agendamento.DonoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar informações do dono: %w", err)
	}

	petshop, err := s.petshopRepository.GetByID(agendamento.PetshopID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar informações do petshop: %w", err)
	}

	// Converter para DTO de resposta
	return s.entityToResponseDTO(agendamento, pet.Nome, dono.Nome, petshop.Nome), nil
}

// GetByDonoID lista todos os agendamentos de um dono
func (s *AgendamentoService) GetByDonoID(donoID ksuid.KSUID) ([]dtos.AgendamentoResponseDTO, error) {
	// Verificar se o dono existe
	_, err := s.donoRepository.GetByID(donoID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, fmt.Errorf("dono não encontrado")
		}
		return nil, fmt.Errorf("erro ao verificar dono: %w", err)
	}

	// Buscar agendamentos do dono
	agendamentos, err := s.agendamentoRepository.GetByDonoID(donoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agendamentos: %w", err)
	}

	// Converter para DTO de resposta
	var agendamentosDTO []dtos.AgendamentoResponseDTO
	for _, agendamento := range agendamentos {
		// Buscar informações adicionais
		pet, err := s.petRepository.GetByID(agendamento.PetID)
		if err != nil {
			continue // Pular este agendamento se não for possível buscar o pet
		}

		dono, err := s.donoRepository.GetByID(agendamento.DonoID)
		if err != nil {
			continue // Pular este agendamento se não for possível buscar o dono
		}

		petshop, err := s.petshopRepository.GetByID(agendamento.PetshopID)
		if err != nil {
			continue // Pular este agendamento se não for possível buscar o petshop
		}

		// Adicionar agendamento convertido
		agendamentosDTO = append(agendamentosDTO, *s.entityToResponseDTO(&agendamento, pet.Nome, dono.Nome, petshop.Nome))
	}

	return agendamentosDTO, nil
}

// GetByPetshopID lista todos os agendamentos de um petshop
func (s *AgendamentoService) GetByPetshopID(petshopID ksuid.KSUID) ([]dtos.AgendamentoResponseDTO, error) {
	// Verificar se o petshop existe
	_, err := s.petshopRepository.GetByID(petshopID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, fmt.Errorf("petshop não encontrado")
		}
		return nil, fmt.Errorf("erro ao verificar petshop: %w", err)
	}

	// Buscar agendamentos do petshop
	agendamentos, err := s.agendamentoRepository.GetByPetshopID(petshopID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agendamentos: %w", err)
	}

	// Converter para DTO de resposta
	var agendamentosDTO []dtos.AgendamentoResponseDTO
	for _, agendamento := range agendamentos {
		// Buscar informações adicionais
		pet, err := s.petRepository.GetByID(agendamento.PetID)
		if err != nil {
			continue // Pular este agendamento se não for possível buscar o pet
		}

		dono, err := s.donoRepository.GetByID(agendamento.DonoID)
		if err != nil {
			continue // Pular este agendamento se não for possível buscar o dono
		}

		petshop, err := s.petshopRepository.GetByID(agendamento.PetshopID)
		if err != nil {
			continue // Pular este agendamento se não for possível buscar o petshop
		}

		// Adicionar agendamento convertido
		agendamentosDTO = append(agendamentosDTO, *s.entityToResponseDTO(&agendamento, pet.Nome, dono.Nome, petshop.Nome))
	}

	return agendamentosDTO, nil
}

// UpdateStatus atualiza o status de um agendamento
func (s *AgendamentoService) UpdateStatus(id ksuid.KSUID, dto *dtos.AgendamentoUpdateStatusDTO) (*dtos.AgendamentoResponseDTO, error) {
	// Verificar se o agendamento existe
	agendamento, err := s.agendamentoRepository.GetByID(id)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("erro ao verificar agendamento: %w", err)
	}

	// Validar a transição de status
	statusAtual := agendamento.Status
	novoStatus := entities.StatusAgendamento(dto.Status)

	// Validações específicas de acordo com regras de negócio
	if statusAtual == entities.StatusCancelado && novoStatus != entities.StatusCancelado {
		return nil, fmt.Errorf("não é possível alterar o status de um agendamento cancelado")
	}

	if statusAtual == entities.StatusConcluido && novoStatus != entities.StatusConcluido {
		return nil, fmt.Errorf("não é possível alterar o status de um agendamento concluído")
	}

	// Atualizar status no banco de dados
	if err := s.agendamentoRepository.UpdateStatus(id, novoStatus); err != nil {
		return nil, fmt.Errorf("erro ao atualizar status: %w", err)
	}

	// Buscar agendamento atualizado
	agendamentoAtualizado, err := s.agendamentoRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agendamento atualizado: %w", err)
	}

	// Buscar informações adicionais para o DTO
	pet, err := s.petRepository.GetByID(agendamentoAtualizado.PetID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar informações do pet: %w", err)
	}

	dono, err := s.donoRepository.GetByID(agendamentoAtualizado.DonoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar informações do dono: %w", err)
	}

	petshop, err := s.petshopRepository.GetByID(agendamentoAtualizado.PetshopID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar informações do petshop: %w", err)
	}

	// Converter para DTO de resposta
	return s.entityToResponseDTO(agendamentoAtualizado, pet.Nome, dono.Nome, petshop.Nome), nil
}

// Update atualiza os dados de um agendamento
func (s *AgendamentoService) Update(id ksuid.KSUID, dto *dtos.AgendamentoUpdateDTO) (*dtos.AgendamentoResponseDTO, error) {
	// Verificar se o agendamento existe
	agendamento, err := s.agendamentoRepository.GetByID(id)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("erro ao verificar agendamento: %w", err)
	}

	// Não permitir atualização de agendamentos cancelados ou concluídos
	if agendamento.Status == entities.StatusCancelado {
		return nil, fmt.Errorf("não é possível atualizar um agendamento cancelado")
	}
	if agendamento.Status == entities.StatusConcluido {
		return nil, fmt.Errorf("não é possível atualizar um agendamento concluído")
	}

	// Converter data de string para time.Time
	dataAgendada, err := time.Parse("2006-01-02T15:04:05Z07:00", dto.DataAgendada)
	if err != nil {
		return nil, fmt.Errorf("formato de data inválido, use ISO 8601: %w", err)
	}

	// Validar que a data não é passada
	if dataAgendada.Before(time.Now()) {
		return nil, fmt.Errorf("a data agendada não pode ser no passado")
	}

	// Atualizar campos
	agendamento.DataAgendada = dataAgendada
	agendamento.Observacoes = dto.Observacoes
	agendamento.TotalPrevisto = dto.TotalPrevisto

	// Processar itens do agendamento
	itens := []entities.ItemAgendamento{}
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
				return nil, fmt.Errorf("serviço não encontrado")
			}
			return nil, fmt.Errorf("erro ao verificar serviço: %w", err)
		}

		if servico.PetshopID != agendamento.PetshopID {
			return nil, fmt.Errorf("o serviço não pertence ao petshop especificado")
		}

		if !servico.Ativo {
			return nil, fmt.Errorf("o serviço '%s' não está ativo", servico.Nome)
		}

		// Adicionar item ao agendamento
		itens = append(itens, entities.ItemAgendamento{
			AgendamentoID: agendamento.ID,
			ServicoID:     servicoID,
			NomeServico:   servico.Nome,
			PrecoPrevisto: itemDTO.PrecoPrevisto,
		})

		totalCalculado += itemDTO.PrecoPrevisto
	}

	// Validar o total previsto
	if totalCalculado != dto.TotalPrevisto {
		return nil, fmt.Errorf("o total previsto (%f) não corresponde à soma dos preços dos itens (%f)", dto.TotalPrevisto, totalCalculado)
	}

	// Atualizar itens
	agendamento.Itens = itens

	// Salvar no repositório
	if err := s.agendamentoRepository.Update(agendamento); err != nil {
		return nil, fmt.Errorf("falha ao atualizar agendamento: %w", err)
	}

	// Buscar informações adicionais para o DTO
	pet, err := s.petRepository.GetByID(agendamento.PetID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar informações do pet: %w", err)
	}

	dono, err := s.donoRepository.GetByID(agendamento.DonoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar informações do dono: %w", err)
	}

	petshop, err := s.petshopRepository.GetByID(agendamento.PetshopID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar informações do petshop: %w", err)
	}

	// Converter para DTO de resposta
	return s.entityToResponseDTO(agendamento, pet.Nome, dono.Nome, petshop.Nome), nil
}

// Helper para converter entidade Agendamento para DTO de resposta
func (s *AgendamentoService) entityToResponseDTO(agendamento *entities.Agendamento, nomePet string, nomeDono string, nomePetshop string) *dtos.AgendamentoResponseDTO {
	// Converter itens
	var itensDTO []dtos.ItemAgendamentoResponseDTO
	for _, item := range agendamento.Itens {
		itensDTO = append(itensDTO, dtos.ItemAgendamentoResponseDTO{
			ID:            item.ID.String(),
			ServicoID:     item.ServicoID.String(),
			NomeServico:   item.NomeServico,
			PrecoPrevisto: item.PrecoPrevisto,
		})
	}

	return &dtos.AgendamentoResponseDTO{
		ID:            agendamento.ID.String(),
		DonoID:        agendamento.DonoID.String(),
		NomeDono:      nomeDono,
		PetID:         agendamento.PetID.String(),
		NomePet:       nomePet,
		PetshopID:     agendamento.PetshopID.String(),
		NomePetshop:   nomePetshop,
		DataAgendada:  agendamento.DataAgendada.Format(time.RFC3339),
		Status:        string(agendamento.Status),
		Observacoes:   agendamento.Observacoes,
		TotalPrevisto: agendamento.TotalPrevisto,
		Itens:         itensDTO,
		CreatedAt:     agendamento.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     agendamento.UpdatedAt.Format(time.RFC3339),
	}
}
