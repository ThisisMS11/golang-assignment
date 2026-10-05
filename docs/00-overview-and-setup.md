# Retail Pulse Image Processing Service - Overview & Setup

## Project Purpose

The **Retail Pulse Image Processing Service** is a Go-based microservice designed to process large volumes of retail store images. It provides a RESTful API for submitting batch image processing jobs, tracking their progress, and retrieving results. The service calculates the perimeter of each downloaded image (2 × [Height + Width]) and simulates GPU-intensive processing workloads.

## Core Problem Solved

Retail organizations collect thousands of images from store locations for planogram compliance, shelf auditing, and inventory verification. This service addresses:

- **Scale**: Concurrent processing of multiple jobs, each containing thousands of images
- **Reliability**: Thread-safe job tracking with comprehensive error handling and status reporting
- **Data Enrichment**: Validates store IDs against a master dataset and enriches results with store metadata
- **Observability**: Real-time job status tracking (pending, processing, completed, failed) with detailed error reporting
- **Operational Simplicity**: Single binary deployment with minimal dependencies, container-ready

## Key Features

| Feature | Description |
|---------|-------------|
| **Job Submission** | `POST /api/submit` — Accepts image URLs and store IDs, returns unique job ID |
| **Status Polling** | `GET /api/status?jobid={id}` — Real-time progress, per-image results, and error details |
| **Job Listing** | `GET /api/jobs` — Aggregated view of all jobs with summary statistics |
| **Concurrent Processing** | Goroutine-based parallel image download and perimeter calculation |
| **Store Validation** | Pre-loads `storeMaster.csv` for O(1) store ID lookup and metadata enrichment |
| **Failure Isolation** | Individual image failures don't block job; job marked failed if any image fails |
| **Format Support** | JPEG and PNG (extensible via processor configuration) |
| **GPU Simulation** | Configurable random sleep to simulate compute-intensive inference workloads |

## Tech Stack Summary

| Layer | Technology |
|-------|------------|
| **Language** | Go 1.16+ |
| **HTTP Router** | Standard library `net/http` with custom router (`internal/api/router.go`) |
| **Concurrency** | Goroutines + `sync.Mutex` / `sync.RWMutex` for thread-safe job store |
| **Image Processing** | `image` (stdlib) + `github.com/disintegration/imaging` for decode/perimeter |
| **CSV Parsing** | `encoding/csv` (stdlib) for `storeMaster.csv` ingestion |
| **Containerization** | Docker (multi-stage build) + Docker Compose |
| **Deployment** | Single binary, binds to `PORT` env var (default 8080) |

**Project Structure Highlights**

```
├── cmd/server/main.go          # Entry point, server bootstrap
├── internal/
│   ├── api/handlers.go         # HTTP handlers (submit, status, jobs)
│   ├── api/router.go           # Route registration
│   ├── config/config.go        # Configuration (env-driven)
│   ├── job/job.go              # Job struct, status constants, thread-safe store
│   ├── models/models.go        # Request/response DTOs
│   ├── processor/processor.go  # Image download, decode, perimeter calc, GPU sim
│   └── store/store.go          # StoreMaster CSV loading, store lookup
├── storeMaster.csv             # Store metadata (ID, name, location, etc.)
├── Dockerfile                  # Multi-stage build
└── docker-compose.yml          # Local dev / prod compose
```

---

## Local Setup Instructions

### Prerequisites

- **Go** 1.16 or later (`go version`)
- **Docker** 20.10+ and **Docker Compose** v2 (optional, for containerized run)
- **Git** for cloning
- Network access to download images from provided URLs (public HTTP/HTTPS)

---

### Option 1: Run Natively with Go

#### 1. Clone the Repository
```bash
git clone https://github.com/ThisisMS11/golang-assignment.git
cd golang-assignment
```

#### 2. Install Dependencies
```bash
go mod tidy
```
> The project uses Go modules. `go mod tidy` downloads `github.com/disintegration/imaging` and any transitive dependencies.

#### 3. Verify Store Master Data
Ensure `storeMaster.csv` exists in the project root (committed in repo). The service loads it at startup.

```bash
head -n 5 storeMaster.csv
# Expected columns: store_id,store_name,city,state,country,...
```

#### 4. Configure Environment Variables (Optional)

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server listen port |
| `STORE_MASTER_PATH` | `./storeMaster.csv` | Path to store master CSV |
| `MAX_CONCURRENT_IMAGES` | `50` | Max parallel image downloads per job |
| `GPU_SIM_MIN_MS` | `100` | Min simulated GPU processing time (ms) |
| `GPU_SIM_MAX_MS` | `500` | Max simulated GPU processing time (ms) |

Create a `.env` file or export in shell:
```bash
export PORT=8080
export STORE_MASTER_PATH=./storeMaster.csv
export MAX_CONCURRENT_IMAGES=50
export GPU_SIM_MIN_MS=100
export GPU_SIM_MAX_MS=500
```

#### 5. Run the Server
```bash
go run cmd/server/main.go
```
Output:
```
2024/01/15 10:30:45 Loading store master data from ./storeMaster.csv
2024/01/15 10:30:45 Loaded 1250 stores
2024/01/15 10:30:45 Server starting on :8080
```

#### 6. Verify Health
```bash
curl -s http://localhost:8080/api/jobs | jq .
# Expected: {"jobs":[]} (empty list on fresh start)
```

---

### Option 2: Run with Docker Compose (Recommended for Parity)

#### 1. Clone & Navigate
```bash
git clone https://github.com/ThisisMS11/golang-assignment.git
cd golang-assignment
```

#### 2. Build and Start
```bash
docker compose up --build -d
```
- Builds the multi-stage Dockerfile (compiles binary in builder stage, copies to `scratch`/`distroless` runtime)
- Mounts `storeMaster.csv` as a volume for live edits
- Maps container port 8080 → host port 8080
- Restarts on failure

#### 3. View Logs
```bash
docker compose logs -f kiranaclub
```

#### 4. Verify
```bash
curl -s http://localhost:8080/api/jobs | jq .
```

#### 5. Stop & Cleanup
```bash
docker compose down        # Stop containers
docker compose down -v     # Stop + remove volumes (if any)
```

---

### Option 3: Build & Run Docker Image Manually

```bash
# Build
docker build -t retail-pulse:local .

# Run (mount CSV, set port)
docker run -d \
  --name retail-pulse \
  -p 8080:8080 \
  -v "$(pwd)/storeMaster.csv:/app/storeMaster.csv" \
  -e PORT=8080 \
  retail-pulse:local
```

---

## Quick API Smoke Test

```bash
# 1. Submit a job
JOB_RESPONSE=$(curl -s -X POST http://localhost:8080/api/submit \
  -H "Content-Type: application/json" \
  -d '{
    "images": [
      {"url": "https://example.com/image1.jpg", "store_id": "101"},
      {"url": "https://example.com/image2.png", "store_id": "102"}
    ]
  }')
echo "$JOB_RESPONSE" | jq .

# Extract job_id
JOB_ID=$(echo "$JOB_RESPONSE" | jq -r .job_id)

# 2. Poll status until completed/failed
while true; do
  STATUS=$(curl -s "http://localhost:8080/api/status?jobid=$JOB_ID" | jq -r .status)
  echo "Job $JOB_ID status: $STATUS"
  [[ "$STATUS" == "completed" || "$STATUS" == "failed" ]] && break
  sleep 2
done

# 3. Final details
curl -s "http://localhost:8080/api/status?jobid=$JOB_ID" | jq .

# 4. List all jobs
curl -s http://localhost:8080/api/jobs | jq .
```

---

## Troubleshooting

| Symptom | Likely Cause | Resolution |
|---------|--------------|------------|
| `storeMaster.csv not found` | Wrong working directory / volume mount missing | Run from project root; verify `STORE_MASTER_PATH` |
| `bind: address already in use` | Port 8080 occupied | Change `PORT` env var or stop conflicting process |
| Images stuck in `downloading` | Network egress blocked / URL unreachable | Ensure container/host can reach image URLs; test with `curl <url>` |
| Job fails immediately | Invalid store_id in request | Verify store_id exists in `storeMaster.csv` |
| High memory under load | Unbounded concurrent downloads | Tune `MAX_CONCURRENT_IMAGES` lower |

---

## Next Steps

- See **[API Reference](../01-api-reference.md)** for detailed request/response schemas
- See **[Architecture & Concurrency Model](../02-architecture.md)** for internal design deep-dive
- See **[Deployment Guide](../03-deployment.md)** for production hardening (TLS, reverse proxy, health checks, observability)