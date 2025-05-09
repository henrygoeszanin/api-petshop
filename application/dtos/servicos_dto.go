package dtos

import (
	"github.com/segmentio/ksuid"
)

// ServicoCreateDTO representa a estrutura de dados para criação de um novo serviço
type ServicoCreateDTO struct {
	Nome      string  `json:"nome" binding:"required"`
	Descricao string  `json:"descricao" binding:"required"`
	PrecoBase float64 `json:"preco_base" binding:"required,min=0"`
}

// ServicoUpdateDTO representa a estrutura de dados para atualização de um serviço
type ServicoUpdateDTO struct {
	Nome      string  `json:"nome" binding:"required"`
	Descricao string  `json:"descricao" binding:"required"`
	PrecoBase float64 `json:"preco_base" binding:"required,min=0"`
}

// ServicoResponseDTO representa a estrutura de dados de resposta para um serviço
type ServicoResponseDTO struct {
	ID        ksuid.KSUID `json:"id"`
	PetshopID ksuid.KSUID `json:"petshop_id"`
	Nome      string      `json:"nome"`
	Descricao string      `json:"descricao"`
	PrecoBase float64     `json:"preco_base"`
	Ativo     bool        `json:"ativo"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}
