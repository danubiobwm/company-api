# Makefile - Company API

APP_NAME=company-api
MAIN_FILE=cmd/api/main.go
BINARY=bin/$(APP_NAME)
SWAGGER_DIR=docs
COVERAGE_FILE=coverage.out

GO=go
SWAG=swag

DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_COMPOSE=docker-compose

.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  make tidy             - Atualiza dependências"
	@echo "  make build            - Compila o projeto"
	@echo "  make run              - Executa a aplicação localmente"
	@echo "  make test             - Executa testes unitários com cobertura"
	@echo "  make cover            - Mostra relatório de cobertura no navegador"
	@echo "  make swag             - Gera documentação Swagger"
	@echo "  make docker-build     - Build da imagem Docker"
	@echo "  make docker-up        - Sobe containers com Docker Compose"
	@echo "  make docker-down      - Para containers Docker"
	@echo "  make clean            - Remove binários e arquivos temporários"

tidy:
	$(GO) mod tidy

build:
	$(GO) build -o $(BINARY) $(MAIN_FILE)

run:
	$(GO) run $(MAIN_FILE)

test:
	$(GO) test ./... -coverprofile=$(COVERAGE_FILE)

cover:
	$(GO) tool cover -html=$(COVERAGE_FILE)

swag:
	$(SWAG) init -g $(MAIN_FILE) -o $(SWAGGER_DIR)

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

clean:
	rm -rf $(BINARY) $(COVERAGE_FILE)
