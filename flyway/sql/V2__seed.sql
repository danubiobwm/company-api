-- Inserir departamentos
INSERT INTO departamentos (nome, descricao)
VALUES
  ('Tecnologia da Informação', 'Responsável por infraestrutura, sistemas e suporte.'),
  ('Recursos Humanos', 'Gestão de pessoas e processos seletivos.'),
  ('Financeiro', 'Controle de fluxo de caixa, pagamentos e investimentos.'),
  ('Marketing', 'Divulgação e campanhas da empresa.');

-- Inserir colaboradores
INSERT INTO colaboradores (nome, cpf, rg, cargo, salario, departamento_id)
VALUES
  ('Danubio Martins', '00615075398', 'PR987654', 'Arquiteto de Software', 12000.00,
    (SELECT id FROM departamentos WHERE nome = 'Tecnologia da Informação' LIMIT 1)),
  ('João Silva', '11122233344', 'SP123456', 'Analista de RH', 5500.00,
    (SELECT id FROM departamentos WHERE nome = 'Recursos Humanos' LIMIT 1)),
  ('Maria Souza', '22233344455', 'RJ654321', 'Analista Financeiro', 6000.00,
    (SELECT id FROM departamentos WHERE nome = 'Financeiro' LIMIT 1));

-- Inserir gerente
INSERT INTO gerentes (colaborador_id, departamento_id)
VALUES (
  (SELECT id FROM colaboradores WHERE cpf = '00615075398' LIMIT 1),
  (SELECT id FROM departamentos WHERE nome = 'Tecnologia da Informação' LIMIT 1)
);
