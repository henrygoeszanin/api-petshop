package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/presentation/handlers"
	"github.com/henrygoeszanin/api_petshop/presentation/middlewares"
)

// SetupPetshopRoutes configura as rotas para operações relacionadas a petshops
func SetupPetshopRoutes(router *gin.Engine, petshopHandler *handlers.PetshopHandler, authMiddleware *jwt.GinJWTMiddleware) { // Grupo de rotas para petshops
	petshops := router.Group("/petshops")
	{
		// Rotas públicas (não requerem autenticação)
		{
			// Rota pública para criar um novo petshop
			petshops.POST("", petshopHandler.Create)

			// Rota pública para buscar petshops por cidade
			petshops.GET("", petshopHandler.FindByCity)

			// Rota pública para buscar petshop específico por ID
			petshops.GET("/:id", petshopHandler.GetByID)
		}

		// Rotas protegidas (requerem autenticação)
		protected := petshops.Group("/")
		protected.Use(authMiddleware.MiddlewareFunc())
		{
			// Verificação adicional para garantir que apenas o próprio petshop ou admin pode atualizar
			// Rota para atualizar dados do petshop
			protected.PUT("/:id", middlewares.PetshopOwnershipRequired(), petshopHandler.Update)

			// Rota para atualizar endereço do petshop
			protected.PUT("/:id/endereco", middlewares.PetshopOwnershipRequired(), petshopHandler.UpdateEndereco)
		}
	}
}
