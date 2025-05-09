package dtos

import "github.com/segmentio/ksuid"

// LoginRequestDTO representa a estrutura de dados para requisição de login
type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponseDTO representa a estrutura de dados comum para resposta após autenticação
type AuthResponseDTO struct {
	ID    ksuid.KSUID `json:"id"`
	Email string      `json:"email"`
	Nome  string      `json:"nome"`
	Tipo  string      `json:"tipo"` // "dono" ou "petshop"
}

// DonoResponseDTO representa a estrutura de dados para resposta após autenticação de um dono
type DonoResponseDTO struct {
	AuthResponseDTO
	Telefone string `json:"telefone"`
}

// PetshopResponseDTO representa a estrutura de dados para resposta após autenticação de um petshop
type PetshopResponseDTO struct {
	AuthResponseDTO
	Telefone  string  `json:"telefone"`
	Descricao string  `json:"descricao,omitempty"`
	Nota      float32 `json:"nota,omitempty"`
}

// DonoRegisterDTO representa a estrutura de dados para registro de um dono de pet
type DonoRegisterDTO struct {
	Nome        string `json:"nome" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Telefone    string `json:"telefone" binding:"required"`
	CEP         string `json:"cep" binding:"required"`
	Rua         string `json:"rua" binding:"required"`
	Bairro      string `json:"bairro" binding:"required"`
	Cidade      string `json:"cidade" binding:"required"`
	Estado      string `json:"estado" binding:"required"`
	Numero      string `json:"numero" binding:"required"`
	Complemento string `json:"complemento"`
}

// PetshopRegisterDTO representa a estrutura de dados para registro de um petshop
type PetshopRegisterDTO struct {
	Nome        string `json:"nome" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Telefone    string `json:"telefone" binding:"required"`
	CEP         string `json:"cep" binding:"required"`
	Rua         string `json:"rua" binding:"required"`
	Bairro      string `json:"bairro" binding:"required"`
	Cidade      string `json:"cidade" binding:"required"`
	Estado      string `json:"estado" binding:"required"`
	Numero      string `json:"numero" binding:"required"`
	Complemento string `json:"complemento"`
	Descricao   string `json:"descricao"`
}
