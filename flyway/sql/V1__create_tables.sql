-- V1__create_tables.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS departamentos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    gerente_id UUID NULL,
    departamento_superior_id UUID NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS colaboradores (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(100) NOT NULL,
    cpf VARCHAR(11) NOT NULL UNIQUE,
    rg VARCHAR(50) UNIQUE,
    departamento_id UUID NOT NULL REFERENCES departamentos(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Constraint opcional de gerente (FK recursiva)
ALTER TABLE departamentos
    ADD CONSTRAINT fk_departamento_gerente FOREIGN KEY (gerente_id) REFERENCES colaboradores(id) ON DELETE SET NULL;

-- Constraint opcional de departamento superior (FK recursiva)
ALTER TABLE departamentos
    ADD CONSTRAINT fk_departamento_superior FOREIGN KEY (departamento_superior_id) REFERENCES departamentos(id) ON DELETE SET NULL;
