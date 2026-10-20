# POST /api/submit/

## Overview
Creates a new job for processing store images. The endpoint accepts a list of publicly accessible image URLs and a store identifier, enqueues the images for perimeter calculation, and returns a unique job ID that can be used to track the job’s progress.

## HTTP Method
**POST** `/api/submit/`

## Authentication
No authentication is required. The endpoint is publicly accessible.

## Request Body

The request must be sent as `application/json` with a JSON object containing the following fields:

| Field | Type | Description |
|-------|------|-------------|
| `image_urls` | `string[]` | An array of publicly accessible URLs pointing to JPEG or PNG images that will be downloaded and processed. |
| `store_id` | `integer` | The identifier of the retail store. This ID must exist in the `storeMaster.csv` file; otherwise the job is marked as failed. |

### Example Request

```json
{
  "image_urls": [
    "https://example.com/images/store12/photo1.jpg",
    "https://example.com/images/store12/photo2.png",
    "https://example.com/images/store12/photo3.jpg"
  ],
  "store_id": 12
}
```

```bash
curl -X POST http://localhost:8080/api/submit/ \
  -H "Content-Type: application/json" \
  -d '{
    "image_urls": [
      "https://example.com/images/store12/photo1.jpg",
      "https://example.com/images/store12/photo2.png",
      "https://example.com/images/store12/photo3.jpg"
    ],
    "store_id": 12
  }'
```

## Response Body

On successful submission, the service returns a JSON object with the following fields:

| Field | Type | Description |
|-------|------|-------------|
| `job_id` | `integer` | The unique sequential identifier for the created job (starting from 1 and incrementing for each new job). |
| `status` | `string` | The initial status of the job; always `"queued"` for a newly submitted job. |
| `message` | `string` (optional) | A human‑readable description of the submission result (e.g., `"Job submitted successfully"`). |

### Example Success Response

```json
{
  "job_id": 1,
  "status": "queued",
  "message": "Job submitted successfully"
}
```

## Error Responses

| HTTP Status | Description | Example Response Body |
|-------------|-------------|-----------------------|
| `400 Bad Request` | The request payload is malformed, missing required fields, or contains invalid data (e.g., `store_id` not found, unsupported image format, empty `image_urls`). | ```json { "error": "store_id 99 does not exist in store master data" } ``` |
| `415 Unsupported Media Type` | The `Content-Type` header is not set to `application/json`. | ```json { "error": "Unsupported Media Type" } ``` |
| `500 Internal Server Error` | An unexpected error occurred on the server while processing the request (e.g., database connection failure, panic). | ```json { "error": "Internal server error" } ``` |

### Example Error Response (400)

```json
{
  "error": "store_id 99 does not exist in store master data"
}
```

## Notes

- **Store Validation** – The `store_id` is checked against `storeMaster.csv`. If the ID is not present, the job is immediately marked as `failed`.
- **Image Format** – Only JPEG (`.jpg`, `.jpeg`) and PNG (`.png`) files are supported. Unsupported formats cause the job to be marked as `failed`.
- **Concurrency** – Multiple jobs can be submitted simultaneously; each receives its own `job_id` and is processed independently.
- **No Authentication** – The endpoint does not require any API keys or tokens. Ensure that image URLs are publicly accessible and compliant with the respective image hosts’ usage policies.