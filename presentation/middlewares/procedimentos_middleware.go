package middlewares

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/application/services"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
	"github.com/segmentio/ksuid"
)

// serviço global para uso nos middlewares
var petServiceInstance *services.PetService

// SetPetService configura o serviço de pet para uso nos middlewares
func SetPetService(s *services.PetService) {
	petServiceInstance = s
}

// PetOwnershipRequired verifica se o dono autenticado é o proprietário do pet
func PetOwnershipRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai as claims do token JWT
		claims := jwt.ExtractClaims(c)

		// Verifica o tipo de usuário
		tipo, tipoExists := claims["tipo"]
		if !tipoExists || tipo != "dono" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acesso negado. Apenas donos de pets podem acessar este recurso.",
			})
			c.Abort()
			return
		}

		// Extrai o ID do dono autenticado
		donoIDStr, exists := claims["id"].(string)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Problema de autenticação. ID do usuário não encontrado.",
			})
			c.Abort()
			return
		}

		// Extrai o ID do pet da URL
		petIDStr := c.Param("petId")
		petID, err := ksuid.Parse(petIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID do pet inválido"})
			c.Abort()
			return
		}

		// Busca o pet para verificar seu proprietário
		pet, err := petServiceInstance.GetByID(petID)
		if err != nil {
			if err == errors.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Pet não encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar proprietário do pet"})
			}
			c.Abort()
			return
		}

		// Converte o ID do dono autenticado para KSUID para comparação
		donoID, err := ksuid.Parse(donoIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ID no token é inválido",
			})
			c.Abort()
			return
		}

		// Compara o ID do dono autenticado com o ID do dono do pet
		if donoID.String() != pet.DonoID.String() {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acesso negado. Você não é o dono deste pet.",
			})
			c.Abort()
			return
		}

		// Se passar por todas as verificações, o dono está acessando seu próprio pet
		c.Next()
	}
}
