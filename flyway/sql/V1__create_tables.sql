-- V1__create_tables.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- departments table
CREATE TABLE IF NOT EXISTS departamentos (
    id UUID PRIMARY KEY,
    nome TEXT NOT NULL,
    gerente_id UUID NOT NULL,
    departamento_superior_id UUID,
    CONSTRAINT fk_departamento_superior FOREIGN KEY (departamento_superior_id) REFERENCES departamentos(id) ON DELETE SET NULL
);

-- colaboradores table
CREATE TABLE IF NOT EXISTS colaboradores (
    id UUID PRIMARY KEY,
    nome TEXT NOT NULL,
    cpf VARCHAR(11) NOT NULL UNIQUE,
    rg VARCHAR(50) UNIQUE,
    departamento_id UUID NOT NULL,
    CONSTRAINT fk_departamento FOREIGN KEY (departamento_id) REFERENCES departamentos(id) ON DELETE RESTRICT
);

-- constraint: gerente must be an existing colaborador linked to same department enforced by app logic
-- add indexes
CREATE INDEX IF NOT EXISTS idx_colaboradores_departamento ON colaboradores(departamento_id);
CREATE INDEX IF NOT EXISTS idx_departamentos_superior ON departamentos(departamento_superior_id);
