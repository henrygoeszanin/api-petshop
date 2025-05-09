package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/application/dtos"
	"github.com/henrygoeszanin/api_petshop/application/services"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
)

// PetshopHandler gerencia as requisições relacionadas a petshops
type PetshopHandler struct {
	petshopService *services.PetshopService
}

// NewPetshopHandler cria uma nova instância de PetshopHandler
func NewPetshopHandler(petshopService *services.PetshopService) *PetshopHandler {
	return &PetshopHandler{
		petshopService: petshopService,
	}
}

// Create processa a criação de um novo petshop
func (h *PetshopHandler) Create(c *gin.Context) {
	var dto dtos.PetshopCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.petshopService.Create(&dto)
	if err != nil {
		switch err {
		case errors.ErrAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao criar petshop: %v", err)})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetByID processa a requisição para buscar um petshop por ID
func (h *PetshopHandler) GetByID(c *gin.Context) {
	// Extrair o ID da requisição
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Buscar petshop no serviço
	petshop, err := h.petshopService.GetByID(id)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Petshop não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar petshop: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, petshop)
}

// Update processa a atualização dos dados básicos de um petshop
func (h *PetshopHandler) Update(c *gin.Context) {
	// Extrair o ID da requisição
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Extrair dados do body
	var dto dtos.PetshopUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualizar petshop no serviço
	petshop, err := h.petshopService.Update(id, &dto)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Petshop não encontrado"})
		case errors.ErrAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao atualizar petshop: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, petshop)
}

// UpdateEndereco processa a atualização do endereço de um petshop
func (h *PetshopHandler) UpdateEndereco(c *gin.Context) {
	// Extrair o ID da requisição
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Extrair dados do body
	var dto dtos.PetshopUpdateEnderecoDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualizar endereço do petshop no serviço
	petshop, err := h.petshopService.UpdateEndereco(id, &dto)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Petshop não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao atualizar endereço do petshop: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, petshop)
}

// FindByCity busca petshops por cidade
func (h *PetshopHandler) FindByCity(c *gin.Context) {
	// Extrair parâmetros da query
	cidade := c.Query("cidade")
	if cidade == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'cidade' é obrigatório"})
		return
	}

	// Extrair parâmetros de paginação
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Buscar petshops por cidade
	petshops, err := h.petshopService.FindByCity(cidade, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar petshops: %v", err)})
		return
	}

	c.JSON(http.StatusOK, petshops)
}

// List lista todos os petshops com paginação
func (h *PetshopHandler) List(c *gin.Context) {
	// Extrair parâmetros de paginação
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Listar petshops
	petshops, err := h.petshopService.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao listar petshops: %v", err)})
		return
	}

	c.JSON(http.StatusOK, petshops)
}
