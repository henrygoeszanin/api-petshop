package handlers

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/application/dtos"
	"github.com/henrygoeszanin/api_petshop/application/services"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
)

// AuthHandler gerencia as requisições relacionadas a autenticação
type AuthHandler struct {
	authService    *services.AuthService
	authMiddleware *jwt.GinJWTMiddleware
}

// NewAuthHandler cria uma nova instância de AuthHandler
func NewAuthHandler(authService *services.AuthService, authMiddleware *jwt.GinJWTMiddleware) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}

// Login processa as requisições de login
func (h *AuthHandler) Login(c *gin.Context) {
	h.authMiddleware.LoginHandler(c)
}

// Logout processa as requisições de logout
func (h *AuthHandler) Logout(c *gin.Context) {
	h.authMiddleware.LogoutHandler(c)
}

// RefreshToken atualiza o token JWT
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	h.authMiddleware.RefreshHandler(c)
}

// GetCurrentUser retorna informações do usuário autenticado
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get("identity")

	switch user.(type) {
	case *dtos.DonoResponseDTO:
		c.JSON(http.StatusOK, gin.H{
			"id":       claims["id"],
			"email":    claims["email"],
			"nome":     claims["nome"],
			"tipo":     claims["tipo"],
			"telefone": claims["telefone"],
		})
	case *dtos.PetshopResponseDTO:
		c.JSON(http.StatusOK, gin.H{
			"id":        claims["id"],
			"email":     claims["email"],
			"nome":      claims["nome"],
			"tipo":      claims["tipo"],
			"telefone":  claims["telefone"],
			"descricao": claims["descricao"],
			"nota":      claims["nota"],
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Tipo de usuário não reconhecido",
		})
	}
}

// RegisterDono processa o registro de um novo Dono
func (h *AuthHandler) RegisterDono(c *gin.Context) {
	var dto dtos.DonoRegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.authService.RegisterDono(&dto)
	if err != nil {
		switch err {
		case errors.ErrAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// RegisterPetshop processa o registro de um novo Petshop
func (h *AuthHandler) RegisterPetshop(c *gin.Context) {
	var dto dtos.PetshopRegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.authService.RegisterPetshop(&dto)
	if err != nil {
		switch err {
		case errors.ErrAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, resp)
}
