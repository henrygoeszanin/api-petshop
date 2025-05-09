package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/application/dtos"
	"github.com/henrygoeszanin/api_petshop/application/services"
	"github.com/segmentio/ksuid"
)

// ProcedimentoHandler gerencia as requisições relacionadas a procedimentos
type ProcedimentoHandler struct {
	procedimentoService *services.ProcedimentoService
}

// NewProcedimentoHandler cria uma nova instância de ProcedimentoHandler
func NewProcedimentoHandler(procedimentoService *services.ProcedimentoService) *ProcedimentoHandler {
	return &ProcedimentoHandler{
		procedimentoService: procedimentoService,
	}
}

// Create processa a criação de um novo registro de procedimento
func (h *ProcedimentoHandler) Create(c *gin.Context) {
	var dto dtos.ProcedimentoCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.procedimentoService.Create(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Erro ao criar procedimento: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetByPetID processa a requisição para listar os procedimentos de um pet
func (h *ProcedimentoHandler) GetByPetID(c *gin.Context) {
	// Extrair o ID do pet da requisição
	petIDStr := c.Param("petId")
	petID, err := ksuid.Parse(petIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do pet inválido"})
		return
	}

	// Buscar procedimentos do pet
	procedimentos, err := h.procedimentoService.GetByPetID(petID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar procedimentos: %v", err)})
		return
	}

	c.JSON(http.StatusOK, procedimentos)
}
