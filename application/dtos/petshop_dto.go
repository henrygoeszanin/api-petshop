package dtos

import (
	"github.com/segmentio/ksuid"
)

// PetshopCreateDTO representa a estrutura de dados para criação de um novo petshop
type PetshopCreateDTO struct {
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

// PetshopUpdateDTO representa a estrutura de dados para atualização básica de um petshop
type PetshopUpdateDTO struct {
	Nome      string `json:"nome" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Telefone  string `json:"telefone" binding:"required"`
	Descricao string `json:"descricao"`
}

// PetshopUpdateEnderecoDTO representa a estrutura de dados para atualização do endereço de um petshop
type PetshopUpdateEnderecoDTO struct {
	CEP         string `json:"cep" binding:"required"`
	Rua         string `json:"rua" binding:"required"`
	Bairro      string `json:"bairro" binding:"required"`
	Cidade      string `json:"cidade" binding:"required"`
	Estado      string `json:"estado" binding:"required"`
	Numero      string `json:"numero" binding:"required"`
	Complemento string `json:"complemento"`
}

// PetshopDetailDTO representa a estrutura de dados completa de um petshop para resposta de API
type PetshopDetailDTO struct {
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
	Descricao   string      `json:"descricao,omitempty"`
	Nota        float32     `json:"nota"`
	Ativo       bool        `json:"ativo"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

// PetshopListItemDTO representa a estrutura de dados resumida de um petshop para listagens
type PetshopListItemDTO struct {
	ID        ksuid.KSUID `json:"id"`
	Nome      string      `json:"nome"`
	Cidade    string      `json:"cidade"`
	Estado    string      `json:"estado"`
	Nota      float32     `json:"nota"`
	Descricao string      `json:"descricao,omitempty"`
}
