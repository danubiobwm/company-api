package models

import (
	"time"

	"github.com/google/uuid"
)

type Colaborador struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Nome           string    `gorm:"not null" json:"nome"`
	CPF            string    `gorm:"size:11;not null;uniqueIndex" json:"cpf"`
	RG             *string   `gorm:"uniqueIndex" json:"rg,omitempty"`
	DepartamentoID uuid.UUID `gorm:"type:uuid;not null" json:"departamento_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Departamento struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Nome                   string     `gorm:"not null" json:"nome"`
	GerenteID              uuid.UUID  `gorm:"type:uuid;not null" json:"gerente_id"`
	DepartamentoSuperiorID *uuid.UUID `gorm:"type:uuid" json:"departamento_superior_id,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}
