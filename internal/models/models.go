package models

import (
	"time"

	"github.com/google/uuid"
)

type Colaborador struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Nome           string    `gorm:"not null" json:"nome"`
	CPF            string    `gorm:"size:11;not null;uniqueIndex" json:"cpf"`
	RG             *string   `gorm:"size:50;uniqueIndex" json:"rg,omitempty"`
	DepartamentoID uuid.UUID `gorm:"type:uuid;not null" json:"departamento_id"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Departamento struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Nome                   string     `gorm:"not null" json:"nome"`
	Descricao              *string    `gorm:"type:text" json:"descricao,omitempty"`
	GerenteID              *uuid.UUID `gorm:"type:uuid" json:"gerente_id,omitempty"`
	DepartamentoSuperiorID *uuid.UUID `gorm:"type:uuid" json:"departamento_superior_id,omitempty"`
	CreatedAt              time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// relations (for preload)
	Gerente *Colaborador `gorm:"foreignKey:GerenteID" json:"gerente,omitempty"`
}

func (Colaborador) TableName() string  { return "colaboradores" }
func (Departamento) TableName() string { return "departamentos" }
