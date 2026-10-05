# Infrastructure and Deployment

## Overview
This section describes the containerized deployment of the Retail Pulse Image Processing Service using Docker and Docker Compose. It covers the service definitions, required environment variables, port mappings, volume mounts, and the commands to start and tear down the stack.

## Prerequisites
- **Docker** (version 20.10 or later)
- **Docker Compose** (included with Docker or installed separately)
- Sufficient disk space to build the Go application and store the generated image processing data.

## Docker Compose Configuration

The `docker-compose.yml` file defines a single service, `kiranaclub`, which runs the Go application inside a Docker container.

```yaml
version: '3'

services:
  kiranaclub:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./storeMaster.csv:/app/storeMaster.csv
    environment:
      - PORT=8080
    restart: on-failure
```

### Service Details

| Component | Description |
|-----------|-------------|
| **Service Name** | `kiranaclub` |
| **Build Context** | Current directory (`.`) – uses the repository’s `Dockerfile` |
| **Port Mapping** | Host `8080` → Container `8080` (exposes the HTTP API) |
| **Volume Mount** | `storeMaster.csv` from the host is mounted at `/app/storeMaster.csv` inside the container (required for store‑ID validation) |
| **Environment Variables** | `PORT=8080` – sets the listening port for the Go application |
| **Restart Policy** | `on-failure` – Docker will restart the container automatically if the process exits with a non‑zero status |

## Required Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | The TCP port on which the Go service listens. Must match the port exposed in `docker-compose.yml`. |

> **Note:** The Go application reads the `PORT` environment variable (see `internal/config/config.go`) to determine the binding address. No other variables are required for a basic deployment.

## Docker Image Build

The service is built using the `Dockerfile` located in the repository root. It typically performs the following steps:

1. Sets up a Go workspace.
2. Copies go source files (`./...`).
3. Downloads Go modules (`go mod download`).
4. Compiles the binary (`go build -o app ./cmd/server`).
5. Creates a minimal runtime image (e.g., using a `distroless` base) that contains only the compiled binary and the `storeMaster.csv` file.

You can inspect the Dockerfile to confirm the exact steps, but they are not required for everyday operations—just for rebuilding when dependencies change.

## Starting the Stack

1. **Navigate to the repository root** (where `docker-compose.yml` resides).
2. **Build and start all services** (first time or after changes):
   ```bash
   docker compose up --build -d
   ```
   - `--build` forces a rebuild of the image.
   - `-d` runs the containers in detached mode.

3. **Verify the service is running**:
   ```bash
   docker compose ps
   ```
   You should see the `kiranaclub` container listed as `up`.

4. **Check logs** (optional, for troubleshooting):
   ```bash
   docker compose logs -f kiranaclub
   ```

## Accessing the Service

With the container running, the API is accessible from any host that can reach the Docker host:

```
http://localhost:8080/api/submit/
http://localhost:8080/api/status?jobid=123
http://localhost:8080/api/jobs
```

If you are running Docker inside a VM or remote host, replace `localhost` with the appropriate IP or hostname.

## Stopping and Tearing Down the Stack

When you no longer need the deployment:

1. **Stop and remove containers, networks, and volumes**:
   ```bash
   docker compose down --volumes
   ```
   - `down` stops the containers and removes the automatically created networks.
   - `--volumes` also removes any named volumes (in this setup only the `storeMaster.csv` mount is used, which is a bind mount and thus removed from the host perspective).

2. **Optional: Remove the built image** (not required for a fresh clone):
   ```bash
   docker image rm kiranaclub
   ```

## Rebuilding After Dependencies Change

If the `Dockerfile`, `go.mod`, `go.sum`, or any source files change, you can rebuild without recreating the entire stack:

```bash
docker compose build --no-cache
docker compose up -d
```

## Common Troubleshooting Steps

| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| `docker compose up` fails with “executable file not found” | Missing Docker installed or wrong architecture | Install Docker for your OS; ensure the host architecture matches the container image (e.g., `amd64`). |
| Container exits with status `1` and no logs | Application error (e.g., missing `storeMaster.csv`) | Verify the CSV file is present in the host directory and correctly formatted. The volume mount may need to be changed to an absolute path. |
| Port `8080` already in use | Another process (e.g., a manually run Go binary) is listening | Stop the conflicting process or change the `PORT` environment variable and the corresponding `docker-compose.yml` port mapping. |
| “permission denied” when mounting CSV | Incorrect file permissions | Ensure the CSV file is readable by the user running Docker (often root inside the container). Use `chmod 644 storeMaster.csv` if needed. |

## Summary

- The deployment consists of a single Docker service defined in `docker-compose.yml`.
- It builds from the repository’s `Dockerfile`, maps port `8080`, mounts the required `storeMaster.csv`, and sets the `PORT` environment variable.
- Use `docker compose up --build -d` to start and `docker compose down --volumes` to stop the stack.
- The service is accessible via `http://localhost:8080` and provides the three API endpoints for job submission, status checking, and job listing.