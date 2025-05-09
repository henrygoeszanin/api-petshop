package dtos

import (
	"github.com/segmentio/ksuid"
)

// PetCreateDTO representa a estrutura de dados para criação de um novo pet
type PetCreateDTO struct {
	Nome       string `json:"nome" binding:"required"`
	Especie    string `json:"especie" binding:"required"`
	Raca       string `json:"raca" binding:"required"`
	Nascimento string `json:"nascimento" binding:"required"`
	DonoID     string `json:"dono_id" binding:"required"`
}

// PetResponseDTO representa a estrutura de dados de resposta para um pet
type PetResponseDTO struct {
	ID         ksuid.KSUID `json:"id"`
	Nome       string      `json:"nome"`
	Especie    string      `json:"especie"`
	Raca       string      `json:"raca"`
	Nascimento string      `json:"nascimento"`
	DonoID     ksuid.KSUID `json:"dono_id"`
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  string      `json:"updated_at"`
}

// PetUpdateDTO representa a estrutura de dados para atualização de um pet
type PetUpdateDTO struct {
	Nome       string `json:"nome" binding:"required"`
	Especie    string `json:"especie" binding:"required"`
	Raca       string `json:"raca" binding:"required"`
	Nascimento string `json:"nascimento" binding:"required"`
}
