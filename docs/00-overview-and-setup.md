# Overview and Setup

## Project Purpose
The **Retail Pulse Image Processing Service** is a Go‑based microservice designed to ingest, process, and monitor large batches of retail‑store images. It automates the calculation of image perimeters while handling thousands of concurrent jobs, providing real‑time status tracking and robust error handling. The service simplifies the workflow for retailers who need to quickly analyze visual data from multiple locations.

## Core Problem Solved
Retailers collect massive volumes of images from store cameras or sensors. Manually downloading and measuring each image’s dimensions (height × width) to compute perimeters is impractical and error‑prone. This service:

* Accepts job definitions (store ID + list of image URLs) via a REST API.
* Downloads images concurrently and calculates perimeters (`2 × (Height + Width)`).
* Simulates GPU‑like processing with random delays.
* Validates store IDs against a master CSV file.
* Tracks job progress, successes, and failures.
* Exposes simple API endpoints for job submission, status checks, and audit.

The result is a scalable, thread‑safe platform that turns raw image URLs into measurable metrics with minimal human intervention.

## Key Features
| Feature | Description |
|---------|-------------|
| **Concurrent Job Processing** | Handles thousands of images per job using goroutines and worker pools. |
| **Thread‑Safe State Management** | Mutex‑protected job tracking to support simultaneous API calls. |
| **Store Master Integration** | Validates store IDs against `storeMaster.csv` and enriches results. |
| **REST API** | <ul><li>`POST /api/submit/` – creates a new job.</li><li>`GET /api/status?jobid=` – retrieves a specific job’s status.</li><li>`GET /api/jobs` – lists all jobs.</li></ul> |
| **Error Handling** | Marks a job as *failed* if any image cannot be processed or a store ID is missing. |
| **Image Format Support** | JPEG and PNG (easily extensible to other formats). |
| **Sequenced Job IDs** | Unique, monotonically increasing job identifiers. |
| **Docker‑Ready** | Includes a `Dockerfile` and `docker‑compose.yml` for containerized deployment. |
| **Health & Monitoring** | Simple status endpoint lets operators monitor progress and troubleshoot. |

## Tech Stack Summary
| Layer | Technology |
|-------|------------|
| **Runtime** | Go 1.16+ |
| **Concurrency** | Goroutines, `sync.Mutex`, worker pools |
| **Web Framework** | net/http (standard library) |
| **Data Storage** | In‑memory job store (CSV for static store master) |
| **File Handling** | `os`, `io`, `image/jpeg`, `image/png` packages |
| **Containerization** | Docker, Docker Compose |
| **Build Tools** | `go build`, `go run` |
| **Testing** | No dedicated test suite in the repo (focus on functional docs) |
| **Configuration** | Environment variable `PORT` (default 8080) |

## Local Setup Instructions

### 1️⃣ Clone the Repository
```bash
git clone https://github.com/ThisisMS11/golang-assignment.git
cd golang-assignment
```

### 2️⃣ Install Prerequisites
* **Go** – version 1.16 or newer  
  *Install*: follow the official Go installation guide for your OS.
* **Docker** & **Docker Compose** – optional, for containerized run (see Option 3).

### 3️⃣ Environment Variables (optional)
The service listens on the port defined by the `PORT` environment variable. If not set, it defaults to `8080`.

```bash
export PORT=8080   # or any desired port
```

### 4️⃣ Run with Go (Option 1)

1. **Ensure the store master file is present**  
   The CSV file `storeMaster.csv` must exist in the repository root (it does).

2. **Start the server**  
   ```bash
   go run ./cmd/server
   ```
   The binary will bind to `localhost:8080` (or the port you set via `PORT`). You should see log output indicating the HTTP server is ready.

3. **Verify the API**  
   ```bash
   # Submit a sample job (replace with real URLs and a valid store ID)
   curl -X POST http://localhost:8080/api/submit/ \
        -H "Content-Type: application/json" \
        -d '{"store_id":"123","image_urls":["https://example.com/image1.jpg","https://example.com/image2.png"]}'
   ```

   ```bash
   # Check job status
   curl "http://localhost:8080/api/status?jobid=1"
   ```

   ```bash
   # List all jobs
   curl http://localhost:8080/api/jobs
   ```

### 5️⃣ Run with Docker Compose (Option 2)

1. **Build the image** (first time only)  
   ```bash
   docker-compose build
   ```

2. **Start the service**  
   ```bash
   docker-compose up -d
   ```
   - The container maps host port `8080` to container port `8080`.
   - It mounts the local `storeMaster.csv` into the container at `/app/storeMaster.csv`.
   - The `PORT` environment variable is set to `8080` inside the container.

3. **Verify** – the same curl commands above work against `http://localhost:8080`.

### 6️⃣ Stop the Service
* **Go run**: press `Ctrl+C` in the terminal.
* **Docker Compose**:  
  ```bash
  docker-compose down
  ```

---

**Tip**: For development, you can run the service in one terminal (using `go run`) and test API calls in another terminal using `curl`. If you prefer a persistent environment, use Docker Compose; it also simplifies scaling and integration with orchestration tools.