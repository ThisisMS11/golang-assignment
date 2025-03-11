# Retail Pulse Image Processing Service

A Go-based service that processes images collected from retail stores by downloading them and calculating their perimeters.

## Description

This service provides an API to process thousands of images collected from retail stores. It can perform multiple concurrent jobs, each containing thousands of images. The service:

1. Receives jobs with image URLs and store IDs
2. Downloads and processes images to calculate their perimeters (2 * [Height + Width])
3. Simulates GPU processing with random sleep times
4. Tracks job status and errors
5. Provides API endpoints to submit jobs,check status and get all jobs information

## Architecture

The service is built using:
- Thread-safe job tracking with mutexes to handle concurrent requests
- Preprocesses Store master data integration to validate store IDs and enrich results
- HTTP handlers for API endpoints
- Goroutines for concurrent image processing

## API Endpoints

1. **Submit Job**
   - URL: `/api/submit/`
   - Method: `POST`
   - Creates a new job for processing store images

2. **Get Job Status**
   - URL: `/api/status?jobid=123`
   - Method: `GET`
   - Returns the status and details of a specific job

2. **Get All Jobs Details**
   - URL: `/api/jobs`
   - Method: `GET`
   - Returns information about all the executed and ongoing jobs.

## Assumptions

1. The store master data is provided in a csv file (`storeMaster.csv`) in the root directory.
2. Image URLs are publicly accessible and can be downloaded without authentication.
4. A job is considered "failed" if any of its images fail to process or if a store ID doesn't exist.
5. The service handles only JPEG and PNG image formats (can be extended to support more).
6. Each job has a unique ID allocated in sequence from 1 upwards.

## Installation and Setup

### Prerequisites

- Go 1.16 or later
- Docker and Docker Compose (optional, for containerized deployment)

### Option 1: Run with Go

1. Clone the repository:
   ```bash
   git clone https://github.com/ThisisMS11/golang-assignment
   cd golang-assignment
   ```
2. Run the project:
   ```bash
   go run cmd/server/main.go
   ```
2. Build and run the service:
   ```bash
   go build -o server ./cmd/server
   ./server
   ```

3. The service will be available at `http://localhost:8080`

### Option 2: Run with Docker

1. Clone the repository:
   ```bash
   git clone https://github.com/ThisisMS11/golang-assignment
   cd golang-assignment
   ```

2. Build and run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

3. The service will be available at `http://localhost:8080`

## Testing

### Submit a job

```bash
curl -X POST http://localhost:8080/api/submit/ \
  -H "Content-Type: application/json" \
  -d '{
   "count": 2,
   "visits": [
      {
         "store_id": "RP00001",
         "image_url": [
            "https://www.gstatic.com/webp/gallery/2.jpg",
            "https://www.gstatic.com/webp/gallery/3.jpg"
         ],
         "visit_time": "2025-03-11T10:30:00Z"
      },
      {
         "store_id": "RP00002",
         "image_url": [
            "https://www.gstatic.com/webp/gallery/3.jpg"
         ],
         "visit_time": "2025-03-11T12:15:00Z"
      }
   ]
}'
```

### Check job status

```bash
curl http://localhost:8080/api/status?jobid=1
```

### Get all jobs

```bash
curl http://localhost:8080/api/jobs
```


## Future Improvements

Given more time, I would enhance this service with:

1. **Persistence**: Add a database (PostgreSQL or MongoDB) to store job data and results, ensuring data isn't lost on service restart.
2. **Authentication & Rate Limiting**: Implement an API key system and rate limits to secure and manage API usage.
3. **Metrics & Monitoring**: Add Prometheus metrics and logging with tools like Winston
4. **More Robust Error Handling**: Implement retries for image downloads.
5. **Unit & Integration Tests**: Add comprehensive test coverage for all components.
7. **API Documentation**: Generate API documentation using Swagger/OpenAPI.
8. **Health Checks**: Add health check endpoints for better integration with orchestration systems.
9. **Input Validation**: Add more thorough validation of input data.
