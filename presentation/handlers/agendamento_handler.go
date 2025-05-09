package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/application/dtos"
	"github.com/henrygoeszanin/api_petshop/application/services"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
)

// AgendamentoHandler gerencia as requisições relacionadas a agendamentos
type AgendamentoHandler struct {
	agendamentoService *services.AgendamentoService
}

// NewAgendamentoHandler cria uma nova instância de AgendamentoHandler
func NewAgendamentoHandler(agendamentoService *services.AgendamentoService) *AgendamentoHandler {
	return &AgendamentoHandler{
		agendamentoService: agendamentoService,
	}
}

// Create processa a criação de um novo agendamento
func (h *AgendamentoHandler) Create(c *gin.Context) {
	var dto dtos.AgendamentoCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.agendamentoService.Create(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Erro ao criar agendamento: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetByID processa a requisição para buscar um agendamento por ID
func (h *AgendamentoHandler) GetByID(c *gin.Context) {
	// Extrair o ID da requisição
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Buscar agendamento no serviço
	agendamento, err := h.agendamentoService.GetByID(id)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar agendamento: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, agendamento)
}

// GetByDonoID processa a requisição para listar todos os agendamentos de um dono
func (h *AgendamentoHandler) GetByDonoID(c *gin.Context) {
	// Extrair o ID do dono da requisição
	donoIDStr := c.Param("donoId")
	donoID, err := ksuid.Parse(donoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do dono inválido"})
		return
	}

	// Buscar agendamentos do dono
	agendamentos, err := h.agendamentoService.GetByDonoID(donoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar agendamentos: %v", err)})
		return
	}

	c.JSON(http.StatusOK, agendamentos)
}

// GetByPetshopID processa a requisição para listar todos os agendamentos de um petshop
func (h *AgendamentoHandler) GetByPetshopID(c *gin.Context) {
	// Extrair o ID do petshop da requisição
	petshopIDStr := c.Param("petshopId")
	petshopID, err := ksuid.Parse(petshopIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do petshop inválido"})
		return
	}

	// Buscar agendamentos do petshop
	agendamentos, err := h.agendamentoService.GetByPetshopID(petshopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar agendamentos: %v", err)})
		return
	}

	c.JSON(http.StatusOK, agendamentos)
}

// UpdateStatus processa a atualização do status de um agendamento
func (h *AgendamentoHandler) UpdateStatus(c *gin.Context) {
	// Extrair o ID da requisição
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Extrair dados do body
	var dto dtos.AgendamentoUpdateStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualizar status do agendamento
	agendamento, err := h.agendamentoService.UpdateStatus(id, &dto)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Erro ao atualizar status do agendamento: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, agendamento)
}

// Update processa a atualização de um agendamento
func (h *AgendamentoHandler) Update(c *gin.Context) {
	// Extrair o ID da requisição
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Extrair dados do body
	var dto dtos.AgendamentoUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualizar agendamento
	agendamento, err := h.agendamentoService.Update(id, &dto)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Erro ao atualizar agendamento: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, agendamento)
}
