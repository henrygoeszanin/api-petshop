package dtos

// ItemProcedimentoCreateDTO representa um item de serviço para criação de um procedimento
type ItemProcedimentoCreateDTO struct {
	ServicoID  string  `json:"servico_id" binding:"required"`
	PrecoFinal float64 `json:"preco_final" binding:"required,min=0"`
}

// ProcedimentoCreateDTO representa dados para criação de um novo procedimento
type ProcedimentoCreateDTO struct {
	PetID          string                      `json:"pet_id" binding:"required"`
	PetshopID      string                      `json:"petshop_id" binding:"required"`
	DataRealizacao string                      `json:"data_realizacao" binding:"required"`
	Observacoes    string                      `json:"observacoes"`
	Total          float64                     `json:"total" binding:"required,min=0"`
	Itens          []ItemProcedimentoCreateDTO `json:"itens" binding:"required,min=1,dive"`
}

// ItemProcedimentoResponseDTO representa um item de serviço na resposta de um procedimento
type ItemProcedimentoResponseDTO struct {
	ID          string  `json:"id"`
	ServicoID   string  `json:"servico_id"`
	NomeServico string  `json:"nome_servico"`
	PrecoFinal  float64 `json:"preco_final"`
}

// ProcedimentoResponseDTO representa a estrutura de dados de resposta para um procedimento
type ProcedimentoResponseDTO struct {
	ID             string                        `json:"id"`
	PetID          string                        `json:"pet_id"`
	NomePet        string                        `json:"nome_pet,omitempty"`
	PetshopID      string                        `json:"petshop_id"`
	NomePetshop    string                        `json:"nome_petshop"`
	DataRealizacao string                        `json:"data_realizacao"`
	Observacoes    string                        `json:"observacoes,omitempty"`
	Total          float64                       `json:"total"`
	Itens          []ItemProcedimentoResponseDTO `json:"itens"`
	CreatedAt      string                        `json:"created_at"`
	UpdatedAt      string                        `json:"updated_at"`
}
