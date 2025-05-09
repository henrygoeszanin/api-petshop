package middlewares

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/application/services"
	"github.com/segmentio/ksuid"
)

// serviço global para uso nos middlewares
var servicoServiceInstance *services.ServicoService

// SetServicoService configura o serviço de serviço para uso nos middlewares
func SetServicoService(s *services.ServicoService) {
	servicoServiceInstance = s
}

// ServicoOwnershipRequired verifica se o petshop autenticado é o proprietário do serviço
func ServicoOwnershipRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai as claims do token JWT
		claims := jwt.ExtractClaims(c)

		// Verifica o tipo de usuário
		tipo, tipoExists := claims["tipo"]
		if !tipoExists || tipo != "petshop" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acesso negado. Apenas petshops podem gerenciar serviços.",
			})
			c.Abort()
			return
		}

		// Extrai o ID do petshop autenticado
		petshopIDStr, exists := claims["id"].(string)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Problema de autenticação. ID do usuário não encontrado.",
			})
			c.Abort()
			return
		}

		// Extrai o ID do serviço da URL
		servicoIDStr := c.Param("id")
		servicoID, err := ksuid.Parse(servicoIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço inválido"})
			c.Abort()
			return
		}

		// Busca o serviço para verificar seu proprietário
		servico, err := servicoServiceInstance.GetByID(servicoID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
			c.Abort()
			return
		}

		// Converte o ID do petshop autenticado para KSUID para comparação
		petshopID, err := ksuid.Parse(petshopIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ID no token é inválido",
			})
			c.Abort()
			return
		}

		// Compara o ID do petshop autenticado com o ID do petshop dono do serviço
		if petshopID != servico.PetshopID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acesso negado. Você não tem permissão para gerenciar este serviço.",
			})
			c.Abort()
			return
		}

		// Se passar por todas as verificações, o petshop está gerenciando seu próprio serviço
		c.Next()
	}
}

// PetshopOwnershipFromParamRequired verifica se o petshop autenticado corresponde ao ID do petshop na URL
func PetshopOwnershipFromParamRequired(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai as claims do token JWT
		claims := jwt.ExtractClaims(c)

		// Verifica o tipo de usuário
		tipo, tipoExists := claims["tipo"]
		if !tipoExists || tipo != "petshop" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acesso negado. Apenas petshops podem gerenciar serviços.",
			})
			c.Abort()
			return
		}

		// Extrai o ID do petshop autenticado
		petshopIDStr, exists := claims["id"].(string)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Problema de autenticação. ID do usuário não encontrado.",
			})
			c.Abort()
			return
		}

		// Extrai o ID do petshop da URL
		urlPetshopIDStr := c.Param(paramName)

		// Converte os IDs para KSUID e compara
		petshopID, err1 := ksuid.Parse(petshopIDStr)
		urlPetshopID, err2 := ksuid.Parse(urlPetshopIDStr)

		if err1 != nil || err2 != nil || petshopID != urlPetshopID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acesso negado. Você não tem permissão para gerenciar serviços deste petshop.",
			})
			c.Abort()
			return
		}

		// Se passar por todas as verificações, o petshop está gerenciando seus próprios serviços
		c.Next()
	}
}
