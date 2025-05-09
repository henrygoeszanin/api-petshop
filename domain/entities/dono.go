package entities

import (
	"time"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Dono struct {
	ID        ksuid.KSUID `gorm:"type:varchar(27);primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Nome        string `json:"nome" gorm:"not null"`
	Email       string `json:"email" gorm:"not null;unique"`
	Telefone    string `json:"telefone" gorm:"not null"`
	Password    string `json:"-" gorm:"not null"`
	CEP         string `json:"cep" gorm:"not null"`
	Rua         string `json:"rua" gorm:"not null"`
	Bairro      string `json:"bairro" gorm:"not null"`
	Cidade      string `json:"cidade" gorm:"not null"`
	Estado      string `json:"estado" gorm:"not null"`
	Complemento string `json:"complemento"`
	Numero      string `json:"numero" gorm:"not null"`
}

// Antes de criar um registro o ID é gerado automaticamente
func (d *Dono) BeforeCreate(tx *gorm.DB) error {
	id, err := ksuid.NewRandom()
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}

// SetPassword gera um hash da senha para armazenamento seguro
func (d *Dono) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	d.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifica se a senha fornecida corresponde à senha armazenada
func (d *Dono) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(d.Password), []byte(password))
	return err == nil
}
