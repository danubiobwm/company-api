# Company API (Go + Gin + GORM + Postgres + Flyway)

## Requisitos
- Docker & docker-compose
- Go 1.22 (se desejar rodar local sem Docker)
- swag (opcional, para gerar docs)

## Rodando com Docker (modo recomendado)
1. Copie `.env.example` para `.env` e ajuste se necessário.
2. Suba os containers (DB + Flyway migrate + API):
```bash
docker-compose up --build
