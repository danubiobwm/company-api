# Company API

API RESTful para gerenciamento de **Departamentos**, **Colaboradores** e **Gerentes**, desenvolvida em Go com **Gin**, **GORM**, e documentação automática via **Swagger**.

---

## Estrutura do Projeto

```
company-api/
│
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repositories/
│   └── services/
├── docs/
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Executando o Projeto

### Pré-requisitos

- Go 1.21+
- Docker e Docker Compose
- Swag CLI (`go install github.com/swaggo/swag/cmd/swag@latest`)

---

### Executar localmente

```bash
make tidy
make build
make run
```

API: `http://localhost:8080`

---

### Swagger

```bash
make swag
```

`http://localhost:8080/swagger/index.html`

---

### Testes

```bash
make test
make cover
```

---

### Docker

```bash
make docker-build
make docker-up
make docker-down
```

---

### Variáveis de Ambiente

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=company
PORT=8080
```

---

## Licença

MIT License

##
Danubio de araujo
