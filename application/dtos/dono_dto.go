package dtos

import (
	"github.com/segmentio/ksuid"
)

// DonoCreateDTO representa a estrutura de dados para criação de um novo dono
type DonoCreateDTO struct {
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

// DonoUpdateDTO representa a estrutura de dados para atualização de dados básicos do dono
type DonoUpdateDTO struct {
	Nome     string `json:"nome" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Telefone string `json:"telefone" binding:"required"`
}

// DonoUpdateLocalizacaoDTO representa a estrutura de dados para atualização de localização de um dono
type DonoUpdateLocalizacaoDTO struct {
	CEP         string `json:"cep" binding:"required"`
	Rua         string `json:"rua" binding:"required"`
	Bairro      string `json:"bairro" binding:"required"`
	Cidade      string `json:"cidade" binding:"required"`
	Estado      string `json:"estado" binding:"required"`
	Numero      string `json:"numero" binding:"required"`
	Complemento string `json:"complemento"`
}

// DonoDetailDTO representa a estrutura de dados completa de um dono para resposta de API
type DonoDetailDTO struct {
	ID          ksuid.KSUID `json:"id"`
	Nome        string      `json:"nome"`
	Email       string      `json:"email"`
	Telefone    string      `json:"telefone"`
	CEP         string      `json:"cep"`
	Rua         string      `json:"rua"`
	Bairro      string      `json:"bairro"`
	Cidade      string      `json:"cidade"`
	Estado      string      `json:"estado"`
	Numero      string      `json:"numero"`
	Complemento string      `json:"complemento,omitempty"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}
