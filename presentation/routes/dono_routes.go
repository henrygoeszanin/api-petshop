package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/presentation/handlers"
	"github.com/henrygoeszanin/api_petshop/presentation/middlewares"
)

// SetupDonoRoutes configura as rotas para operações relacionadas a donos
func SetupDonoRoutes(router *gin.Engine, donoHandler *handlers.DonoHandler, authMiddleware *jwt.GinJWTMiddleware) { // Grupo de rotas para donos
	donos := router.Group("/donos")
	{
		// Rotas públicas (não requerem autenticação)
		{
			// Rota pública para criar um novo dono (POST /donos)
			donos.POST("", donoHandler.Create)
		}

		// Rotas protegidas (requerem autenticação)
		protected := donos.Group("/")
		protected.Use(authMiddleware.MiddlewareFunc())
		{
			// Rota para atualizar dados do dono (PUT /donos/:id)
			// Verifica se o usuário é dono do recurso ou admin
			protected.PUT(":id", middlewares.DonoOwnershipRequired(), donoHandler.Update)

			// Rota para atualizar localização do dono (PUT /donos/:id/localizacao)
			// Verifica se o usuário é dono do recurso ou admin
			protected.PUT(":id/localizacao", middlewares.DonoOwnershipRequired(), donoHandler.UpdateLocalizacao)
		}
	}
}
