package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/presentation/handlers"
	"github.com/henrygoeszanin/api_petshop/presentation/middlewares"
)

// SetupAgendamentoRoutes configura as rotas para operações relacionadas a agendamentos
func SetupAgendamentoRoutes(router *gin.Engine, agendamentoHandler *handlers.AgendamentoHandler, authMiddleware *jwt.GinJWTMiddleware) {
	// Grupo de rotas para agendamentos
	agendamentos := router.Group("/agendamentos")
	{
		// Todas as rotas de agendamentos requerem autenticação
		protected := agendamentos.Group("/")
		protected.Use(authMiddleware.MiddlewareFunc())
		{ // POST /agendamentos - Criar um novo agendamento
			protected.POST("", agendamentoHandler.Create)

			// GET /agendamentos/:id - Buscar agendamento por ID
			// Requer verificação de propriedade (dono ou petshop associado)
			protected.GET("/:id", middlewares.AgendamentoOwnershipRequired(), agendamentoHandler.GetByID)

			// PUT /agendamentos/:id - Atualizar agendamento
			// Requer verificação de propriedade (dono ou petshop associado)
			agendamentos.PUT("/:id", middlewares.AgendamentoOwnershipRequired(), agendamentoHandler.Update)

			// PUT /agendamentos/:id/status - Atualizar status do agendamento
			// Requer verificação de propriedade (dono ou petshop associado)
			agendamentos.PUT("/:id/status", middlewares.AgendamentoOwnershipRequired(), agendamentoHandler.UpdateStatus)
		}
	}

	// Rota para listar agendamentos de um dono
	donos := router.Group("/donos")
	{
		// Rotas protegidas (requerem autenticação)
		protected := donos.Group("/")
		protected.Use(authMiddleware.MiddlewareFunc())
		{
			// GET /donos/:donoId/agendamentos - Listar todos os agendamentos de um dono
			// Middleware verifica se o usuário autenticado é o próprio dono
			protected.GET("/:id/agendamentos", middlewares.DonoOwnershipRequired(), agendamentoHandler.GetByDonoID)
		}
	}

	// Rota para listar agendamentos de um petshop
	petshops := router.Group("/petshops")
	{
		// Rotas protegidas (requerem autenticação)
		protected := petshops.Group("/")
		protected.Use(authMiddleware.MiddlewareFunc())
		{
			// GET /petshops/:petshopId/agendamentos - Listar todos os agendamentos de um petshop
			// Middleware verifica se o usuário autenticado é o próprio petshop
			protected.GET("/:petshopId/agendamentos", middlewares.PetshopOwnershipRequired(), agendamentoHandler.GetByPetshopID)
		}
	}
}
