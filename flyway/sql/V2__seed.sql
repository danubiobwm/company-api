-- V2__seed.sql
INSERT INTO departamentos (id, nome, descricao)
VALUES
('018f3c3e-5c79-7b21-b7e1-d45f80cfa5ac', 'Tecnologia da Informação', 'Responsável por infraestrutura e sistemas internos'),
('018f3c3e-5c79-7b21-b7e1-d45f80cfa5ad', 'Recursos Humanos', 'Responsável por recrutamento e gestão de pessoas');

INSERT INTO colaboradores (id, nome, cpf, rg, departamento_id)
VALUES
('018f3c3e-5c79-7b21-b7e1-d45f80cfa6aa', 'João Silva', '00615075398', 'PR556677', '018f3c3e-5c79-7b21-b7e1-d45f80cfa5ac'),
('018f3c3e-5c79-7b21-b7e1-d45f80cfa6ab', 'Maria Oliveira', '12345678901', 'SP998877', '018f3c3e-5c79-7b21-b7e1-d45f80cfa5ad');

-- Define João Silva como gerente do TI
UPDATE departamentos
SET gerente_id = '018f3c3e-5c79-7b21-b7e1-d45f80cfa6aa'
WHERE id = '018f3c3e-5c79-7b21-b7e1-d45f80cfa5ac';
