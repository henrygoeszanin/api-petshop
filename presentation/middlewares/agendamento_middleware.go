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
var agendamentoServiceInstance *services.AgendamentoService

// SetAgendamentoService configura o serviço de agendamento para uso nos middlewares
func SetAgendamentoService(s *services.AgendamentoService) {
	agendamentoServiceInstance = s
}

// AgendamentoOwnershipRequired verifica se o usuário autenticado é o dono do agendamento ou o petshop associado
func AgendamentoOwnershipRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar se o serviço foi configurado
		if agendamentoServiceInstance == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Serviço de agendamento não configurado"})
			c.Abort()
			return
		}

		// Extrai as claims do token JWT
		claims := jwt.ExtractClaims(c)

		// Verifica o tipo de usuário
		tipo, tipoExists := claims["tipo"]
		if !tipoExists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Tipo de usuário não especificado no token"})
			c.Abort()
			return
		}

		// Extrai o ID do agendamento da URL
		agendamentoIDStr := c.Param("id")
		agendamentoID, err := ksuid.Parse(agendamentoIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID do agendamento inválido"})
			c.Abort()
			return
		}

		// Busca o agendamento
		agendamento, err := agendamentoServiceInstance.GetByID(agendamentoID)
		if err != nil {
			if err == errors.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar agendamento"})
			}
			c.Abort()
			return
		}

		// Extrai o ID do usuário autenticado
		userIDStr, exists := claims["id"].(string)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "ID do usuário não encontrado no token"})
			c.Abort()
			return
		}

		// Verifica o acesso com base no tipo de usuário
		if tipo == "dono" {
			// Se for dono, verifica se o agendamento pertence a ele
			if userIDStr != agendamento.DonoID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para acessar este agendamento"})
				c.Abort()
				return
			}
		} else if tipo == "petshop" {
			// Se for petshop, verifica se o agendamento está associado a ele
			if userIDStr != agendamento.PetshopID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Este agendamento não pertence ao seu petshop"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Tipo de usuário não autorizado"})
			c.Abort()
			return
		}

		// Se passou por todas as verificações, o usuário pode acessar o agendamento
		c.Next()
	}
}
