# Makefile para facilitar a gestão do ambiente Docker Go

# --- Configurações ---
.DEFAULT_GOAL := help
.SILENT:

# Nome do serviço principal da aplicação no docker-compose.yml
SERVICE_NAME := zod-api

COMPOSE_FILES := $(wildcard docker-compose.yml compose.debug.yaml compose.debug-test.yaml)
COMPOSE_ARGS := $(foreach file,$(COMPOSE_FILES),-f $(file))

COMPOSE_CMD := docker compose
GOOSE_CMD := $(COMPOSE_CMD) exec $(SERVICE_NAME) goose

# Lista de arquivos de configuração de desenvolvimento a serem gerenciados.
DEV_FILES := \
    compose.debug-test.yaml \
    compose.debug.yaml \
    docker-compose.yml \
    Dockerfile \
    Dockerfile.debug \
    Dockerfile.debug-test \
    requests.http

# Define o diretório de origem para os arquivos de configuração.
DEV_ENV_DIR := _env/dev

.PHONY: help dev debug-app debug-test up down stop logs ps destroy build setup-dev clean-dev test lint go-mod-tidy sec-check format stats shell migrate-up migrate-down migrate-reset migrate-create

# --- Gerenciamento de Ambiente ---
setup-dev: ## Copia os arquivos de _env/dev para a raiz e cria o diretório tmp.
	echo "🔧  Configurando o ambiente de desenvolvimento local..."
	for file in $(DEV_FILES); do \
		if [ -f "$(DEV_ENV_DIR)/$${file}" ]; then \
			echo "Copiando $${file}..."; \
			cp "$(DEV_ENV_DIR)/$${file}" .; \
		else \
			echo "Aviso: Arquivo $${file} não encontrado em $(DEV_ENV_DIR), pulando."; \
		fi \
	done
	echo "Copiando .env.example para .env com UID/GID dinâmico..."
	if [ -f "$(DEV_ENV_DIR)/.env.example" ]; then \
		export HOST_UID=$(shell id -u); \
		export HOST_GID=$(shell id -g); \
		sed -e "s/^APP_HOST_UID=.*/APP_HOST_UID=$${HOST_UID}/" \
			-e "s/^APP_HOST_GID=.*/APP_HOST_GID=$${HOST_GID}/" "$(DEV_ENV_DIR)/.env.example" > .env; \
		echo "Arquivo .env configurado com HOST_UID=$${HOST_UID} e HOST_GID=$${HOST_GID}."; \
	else \
		echo "Aviso: Arquivo .env.example não encontrado em $(DEV_ENV_DIR), pulando criação do .env."; \
	fi
	echo "Criando diretório tmp..."
	mkdir -p tmp
	echo "Ambiente local configurado com sucesso."

clean-dev: ## Remove os arquivos de configuração da raiz e o diretório tmp.
	echo "🗑️  Limpando o ambiente de desenvolvimento local..."
	echo "Removendo arquivos de configuração..."
	rm -f $(DEV_FILES) .env
	echo "Removendo diretório tmp..."
	rm -rf tmp
	echo "Ambiente local limpo com sucesso."

# --- Comandos Principais ---
dev: ## Configura e inicia o ambiente de desenvolvimento com Hot-Reload (Air)
	$(MAKE) setup-dev
	echo "🚀  Iniciando ambiente de desenvolvimento com Hot-Reload..."
	$(COMPOSE_CMD) up --build

debug-app: ## Configura e inicia o ambiente para DEPURAR A APLICAÇÃO com Delve
	$(MAKE) setup-dev
	echo "🐞  Iniciando ambiente de depuração da aplicação (use 'Attach to Docker' no VS Code)..."
	$(COMPOSE_CMD) -f compose.debug.yaml up --build

debug-test: ## Configura e inicia o ambiente para DEPURAR OS TESTES com Delve
	$(MAKE) setup-dev
	echo "🧪  Iniciando ambiente de depuração de testes (use 'Attach to Docker' no VS Code)..."
	$(COMPOSE_CMD) -f compose.debug-test.yaml up --build

# --- Comandos de Ciclo de Vida ---
up: ## Configura o ambiente e sobe os serviços em background
	$(MAKE) setup-dev
	echo "⬆️  Subindo serviços em background..."
	$(COMPOSE_CMD) up -d --build

down: ## Para e remove os contêineres (sem limpar os arquivos de ambiente)
	echo "⬇️  Parando e removendo contêineres..."
	$(COMPOSE_CMD) $(COMPOSE_ARGS) down

stop: ## Para os contêineres de todos os ambientes sem removê-los
	echo "🛑  Parando contêineres..."
	$(COMPOSE_CMD) $(COMPOSE_ARGS) stop

destroy: ## Destrói tudo (contêineres, volumes) e limpa o ambiente local
	echo "🔥  Destruindo tudo (contêineres, volumes anônimos)..."
	$(COMPOSE_CMD) $(COMPOSE_ARGS) down -v --remove-orphans
	$(MAKE) clean-dev

# --- Comandos de Desenvolvimento (Go) ---
test: ## Executa os testes unitários dentro do contêiner
	echo "🧪  Executando testes..."
	$(COMPOSE_CMD) exec $(SERVICE_NAME) go test -v ./...

lint: ## Executa o linter (golangci-lint) dentro do contêiner
	echo "🔍  Analisando o código com o linter..."
	$(COMPOSE_CMD) exec $(SERVICE_NAME) golangci-lint run

go-mod-tidy: ## Organiza as dependências do projeto (go mod tidy)
	echo "📦  Organizando dependências do Go..."
	$(COMPOSE_CMD) run --rm -v .:/app:z $(SERVICE_NAME) go mod tidy

sec-check: ## Executa a análise de vulnerabilidades de segurança (gosec).
	echo "🔐  Verificando vulnerabilidades de segurança..."
	$(COMPOSE_CMD) exec $(SERVICE_NAME) gosec ./...

format: ## Formata o código Go com regras mais rígidas (gofumpt).
	echo "🎨  Formatando o código..."
	$(COMPOSE_CMD) exec $(SERVICE_NAME) gofumpt -w .

# --- Comandos de Migração (Goose) ---
migrate-up: ## Executa as migrações pendentes (up)
	echo "Applying database migrations..."
	$(GOOSE_CMD) up

migrate-down: ## Reverte a última migração aplicada (down)
	echo "Reverting last database migration..."
	$(GOOSE_CMD) down

migrate-reset: ## Reverte todas as migrações (down to 0)
	echo "Reverting all database migrations..."
	$(GOOSE_CMD) reset

migrate-create: ## Cria um novo arquivo de migração SQL
	$(if $(name),,$(error "O argumento 'name' é obrigatório. Ex: make migrate-create name=add_users_table"))
	echo "INFO: Criando nova migração: $(name)..."
	$(GOOSE_CMD) create $(name) sql

# --- Comandos Utilitários ---
build: ## Força a reconstrução das imagens sem iniciar os contêineres
	echo "🛠️  Construindo imagens..."
	$(COMPOSE_CMD) $(COMPOSE_ARGS) build

logs: ## Exibe os logs dos serviços do docker-compose.yml
	echo "📜  Exibindo logs..."
	$(COMPOSE_CMD) logs -f

ps: ## Lista os contêineres em execução
	echo "📋  Listando contêineres..."
	$(COMPOSE_CMD) ps

stats: ## Exibe o uso de recursos (CPU, Memória) dos contêineres em tempo real.
	echo "📊  Exibindo estatísticas de uso de recursos..."
	$(COMPOSE_CMD) stats

shell: ## Acessa o shell (sh) do contêiner da aplicação
	echo "Acessando o shell do serviço $(SERVICE_NAME)..."
	$(COMPOSE_CMD) exec $(SERVICE_NAME) sh

help: ## Exibe esta mensagem de ajuda
	echo "Commands available:"
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
