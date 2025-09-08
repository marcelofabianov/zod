# ZOD

ZOD é um projeto de API em Go para desenvolvimento local com suporte a hot-reload, usando Docker e Docker Compose. Ele integra PostgreSQL (DB), RabbitMQ e Redis para um ambiente completo de microsserviço. O gerenciamento é facilitado por um Makefile que automatiza a configuração, build, execução e comandos de desenvolvimento.

## Pré-requisitos

- Docker: Versão 20+ com Docker Compose (v2+).
- Go: Versão 1.25+ (instalado localmente apenas se necessário para desenvolvimento fora do container; o Dockerfile usa golang:1.25-alpine).
- Make: Disponível no sistema (padrão em Linux/macOS; no Windows, use Git Bash ou WSL).
- Ferramentas de Dev: VS Code com extensões Docker e Go (opcional, mas recomendado para debug).

### Instalação e Configuração

1. Clone o repositório

```sh
git clone <url-do-repositorio>
```

2. Configure o ambiente de desenvolvimento

Copia os arquivos de configuração de ambiente para raiz do projeto.
Cria o direório tmp.

```sh
make setup-dev
```

### Rodando o Projeto

Use o Makefile para todos os comandos.

#### Modo de Desenvolvimento

Para rodar com hot-reload automático

```sh
make dev
```

Para rodar sem bloquear o terminal

```sh
make up
```

Para baixar os containers

```sh
make down
```

Para limpar arquivos do ambiente

```sh
make clean-dev
```

Para destroy containers e volumes

```sh
make destroy
```

Para executar a app em modo debug

```sh
make debug-app
```

Para executar debug a partir dos testes

```sh
make debug-test
```

Para executar os testes completos

```sh
make test
```

Para executar o lint correcoes de codigo

```sh
make lint
```

Para executar o format chegar regras de formatação

```sh
make format
```


Para resolver dependencias do projeto

```sh
make go-mod-tidy
```

Para verificar vulnerabilidades de segurança

```sh
make sec-check
```

Para modificações nos arquivos do Docker e precisar de um build

```sh
make build
```

Para visualizar logs dos containers

```sh
make logs
```

Para visualizar grafico de utilização de recursos da maquina

```sh
make stats
```

Quer saber mais execute o comando

```sh
make help
```

## Licença

MIT (veja LICENSE).

## Contribuição

1. Fork do repositório
2. `make up`
3. Faça as mudanças, commit.
4. Execute os testes
5. Envie a PR.

__Obrigado !__
