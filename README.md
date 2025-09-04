# Zod URL Shortener

Zod is a high-performance and resilient URL shortening service built with Go, designed to handle high loads efficiently. This project leverages Docker for environment management and includes tools for development, testing, and debugging.

## Prerequisites

- **Docker** and **Docker Compose**: Ensure Docker and Docker Compose are installed on your system.
- **Make**: Required to run the Makefile commands.
- **Go**: Required if you want to run Go commands outside Docker (optional, as the project uses Docker for most tasks).

## Setup and Usage

The project uses a `Makefile` to streamline environment setup, service management, and development tasks. Below are the instructions for using the project based on the available `make` commands.

### 1. Set Up the Development Environment

To configure the local development environment, copy the necessary configuration files and set up the `.env` file with dynamic UID/GID:

```bash
make setup-dev
```

This command:
- Copies configuration files (e.g., `docker-compose.yml`, `Dockerfile`, etc.) from `_env/dev` to the project root.
- Creates a `.env` file based on `.env.example` with the current user's UID and GID.
- Creates a `tmp` directory for temporary files.

### 2. Start the Development Environment

To start the application in development mode with **Hot-Reload** using Air:

```bash
make dev
```

This command:
- Runs `setup-dev` to ensure the environment is configured.
- Starts the services defined in `docker-compose.yml` with automatic rebuilding of images.

### 3. Debug the Application

To debug the application using **Delve** (useful with VS Code's "Attach to Docker" feature):

```bash
make debug-app
```

This command:
- Runs `setup-dev` to configure the environment.
- Starts the debugging environment using `compose.debug.yaml`.

### 4. Debug Tests

To debug unit tests using **Delve** (also compatible with VS Code's "Attach to Docker"):

```bash
make debug-test
```

This command:
- Runs `setup-dev` to configure the environment.
- Starts the test debugging environment using `compose.debug-test.yaml`.

### 5. Run Services in Background

To start the services in the background (detached mode):

```bash
make up
```

This command:
- Runs `setup-dev` to configure the environment.
- Starts the services in the background with automatic rebuilding.

### 6. Stop and Remove Containers

To stop and remove containers without affecting configuration files:

```bash
make down
```

To stop containers without removing them:

```bash
make stop
```

### 7. Destroy Everything

To completely remove containers, volumes, and clean up the local environment:

```bash
make destroy
```

This command:
- Stops and removes containers and anonymous volumes for all environments (`docker-compose.yml`, `compose.debug.yaml`, `compose.debug-test.yaml`).
- Cleans up configuration files and the `tmp` directory.

### 8. Run Tests

To execute unit tests inside the container:

```bash
make test
```

This runs `go test -v ./...` within the `zod-api` service container.

### 9. Run Linter

To analyze the codebase with **golangci-lint**:

```bash
make lint
```

### 10. Organize Go Dependencies

To tidy up Go module dependencies:

```bash
make go-mod-tidy
```

This runs `go mod tidy` inside the container.

### 11. Check for Security Vulnerabilities

To analyze the codebase for security issues using **gosec**:

```bash
make sec-check
```

### 12. Format Code

To format the Go codebase using **gofumpt**:

```bash
make format
```

### 13. Build Docker Images

To rebuild Docker images without starting containers:

```bash
make build
```

This builds images for `docker-compose.yml`, `compose.debug.yaml`, and `compose.debug-test.yaml`.

### 14. View Logs

To view logs for running services:

```bash
make logs
```

### 15. List Running Containers

To list running containers:

```bash
make ps
```

### 16. Monitor Resource Usage

To monitor CPU and memory usage of containers in real-time:

```bash
make stats
```

### 17. Clean Up Development Environment

To remove configuration files and the `tmp` directory:

```bash
make clean-dev
```

### 18. View All Commands

To display the list of available `make` commands:

```bash
make help
```

## Additional Notes

- **Hot-Reload**: The `dev` command uses Air for live reloading during development.
- **Debugging**: Use the `debug-app` and `debug-test` commands with VS Code's debugging tools for an enhanced debugging experience.
- **Environment Variables**: The `.env` file is generated dynamically during `setup-dev`. Do not version control this file.
- **Docker Compose**: The project includes multiple Docker Compose configurations for development, debugging, and testing. Ensure the correct file is used for your workflow.

## License

This project is licensed under the terms specified in the `LICENSE` file.
