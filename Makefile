# Makefile para facilitar a gestão do ambiente Docker Go

# --- Configurações ---
.DEFAULT_GOAL := help

# Nome do serviço principal da aplicação no docker-compose.yml
SERVICE_NAME := zod-api

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

# Garante que os comandos funcionem mesmo que exista um arquivo com o mesmo nome
.PHONY: help dev debug-app debug-test up down stop logs ps destroy build setup-dev clean-dev test lint go-mod-tidy sec-check format stats

# --- Gerenciamento de Ambiente ---
setup-dev: ## Copia os arquivos de _env/dev para a raiz e cria o diretório tmp.
	@echo "🔧  Configurando o ambiente de desenvolvimento local..."
	@for file in $(DEV_FILES); do \
		if [ -f "$(DEV_ENV_DIR)/$${file}" ]; then \
			echo "Copiando $${file}..."; \
			cp "$(DEV_ENV_DIR)/$${file}" .; \
		else \
			echo "Aviso: Arquivo $${file} não encontrado em $(DEV_ENV_DIR), pulando."; \
		fi \
	done
	@echo "Copiando .env.example para .env com UID/GID dinâmico..."
	@if [ -f "$(DEV_ENV_DIR)/.env.example" ]; then \
		export HOST_UID=$(shell id -u); \
		export HOST_GID=$(shell id -g); \
		cp "$(DEV_ENV_DIR)/.env.example" .env; \
		sed -i -e "s/^APP_HOST_UID=.*/APP_HOST_UID=$${HOST_UID}/" \
			   -e "s/^APP_HOST_GID=.*/APP_HOST_GID=$${HOST_GID}/" .env; \
		echo "HOST_UID=$${HOST_UID}" >> .env; \
		echo "HOST_GID=$${HOST_GID}" >> .env; \
		echo "Arquivo .env configurado com HOST_UID=$${HOST_UID} e HOST_GID=$${HOST_GID}."; \
	else \
		echo "Aviso: Arquivo .env.example não encontrado em $(DEV_ENV_DIR), pulando criação do .env."; \
	fi
	@echo "Criando diretório tmp..."
	@mkdir -p tmp
	@echo "Ambiente local configurado com sucesso."

clean-dev: ## Remove os arquivos de configuração da raiz e o diretório tmp.
	@echo "🗑️  Limpando o ambiente de desenvolvimento local..."
	@echo "Removendo arquivos de configuração..."
	@rm -f $(DEV_FILES) .env
	@echo "Removendo diretório tmp..."
	@rm -rf tmp
	@echo "Ambiente local limpo com sucesso."

# --- Comandos Principais ---

dev: ## Configura e inicia o ambiente de desenvolvimento com Hot-Reload (Air)
	@$(MAKE) setup-dev
	@echo "🚀  Iniciando ambiente de desenvolvimento com Hot-Reload..."
	@docker compose up --build

debug-app: ## Configura e inicia o ambiente para DEPURAR A APLICAÇÃO com Delve
	@$(MAKE) setup-dev
	@echo "🐞  Iniciando ambiente de depuração da aplicação (use 'Attach to Docker' no VS Code)..."
	@docker compose -f compose.debug.yaml up --build

debug-test: ## Configura e inicia o ambiente para DEPURAR OS TESTES com Delve
	@$(MAKE) setup-dev
	@echo "🧪  Iniciando ambiente de depuração de testes (use 'Attach to Docker' no VS Code)..."
	@docker compose -f compose.debug-test.yaml up --build

# --- Comandos de Ciclo de Vida ---

up: ## Configura o ambiente e sobe os serviços em background
	@$(MAKE) setup-dev
	@echo "⬆️  Subindo serviços em background..."
	@docker compose up -d --build

down: ## Para e remove os contêineres (sem limpar os arquivos de ambiente)
	@echo "⬇️  Parando e removendo contêineres..."
	@docker compose down
	@if [ -f compose.debug.yaml ]; then docker compose -f compose.debug.yaml down; fi
	@if [ -f compose.debug-test.yaml ]; then docker compose -f compose.debug-test.yaml down; fi

stop: ## Para os contêineres de todos os ambientes sem removê-los
	@echo "🛑  Parando contêineres..."
	@docker compose stop
	@if [ -f compose.debug.yaml ]; then docker compose -f compose.debug.yaml stop; fi
	@if [ -f compose.debug-test.yaml ]; then docker compose -f compose.debug-test.yaml stop; fi

destroy: ## Destrói tudo (contêineres, volumes) e limpa o ambiente local
	@echo "🔥  Destruindo tudo (contêineres, volumes anônimos)..."
	@docker compose down -v --remove-orphans
	@if [ -f compose.debug.yaml ]; then docker compose -f compose.debug.yaml down -v --remove-orphans; fi
	@if [ -f compose.debug-test.yaml ]; then docker compose -f compose.debug-test.yaml down -v --remove-orphans; fi
	@$(MAKE) clean-dev

# --- Comandos de Desenvolvimento (Go) ---

test: ## Executa os testes unitários dentro do contêiner
	@echo "🧪  Executando testes..."
	@docker compose exec $(SERVICE_NAME) go test -v ./...

lint: ## Executa o linter (golangci-lint) dentro do contêiner
	@echo "🔍  Analisando o código com o linter..."
	@docker compose exec $(SERVICE_NAME) golangci-lint run

go-mod-tidy: ## Organiza as dependências do projeto (go mod tidy)
	@echo "📦  Organizando dependências do Go..."
	@docker compose run --rm -v .:/app:z $(SERVICE_NAME) go mod tidy

sec-check: ## Executa a análise de vulnerabilidades de segurança (gosec).
	@echo "🔐  Verificando vulnerabilidades de segurança..."
	@docker compose exec $(SERVICE_NAME) gosec ./...

format: ## Formata o código Go com regras mais rígidas (gofumpt).
	@echo "🎨  Formatando o código..."
	@docker compose exec $(SERVICE_NAME) gofumpt -w .

# --- Comandos Utilitários ---

build: ## Força a reconstrução das imagens sem iniciar os contêineres
	@echo "🛠️  Construindo imagens..."
	@docker compose build
	@docker compose -f compose.debug.yaml build
	@docker compose -f compose.debug-test.yaml build

logs: ## Exibe os logs dos serviços do docker-compose.yml
	@echo "📜  Exibindo logs..."
	@docker compose logs -f

ps: ## Lista os contêineres em execução
	@echo "📋  Listando contêineres..."
	@docker compose ps

stats: ## Exibe o uso de recursos (CPU, Memória) dos contêineres em tempo real.
	@echo "📊  Exibindo estatísticas de uso de recursos..."
	@docker compose stats

help: ## Exibe esta mensagem de ajuda
	@echo "Commands available:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

