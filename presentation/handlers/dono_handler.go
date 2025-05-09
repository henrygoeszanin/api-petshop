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

// DonoHandler gerencia as requisições relacionadas a donos
type DonoHandler struct {
	donoService *services.DonoService
}

// NewDonoHandler cria uma nova instância de DonoHandler
func NewDonoHandler(donoService *services.DonoService) *DonoHandler {
	return &DonoHandler{
		donoService: donoService,
	}
}

// Create processa a criação de um novo dono
func (h *DonoHandler) Create(c *gin.Context) {
	var dto dtos.DonoCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.donoService.Create(&dto)
	if err != nil {
		switch err {
		case errors.ErrAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao criar dono: %v", err)})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetByID processa a requisição para buscar um dono por ID
func (h *DonoHandler) GetByID(c *gin.Context) {
	// Extrair o ID da requisição
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Buscar dono no serviço
	dono, err := h.donoService.GetByID(id)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Dono não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar dono: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, dono)
}

// Update processa a atualização dos dados básicos de um dono
func (h *DonoHandler) Update(c *gin.Context) {
	// Extrair o ID da requisição
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Extrair dados do body
	var dto dtos.DonoUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualizar dono no serviço
	dono, err := h.donoService.Update(id, &dto)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Dono não encontrado"})
		case errors.ErrAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao atualizar dono: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, dono)
}

// UpdateLocalizacao processa a atualização da localização de um dono
func (h *DonoHandler) UpdateLocalizacao(c *gin.Context) {
	// Extrair o ID da requisição
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Extrair dados do body
	var dto dtos.DonoUpdateLocalizacaoDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualizar localização do dono no serviço
	dono, err := h.donoService.UpdateLocalizacao(id, &dto)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Dono não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao atualizar localização do dono: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, dono)
}
