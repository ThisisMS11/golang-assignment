# Submit Job API

## Overview
The `/api/submit/` endpoint allows clients to submit a new batch of retail store images for processing. The service validates the store identifier, downloads the images from the provided URLs, calculates each image's perimeter (2 × (Height + Width)), and tracks the overall job status. No authentication is required for this endpoint.

## Request

### Endpoint
```http
POST /api/submit/
```

### Content‑Type
`application/json`

### Schema
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `store_id` | integer | **Yes** | Identifier of the retail store. Must exist in `storeMaster.csv`. |
| `image_urls` | array of strings | **Yes** | List of publicly accessible URLs pointing to JPEG or PNG images. The service only processes these two formats. |
| `job_id` | **read‑only** | No | Assigned by the service after the request is accepted. Clients must **not** include this field in the request body. |

**Example Request**
```json
{
  "store_id": 42,
  "image_urls": [
    "https://example.com/store42/img1.jpg",
    "https://example.com/store42/img2.png"
  ]
}
```

### Validation Rules
1. The JSON body must be well‑formed.
2. Both `store_id` and `image_urls` fields must be present.
3. `store_id` must map to a record in `storeMaster.csv` (see `store` package).
4. Each URL in `image_urls` must be a valid HTTP/HTTPS URL pointing to a JPEG (`.jpg`, `.jpeg`) or PNG (`.png`) file.
5. The array `image_urls` must not be empty.

If any of the above rules fail, the endpoint returns a `400 Bad Request` response (see **Error Codes** section).

## Response

### Success (202 Accepted)
```json
{
  "job_id": 1,
  "status": "submitted"
}
```

- **`job_id`** – Auto‑incremented unique identifier for the newly created job.
- **`status`** – Initial status of the job; always `"submitted"` on creation.

### Error Responses

| HTTP Code | Example Body | Description |
|-----------|--------------|-------------|
| `400` | `{"error":"missing field 'store_id'"}`, `{"error":"invalid store ID"}`, `{"error":"image_urls must not be empty"}`, `{"error":"unsupported image format"}` | Client‑side validation failure (missing/invalid fields, unrecognized store, unsupported file type, etc.). |
| `500` | `{"error":"internal server error"}` | Unexpected server condition (e.g., failure to start job tracker, database error, panic). |

All error bodies are JSON objects with a single `error` field describing the problem.

## Flow Diagram

```mermaid
flowchart TD
    A[Client POSTs to /api/submit/] --> B[Parse & Validate JSON]
    B --> C{Is store_id valid?}
    C -->|No| D[Return 400: invalid store ID]
    C -->|Yes| E[Assign next job_id]
    E --> F[Create Job record & store in memory]
    F --> G[Spawn goroutine for image processing]
    G --> H{For each image URL}
    H --> I[Download image]
    I --> J{Is format JPEG/PNG?}
    J -->|No| K[Mark image error, continue]
    J -->|Yes| L[Compute perimeter (2*(H+W))]
    L --> M[Simulate GPU processing (random sleep)]
    M --> N[Store result in job tracker]
    N --> O{All images processed?}
    O -->|No| H
    O -->|Yes| P[Update job status to "completed" / "failed"]
    P --> Q[Return 202 with job_id & status]
    D --> Q
```

**Explanation of the Diagram**
1. **Client → Parse** – The request payload is parsed and validated.
2. **Store validation** – Checks `storeMaster.csv` for the provided `store_id`.
3. **Job creation** – A unique `job_id` is generated and a new in‑memory job entry is created.
4. **Processing goroutine** – A background goroutine handles the heavy‑lifting to keep the HTTP handler non‑blocking.
5. **Image loop** – For each URL the service downloads, validates format, calculates perimeter, and simulates GPU processing.
6. **Result tracking** – Individual image results are recorded in the job’s internal tracker.
7. **Job finalization** – Once all images are processed, the overall job status is set ( `"completed"` if all images succeeded, `"failed"` otherwise).
8. **Response** – The client receives the assigned `job_id` and initial `"submitted"` status.

## Notes

- **Concurrency** – The service uses mutex‑protected maps to safely track multiple jobs across many goroutines.
- **Store master data** – The CSV file `storeMaster.csv` is loaded once at startup and used for all store validation checks.
- **Image handling** – Only JPEG and PNG files are processed. Unsupported formats cause the image to be marked as errored, but the job continues processing remaining images.
- **Job status lifecycle** – After submission a job is `"submitted"`. It may transition to `"processing"`, `"completed"`, or `"failed"` as processing progresses. The `/api/status` and `/api/jobs` endpoints expose these states to clients.