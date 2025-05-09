package entities

import (
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// ItemProcedimento representa um serviço executado em um procedimento
type ItemProcedimento struct {
	ID             ksuid.KSUID `gorm:"type:varchar(27);primaryKey"`
	ProcedimentoID ksuid.KSUID `gorm:"type:varchar(27);index"`
	ServicoID      ksuid.KSUID `gorm:"type:varchar(27);index"`
	NomeServico    string      `gorm:"type:varchar(100);not null"` // Snapshot do nome do serviço
	PrecoFinal     float64     `gorm:"type:decimal(10,2);not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// Procedimento representa um registro de atendimento/procedimento realizado em um pet
type Procedimento struct {
	ID             ksuid.KSUID        `gorm:"type:varchar(27);primaryKey"`
	PetID          ksuid.KSUID        `gorm:"type:varchar(27);index;not null"`
	PetshopID      ksuid.KSUID        `gorm:"type:varchar(27);index;not null"`
	NomePetshop    string             `gorm:"type:varchar(100);not null"` // Snapshot do nome do petshop
	DataRealizacao time.Time          `gorm:"not null"`
	Observacoes    string             `gorm:"type:text"`
	Total          float64            `gorm:"type:decimal(10,2);not null"`
	Itens          []ItemProcedimento `gorm:"foreignKey:ProcedimentoID"` // Relação um para muitos
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate é chamado pelo GORM antes de criar um registro
func (p *Procedimento) BeforeCreate(tx *gorm.DB) error {
	p.ID = ksuid.New()
	return nil
}

// BeforeCreate é chamado pelo GORM antes de criar um registro
func (i *ItemProcedimento) BeforeCreate(tx *gorm.DB) error {
	i.ID = ksuid.New()
	return nil
}
