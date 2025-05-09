package config

import (
	"os"
)

// Config contém todas as configurações da aplicação
type Config struct {
	// Configurações do banco de dados
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Configurações do servidor
	ServerPort string

	// Configurações de autenticação
	JWTSecret string
}

// LoadConfig carrega as configurações do ambiente
func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "petshop_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "chave_secreta_padrao"),
	}
}

// getEnv retorna a variável de ambiente ou o valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
