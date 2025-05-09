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

// ServicoHandler gerencia as requisições relacionadas a serviços
type ServicoHandler struct {
	servicoService *services.ServicoService
}

// NewServicoHandler cria uma nova instância de ServicoHandler
func NewServicoHandler(servicoService *services.ServicoService) *ServicoHandler {
	return &ServicoHandler{
		servicoService: servicoService,
	}
}

// Create processa a criação de um novo serviço para um petshop
func (h *ServicoHandler) Create(c *gin.Context) {
	// Extrair ID do petshop da URL
	petshopIDStr := c.Param("petshopId")
	petshopID, err := ksuid.Parse(petshopIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do petshop inválido"})
		return
	}

	// Extrair dados do body
	var dto dtos.ServicoCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Criar serviço
	response, err := h.servicoService.Create(petshopID, &dto)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Petshop não encontrado"})
		case errors.ErrAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Já existe um serviço com este nome neste petshop"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao criar serviço: %v", err)})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetByPetshopID lista todos os serviços de um determinado petshop
func (h *ServicoHandler) GetByPetshopID(c *gin.Context) {
	// Extrair ID do petshop da URL
	petshopIDStr := c.Param("petshopId")
	petshopID, err := ksuid.Parse(petshopIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do petshop inválido"})
		return
	}

	// Buscar serviços
	servicos, err := h.servicoService.GetByPetshopID(petshopID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Petshop não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao listar serviços: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, servicos)
}

// GetByID processa a busca de um serviço por ID
func (h *ServicoHandler) GetByID(c *gin.Context) {
	// Extrair ID do serviço da URL
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço inválido"})
		return
	}

	// Buscar serviço
	servico, err := h.servicoService.GetByID(id)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar serviço: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, servico)
}

// Update processa a atualização de um serviço
func (h *ServicoHandler) Update(c *gin.Context) {
	// Extrair ID do serviço da URL
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço inválido"})
		return
	}

	// Extrair dados do body
	var dto dtos.ServicoUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualizar serviço
	servico, err := h.servicoService.Update(id, &dto)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		case errors.ErrAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Já existe um serviço com este nome neste petshop"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao atualizar serviço: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, servico)
}

// Delete processa a exclusão de um serviço
func (h *ServicoHandler) Delete(c *gin.Context) {
	// Extrair ID do serviço da URL
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço inválido"})
		return
	}

	// Excluir serviço
	if err := h.servicoService.Delete(id); err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao excluir serviço: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Serviço excluído com sucesso"})
}
