# Infrastructure and Deployment

This document describes the containerized infrastructure for the Retail Pulse Image Processing Service, including the Docker build process, Docker Compose service definitions, environment configuration, networking, storage, and operational procedures for starting and stopping the stack.

---

## 1. Container Image Build

### 1.1 Dockerfile Overview

The repository includes a multi-stage `Dockerfile` that produces a minimal runtime image.

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/storeMaster.csv .
EXPOSE 8080
ENTRYPOINT ["/app/server"]
```

**Key characteristics:**
- **Base image**: `golang:1.21-alpine` for compilation, `alpine:3.18` for runtime (~5 MB).
- **Static binary**: `CGO_ENABLED=0` ensures a statically linked executable with no external libc dependencies.
- **Artifacts copied**: Compiled binary `server` and the static `storeMaster.csv` reference file.
- **Entrypoint**: Direct execution of the binary; no shell wrapper.
- **Exposed port**: `8080` (documented via `EXPOSE`; actual mapping is controlled by Compose).

### 1.2 Building the Image Manually

```bash
# From repository root
docker build -t retail-pulse:local .
```

---

## 2. Docker Compose Service Definitions

The `docker-compose.yml` file defines a single service named **kiranaclub**.

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

### 2.1 Service: `kiranaclub`

| Attribute | Value | Description |
|-----------|-------|-------------|
| **Build context** | `.` (repository root) | Uses the local `Dockerfile`. |
| **Container name** | `kiranaclub` (implicit) | Compose assigns `project_kiranaclub_1` unless overridden. |
| **Restart policy** | `on-failure` | Restarts only on non-zero exit codes; avoids restart loops on configuration errors. |
| **Network mode** | `bridge` (default) | Isolated bridge network; service reachable via `localhost:8080` on host. |

---

## 3. Environment Variables

| Variable | Default | Required | Purpose |
|----------|---------|----------|---------|
| `PORT` | `8080` | Yes | TCP port the HTTP server binds to inside the container. Must match the Compose port mapping target. |

> **Note**: The application reads `PORT` via `internal/config/config.go`. No other environment variables are currently consumed. Secrets (e.g., database credentials, API keys) are not used in this version.

---

## 4. Port Mappings

| Host Port | Container Port | Protocol | Service |
|-----------|----------------|----------|---------|
| `8080` | `8080` | TCP | `kiranaclub` HTTP API |

- Access the API at `http://localhost:8080` from the Docker host.
- For remote access, bind to `0.0.0.0:8080` (default) or restrict to `127.0.0.1:8080` by changing the mapping to `"127.0.0.1:8080:8080"`.

---

## 5. Volume Mounts

| Host Path | Container Path | Mode | Purpose |
|-----------|----------------|------|---------|
| `./storeMaster.csv` | `/app/storeMaster.csv` | Read-only (implicit) | Provides the store master CSV to the running container. The file is also baked into the image; the mount allows live updates without rebuilding. |

**Behavior:**
- The mount is a **bind mount**; changes on the host are immediately visible inside the container.
- If the file is missing on the host, Compose will create an empty directory at the target path, causing the application to fail at startup. Ensure the file exists before `docker compose up`.

---

## 6. Starting the Stack

### 6.1 Prerequisites
- Docker Engine ≥ 20.10
- Docker Compose v2 (`docker compose`) or v1 (`docker-compose`)
- `storeMaster.csv` present in the repository root

### 6.2 Commands

```bash
# Build and start in foreground (logs streamed)
docker compose up --build

# Build and start in background (detached)
docker compose up --build -d

# Start without rebuilding (use existing images)
docker compose up -d
```

### 6.3 Verification

```bash
# Check container status
docker compose ps

# Tail logs
docker compose logs -f kiranaclub

# Health check (once server is up)
curl -s http://localhost:8080/api/jobs | jq .
```

Expected healthy output: `[]` (empty job list) with HTTP 200.

---

## 7. Tearing Down the Stack

### 7.1 Stop and Remove Containers (Preserve Volumes/Networks)

```bash
docker compose down
```
- Stops and removes the `kiranaclub` container.
- Removes the default project network.
- **Does not** remove the built image or the bind-mounted `storeMaster.csv`.

### 7.2 Full Cleanup (Includes Images and Orphan Volumes)

```bash
docker compose down --rmi all --volumes --remove-orphans
```
- `--rmi all`: Deletes the `retail-pulse` image built by Compose.
- `--volumes`: Removes any named volumes (none defined here, but safe to include).
- `--remove-orphans`: Cleans up containers from previous project runs.

---

## 8. Operational Notes

### 8.1 Updating `storeMaster.csv`
Because the CSV is bind-mounted, replace the file on the host and **restart the container** to reload:

```bash
cp new_storeMaster.csv ./storeMaster.csv
docker compose restart kiranaclub
```

### 8.2 Log Rotation
Docker’s default JSON log driver writes to `/var/lib/docker/containers/<id>/<id>-json.log`. For production, configure a log driver (e.g., `journald`, `fluentd`, `loki`) in `docker-compose.yml`:

```yaml
services:
  kiranaclub:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "5"
```

### 8.3 Resource Limits (Optional)
Add resource constraints to prevent noisy-neighbor issues:

```yaml
services:
  kiranaclub:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 128M
```

---

## 9. Troubleshooting Quick Reference

| Symptom | Likely Cause | Resolution |
|---------|--------------|------------|
| `port 8080 already in use` | Another process on host port 8080 | Change host mapping: `"8081:8080"` or stop conflicting process. |
| `storeMaster.csv: no such file or directory` | File missing in repo root | Restore `storeMaster.csv` from version control or backup. |
| Container exits immediately | Binary fails to start (e.g., config error) | Run `docker compose logs kiranaclub` to see panic / error output. |
| `permission denied` on CSV | File permissions on host prevent read | `chmod 644 storeMaster.csv`. |

---

## 10. Summary Checklist for Deployment

- [ ] `storeMaster.csv` present in repository root.
- [ ] `docker compose up --build -d` completes without errors.
- [ ] `docker compose ps` shows `kiranaclub` status `Up`.
- [ ] `curl http://localhost:8080/api/jobs` returns `[]` (HTTP 200).
- [ ] Log output shows `Server starting on port 8080` (or configured `PORT`).