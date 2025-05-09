package entities

import (
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type Pet struct {
	ID        ksuid.KSUID `gorm:"type:varchar(27);primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Nome       string      `json:"nome" gorm:"not null"`
	Especie    string      `json:"especie" gorm:"not null"`
	Raca       string      `json:"raca" gorm:"not null"`
	Nascimento string      `json:"nascimento" gorm:"not null"`
	DonoID     ksuid.KSUID `json:"dono_id" gorm:"type:varchar(27);not null"`
}

// Antes de criar um registro o ID é gerado automaticamente
func (d *Pet) BeforeCreate(tx *gorm.DB) error {
	id, err := ksuid.NewRandom()
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}
