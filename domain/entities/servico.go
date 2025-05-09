package entities

import (
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type Servico struct {
	ID        ksuid.KSUID `gorm:"type:varchar(26);primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Nome      string      `json:"nome" gorm:"not null"`
	Descricao string      `json:"descricao"`
	Preco     float32     `json:"preco" gorm:"not null"`
	Duracao   string      `json:"duracao"`
	PetshopID ksuid.KSUID `json:"petshop_id" gorm:"type:varchar(26);not null"`
	PrecoBase float64     `json:"preco_base" gorm:"not null"`
	Ativo     bool        `json:"ativo" gorm:"default:true"`
}

// Antes de criar um registro o ID é gerado automaticamente
func (d *Servico) BeforeCreate(tx *gorm.DB) error {
	id, err := ksuid.NewRandom()
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}
