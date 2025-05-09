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

// PetHandler gerencia as requisições relacionadas a pets
type PetHandler struct {
	petService *services.PetService
}

// NewPetHandler cria uma nova instância de PetHandler
func NewPetHandler(petService *services.PetService) *PetHandler {
	return &PetHandler{
		petService: petService,
	}
}

// Create processa a criação de um novo pet
func (h *PetHandler) Create(c *gin.Context) {
	var dto dtos.PetCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.petService.Create(&dto)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Dono não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao criar pet: %v", err)})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetByID processa a requisição para buscar um pet por ID
func (h *PetHandler) GetByID(c *gin.Context) {
	// Extrair o ID da requisição
	idStr := c.Param("id")
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Buscar pet no serviço
	pet, err := h.petService.GetByID(id)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Pet não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar pet: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, pet)
}

// GetByDonoID processa a requisição para listar todos os pets de um dono
func (h *PetHandler) GetByDonoID(c *gin.Context) {
	// Extrair o ID do dono da requisição
	idStr := c.Param("id")
	donoID, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do dono inválido"})
		return
	}

	// Buscar pets do dono no serviço
	pets, err := h.petService.GetByDonoID(donoID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Dono não encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar pets do dono: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, pets)
}
