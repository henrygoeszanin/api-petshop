package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/presentation/handlers"
	"github.com/henrygoeszanin/api_petshop/presentation/middlewares"
)

// SetupPetRoutes configura as rotas para operações relacionadas a pets
func SetupPetRoutes(router *gin.Engine, petHandler *handlers.PetHandler, authMiddleware *jwt.GinJWTMiddleware) {
	// Grupo de rotas para pets
	pets := router.Group("/pets")
	{
		// Rotas protegidas (requerem autenticação)
		protected := pets.Group("/")
		protected.Use(authMiddleware.MiddlewareFunc())
		{
			// POST /pets - Criar um novo pet
			protected.POST("", petHandler.Create)

			// GET /pets/:id - Retornar dados do pet por ID
			protected.GET(":id", petHandler.GetByID)
		}
	}

	// Rota para listar pets de um dono
	donos := router.Group("/donos")
	{
		// Rotas protegidas (requerem autenticação)
		protected := donos.Group("/")
		protected.Use(authMiddleware.MiddlewareFunc())
		{
			// GET /donos/:id/pets - Listar todos os pets de um determinado dono
			// Verifica se o usuário é dono do recurso ou admin
			protected.GET(":id/pets", middlewares.DonoOwnershipRequired(), petHandler.GetByDonoID)
		}
	}
}
