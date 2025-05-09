package dtos

// ItemAgendamentoCreateDTO representa um item de serviço para criação de um agendamento
type ItemAgendamentoCreateDTO struct {
	ServicoID     string  `json:"servico_id" binding:"required"`
	PrecoPrevisto float64 `json:"preco_previsto" binding:"required,min=0"`
}

// AgendamentoCreateDTO representa dados para criação de um novo agendamento
type AgendamentoCreateDTO struct {
	DonoID        string                     `json:"dono_id" binding:"required"`
	PetID         string                     `json:"pet_id" binding:"required"`
	PetshopID     string                     `json:"petshop_id" binding:"required"`
	DataAgendada  string                     `json:"data_agendada" binding:"required"` // Formato ISO8601: 2006-01-02T15:04:05Z07:00
	Observacoes   string                     `json:"observacoes"`
	TotalPrevisto float64                    `json:"total_previsto" binding:"required,min=0"`
	Itens         []ItemAgendamentoCreateDTO `json:"itens" binding:"required,dive"`
}

// ItemAgendamentoResponseDTO representa um item de serviço na resposta de um agendamento
type ItemAgendamentoResponseDTO struct {
	ID            string  `json:"id"`
	ServicoID     string  `json:"servico_id"`
	NomeServico   string  `json:"nome_servico"`
	PrecoPrevisto float64 `json:"preco_previsto"`
}

// AgendamentoResponseDTO representa a estrutura de dados de resposta para um agendamento
type AgendamentoResponseDTO struct {
	ID            string                       `json:"id"`
	DonoID        string                       `json:"dono_id"`
	NomeDono      string                       `json:"nome_dono"`
	PetID         string                       `json:"pet_id"`
	NomePet       string                       `json:"nome_pet"`
	PetshopID     string                       `json:"petshop_id"`
	NomePetshop   string                       `json:"nome_petshop"`
	DataAgendada  string                       `json:"data_agendada"`
	Status        string                       `json:"status"`
	Observacoes   string                       `json:"observacoes"`
	TotalPrevisto float64                      `json:"total_previsto"`
	Itens         []ItemAgendamentoResponseDTO `json:"itens"`
	CreatedAt     string                       `json:"created_at"`
	UpdatedAt     string                       `json:"updated_at"`
}

// AgendamentoUpdateStatusDTO representa dados para atualização do status de um agendamento
type AgendamentoUpdateStatusDTO struct {
	Status string `json:"status" binding:"required,oneof=pendente confirmado cancelado concluido"`
}

// AgendamentoUpdateDTO representa dados para atualização de um agendamento existente
type AgendamentoUpdateDTO struct {
	DataAgendada  string                     `json:"data_agendada" binding:"required"` // Formato ISO8601: 2006-01-02T15:04:05Z07:00
	Observacoes   string                     `json:"observacoes"`
	TotalPrevisto float64                    `json:"total_previsto" binding:"required,min=0"`
	Itens         []ItemAgendamentoCreateDTO `json:"itens" binding:"required,dive"`
}
