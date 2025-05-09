package middlewares

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

// DonoOwnershipRequired verifica se o ID do dono autenticado corresponde ao ID do recurso acessado
// ou se o usuário é um administrador. Caso contrário, aborta a requisição com 403 Forbidden.
func DonoOwnershipRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai as claims do token JWT
		claims := jwt.ExtractClaims(c)

		// Verifica se é admin (eles podem acessar qualquer recurso)
		isAdmin, exists := claims["is_admin"]
		if exists && isAdmin == true {
			c.Next()
			return
		}

		// Verifica o tipo de usuário
		tipo, tipoExists := claims["tipo"]
		if !tipoExists || tipo != "dono" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acesso negado. Apenas donos podem acessar este recurso.",
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

		// Extrai o ID do recurso solicitado
		resourceIDStr := c.Param("id")

		// Converte os IDs para ksuid.KSUID e compara
		donoID, err1 := ksuid.Parse(donoIDStr)
		resourceID, err2 := ksuid.Parse(resourceIDStr)

		if err1 != nil || err2 != nil || donoID != resourceID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acesso negado. Você não tem permissão para acessar este recurso.",
			})
			c.Abort()
			return
		}

		// Se passar por todas as verificações, o dono está acessando seu próprio recurso
		c.Next()
	}
}
