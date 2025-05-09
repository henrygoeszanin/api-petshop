package entities

import (
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// StatusAgendamento representa o status atual do agendamento
type StatusAgendamento string

const (
	// StatusPendente é o status inicial de um agendamento recém-criado
	StatusPendente StatusAgendamento = "pendente"
	// StatusConfirmado é o status após o petshop confirmar o agendamento
	StatusConfirmado StatusAgendamento = "confirmado"
	// StatusCancelado é o status quando o agendamento é cancelado (pelo dono ou petshop)
	StatusCancelado StatusAgendamento = "cancelado"
	// StatusConcluido é o status quando o procedimento foi realizado
	StatusConcluido StatusAgendamento = "concluido"
)

// ItemAgendamento representa um serviço selecionado em um agendamento
type ItemAgendamento struct {
	ID            ksuid.KSUID `gorm:"type:varchar(27);primaryKey"`
	AgendamentoID ksuid.KSUID `gorm:"type:varchar(27);index"`
	ServicoID     ksuid.KSUID `gorm:"type:varchar(27);index"`
	NomeServico   string      `gorm:"type:varchar(100);not null"` // Snapshot do nome do serviço
	PrecoPrevisto float64     `gorm:"type:decimal(10,2);not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// Agendamento representa um agendamento de procedimento a ser realizado em um pet
type Agendamento struct {
	ID            ksuid.KSUID       `gorm:"type:varchar(27);primaryKey"`
	DonoID        ksuid.KSUID       `gorm:"type:varchar(27);index;not null"`
	PetID         ksuid.KSUID       `gorm:"type:varchar(27);index;not null"`
	PetshopID     ksuid.KSUID       `gorm:"type:varchar(27);index;not null"`
	DataAgendada  time.Time         `gorm:"not null;index"`
	Status        StatusAgendamento `gorm:"type:varchar(20);not null;default:'pendente'"`
	Observacoes   string            `gorm:"type:text"`
	TotalPrevisto float64           `gorm:"type:decimal(10,2);not null"`
	Itens         []ItemAgendamento `gorm:"foreignKey:AgendamentoID"` // Relação um para muitos
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate é chamado pelo GORM antes de criar um registro
func (a *Agendamento) BeforeCreate(tx *gorm.DB) error {
	a.ID = ksuid.New()
	if a.Status == "" {
		a.Status = StatusPendente
	}
	return nil
}

// BeforeCreate é chamado pelo GORM antes de criar um registro
func (i *ItemAgendamento) BeforeCreate(tx *gorm.DB) error {
	i.ID = ksuid.New()
	return nil
}
