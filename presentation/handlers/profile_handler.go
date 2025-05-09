package handlers

import (
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/application/services"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
)

// ProfileHandler gerencia as requisições relacionadas ao perfil do usuário logado
type ProfileHandler struct {
	donoService    *services.DonoService
	petshopService *services.PetshopService
}

// NewProfileHandler cria uma nova instância de ProfileHandler
func NewProfileHandler(donoService *services.DonoService, petshopService *services.PetshopService) *ProfileHandler {
	return &ProfileHandler{
		donoService:    donoService,
		petshopService: petshopService,
	}
}

// GetCurrentUserProfile obtém o perfil do usuário atualmente autenticado
func (h *ProfileHandler) GetCurrentUserProfile(c *gin.Context) {
	// Extrair claims do token JWT
	claims := jwt.ExtractClaims(c)

	// Obter ID e tipo do usuário
	idStr, idExists := claims["id"].(string)
	tipo, tipoExists := claims["tipo"].(string)

	// Verificar se os dados necessários estão presentes
	if !idExists || !tipoExists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token inválido ou expirado",
		})
		return
	}

	// Converter string ID para KSUID
	id, err := ksuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ID no token é inválido",
		})
		return
	}

	// Buscar perfil com base no tipo de usuário
	switch tipo {
	case "dono":
		dono, err := h.donoService.GetByID(id)
		if err != nil {
			switch err {
			case errors.ErrNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "Dono não encontrado"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Erro ao buscar perfil do dono: %v", err),
				})
			}
			return
		}
		c.JSON(http.StatusOK, dono)

	case "petshop":
		petshop, err := h.petshopService.GetByID(id)
		if err != nil {
			switch err {
			case errors.ErrNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "Petshop não encontrado"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Erro ao buscar perfil do petshop: %v", err),
				})
			}
			return
		}
		c.JSON(http.StatusOK, petshop)

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Tipo de usuário não reconhecido: %s", tipo),
		})
	}
}
