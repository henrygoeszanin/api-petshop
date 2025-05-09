package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/presentation/handlers"
)

// SetupProfileRoutes configura as rotas para operações relacionadas ao perfil do usuário
func SetupProfileRoutes(router *gin.Engine, profileHandler *handlers.ProfileHandler, authMiddleware *jwt.GinJWTMiddleware) {
	// Grupo de rotas para perfil
	profile := router.Group("/profile")

	// Todas as rotas de perfil requerem autenticação
	profile.Use(authMiddleware.MiddlewareFunc())

	// Rota para obter o próprio perfil do usuário logado
	profile.GET("", profileHandler.GetCurrentUserProfile)
}
