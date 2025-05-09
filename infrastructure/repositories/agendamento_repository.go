package repositories

import (
	"fmt"
	"time"

	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// AgendamentoRepositoryImpl implementa o repositório de Agendamento usando o GORM
type AgendamentoRepositoryImpl struct {
	db *gorm.DB
}

// NewAgendamentoRepository cria uma nova instância do repositório de Agendamento
func NewAgendamentoRepository(db *gorm.DB) *AgendamentoRepositoryImpl {
	return &AgendamentoRepositoryImpl{db: db}
}

// Create insere um novo agendamento no banco de dados
func (r *AgendamentoRepositoryImpl) Create(agendamento *entities.Agendamento) error {
	result := r.db.Create(agendamento)
	if result.Error != nil {
		return fmt.Errorf("erro ao criar agendamento: %w", result.Error)
	}
	return nil
}

// GetByID busca um agendamento pelo ID
func (r *AgendamentoRepositoryImpl) GetByID(id ksuid.KSUID) (*entities.Agendamento, error) {
	var agendamento entities.Agendamento
	result := r.db.Preload("Itens").First(&agendamento, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("erro ao buscar agendamento: %w", result.Error)
	}
	return &agendamento, nil
}

// Update atualiza os dados de um agendamento
func (r *AgendamentoRepositoryImpl) Update(agendamento *entities.Agendamento) error {
	// Primeiro, verificamos se o agendamento existe
	var exists entities.Agendamento
	if err := r.db.First(&exists, "id = ?", agendamento.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return fmt.Errorf("erro ao verificar existência do agendamento: %w", err)
	}

	// Começar uma transação para garantir atomicidade
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Atualizar os itens do agendamento requer excluir os existentes e criar novos
	if err := tx.Where("agendamento_id = ?", agendamento.ID).Delete(&entities.ItemAgendamento{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao excluir itens existentes: %w", err)
	}

	// Salvar o agendamento atualizado
	if err := tx.Save(agendamento).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao atualizar agendamento: %w", err)
	}

	return tx.Commit().Error
}

// UpdateStatus atualiza apenas o status de um agendamento
func (r *AgendamentoRepositoryImpl) UpdateStatus(id ksuid.KSUID, status entities.StatusAgendamento) error {
	result := r.db.Model(&entities.Agendamento{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("erro ao atualizar status do agendamento: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// Delete exclui um agendamento do banco de dados (soft delete)
func (r *AgendamentoRepositoryImpl) Delete(id ksuid.KSUID) error {
	result := r.db.Delete(&entities.Agendamento{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("erro ao excluir agendamento: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// GetByDonoID busca todos os agendamentos de um determinado dono
func (r *AgendamentoRepositoryImpl) GetByDonoID(donoID ksuid.KSUID) ([]entities.Agendamento, error) {
	var agendamentos []entities.Agendamento
	result := r.db.Preload("Itens").Where("dono_id = ?", donoID).Order("data_agendada DESC").Find(&agendamentos)
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao buscar agendamentos do dono: %w", result.Error)
	}
	return agendamentos, nil
}

// GetByPetshopID busca todos os agendamentos de um determinado petshop
func (r *AgendamentoRepositoryImpl) GetByPetshopID(petshopID ksuid.KSUID) ([]entities.Agendamento, error) {
	var agendamentos []entities.Agendamento
	result := r.db.Preload("Itens").Where("petshop_id = ?", petshopID).Order("data_agendada DESC").Find(&agendamentos)
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao buscar agendamentos do petshop: %w", result.Error)
	}
	return agendamentos, nil
}

// GetByPetID busca todos os agendamentos de um determinado pet
func (r *AgendamentoRepositoryImpl) GetByPetID(petID ksuid.KSUID) ([]entities.Agendamento, error) {
	var agendamentos []entities.Agendamento
	result := r.db.Preload("Itens").Where("pet_id = ?", petID).Order("data_agendada DESC").Find(&agendamentos)
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao buscar agendamentos do pet: %w", result.Error)
	}
	return agendamentos, nil
}

// GetAgendamentosFuturos busca todos os agendamentos futuros de um determinado petshop
func (r *AgendamentoRepositoryImpl) GetAgendamentosFuturos(petshopID ksuid.KSUID) ([]entities.Agendamento, error) {
	var agendamentos []entities.Agendamento
	now := time.Now()
	result := r.db.Preload("Itens").
		Where("petshop_id = ? AND data_agendada > ? AND status != ?",
			petshopID, now, entities.StatusCancelado).
		Order("data_agendada ASC").
		Find(&agendamentos)
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao buscar agendamentos futuros do petshop: %w", result.Error)
	}
	return agendamentos, nil
}
