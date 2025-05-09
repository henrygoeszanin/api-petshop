package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/presentation/handlers"
	"github.com/henrygoeszanin/api_petshop/presentation/middlewares"
)

// SetupAuthRoutes configura as rotas de autenticação
func SetupAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, authMiddleware *jwt.GinJWTMiddleware) {
	// Adicionando o middleware TokenExtractor para toda a aplicação
	router.Use(middlewares.TokenExtractor())

	// Grupo de rotas para autenticação
	auth := router.Group("/auth")
	{
		// Rotas públicas (não requerem autenticação)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		// Rotas de registro de novas contas
		auth.POST("/register/dono", authHandler.RegisterDono)
		auth.POST("/register/petshop", authHandler.RegisterPetshop)
		auth.GET("/logout", authHandler.Logout)

		// Rotas protegidas (requerem autenticação)
		protected := auth.Group("/")
		protected.Use(authMiddleware.MiddlewareFunc())
		{
			protected.GET("/me", authHandler.GetCurrentUser)
		}
	}
}
