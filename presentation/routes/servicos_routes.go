package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/presentation/handlers"
	"github.com/henrygoeszanin/api_petshop/presentation/middlewares"
)

// SetupServicoRoutes configura as rotas para operações relacionadas a serviços
func SetupServicoRoutes(router *gin.Engine, servicoHandler *handlers.ServicoHandler, authMiddleware *jwt.GinJWTMiddleware) {
	// Rota para listar serviços de um petshop (pública)
	router.GET("/petshops/:id/servicos", servicoHandler.GetByPetshopID)

	// Rotas protegidas (requerem autenticação)
	protected := router.Group("/petshops/:petshopId/servicos")
	protected.Use(authMiddleware.MiddlewareFunc())
	{
		// Rota para adicionar um novo serviço a um petshop
		// Apenas o próprio petshop pode adicionar serviços
		protected.POST("", middlewares.PetshopOwnershipFromParamRequired("petshopId"), servicoHandler.Create)
	}

	// Rotas para manipular serviços individuais
	servicoProtected := router.Group("/servicos")
	servicoProtected.Use(authMiddleware.MiddlewareFunc())
	{
		// Rota para obter detalhes de um serviço específico
		servicoProtected.GET("/:id", servicoHandler.GetByID)

		// Rotas que exigem verificação de propriedade do serviço
		// Um middleware personalizado seria necessário para verificar se o petshop autenticado
		// é o dono do serviço que está sendo manipulado
		servicoProtected.PUT("/:id", middlewares.ServicoOwnershipRequired(), servicoHandler.Update)
		servicoProtected.DELETE("/:id", middlewares.ServicoOwnershipRequired(), servicoHandler.Delete)
	}
}
