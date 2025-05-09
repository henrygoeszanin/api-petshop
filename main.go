package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/application/services"
	"github.com/henrygoeszanin/api_petshop/config"
	"github.com/henrygoeszanin/api_petshop/infrastructure/database"
	"github.com/henrygoeszanin/api_petshop/infrastructure/repositories"
	"github.com/henrygoeszanin/api_petshop/presentation/handlers"
	"github.com/henrygoeszanin/api_petshop/presentation/middlewares"
	"github.com/henrygoeszanin/api_petshop/presentation/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Carrega as configurações
	cfg := config.LoadConfig()

	// Configura o banco de dados
	db, err := database.SetupDatabase(cfg)
	if err != nil {
		fmt.Printf("Erro ao configurar banco de dados: %v\n", err)
		os.Exit(1)
	} // Inicializa os repositórios
	donoRepo := repositories.NewDonoRepository(db)
	petshopRepo := repositories.NewPetshopRepository(db)
	petRepo := repositories.NewPetRepository(db)

	// Inicializa os serviços
	authService := services.NewAuthService(donoRepo, petshopRepo)
	donoService := services.NewDonoService(donoRepo)
	petService := services.NewPetService(petRepo, donoRepo)

	// Configura o middleware JWT
	authMiddleware, err := middlewares.SetupJWTMiddleware(authService, cfg)
	if err != nil {
		fmt.Printf("Erro ao configurar middleware JWT: %v\n", err)
		os.Exit(1)
	}
	// Inicializa os handlers
	authHandler := handlers.NewAuthHandler(authService, authMiddleware)
	donoHandler := handlers.NewDonoHandler(donoService)
	petHandler := handlers.NewPetHandler(petService)
	// Configura o router Gin
	router := gin.Default()

	// Configura as rotas
	routes.SetupAuthRoutes(router, authHandler, authMiddleware)
	routes.SetupDonoRoutes(router, donoHandler, authMiddleware)
	routes.SetupPetRoutes(router, petHandler, authMiddleware)

	// Inicia o servidor
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("Servidor iniciado em http://localhost%s\n", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		fmt.Printf("Erro ao iniciar servidor: %v\n", err)
		os.Exit(1)
	}
}
