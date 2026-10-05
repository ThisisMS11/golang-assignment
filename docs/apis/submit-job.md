# POST /api/submit/ Endpoint Documentation

## Overview

The **Submit Job** endpoint allows clients to create a new image processing job for a specific retail store. The service will download the provided image URLs, calculate the perimeter of each image (2 × [Height + Width]), simulate GPU processing with a random delay, and track the job’s progress. Authentication is **not** required to use this endpoint.

## Request

### Endpoint
```
POST /api/submit/
```

### Authentication
None required.

### Request Body (JSON)

The request payload must be a JSON object with the following fields (all required unless otherwise noted):

| Field | Type | Description |
|-------|------|-------------|
| `storeID` | `string` | Identifier of the retail store. Must exist in the `storeMaster.csv` file. |
| `imageURLs` | `[]string` | Array of publicly accessible URLs pointing to image files. Supported formats are JPEG and PNG (case‑insensitive). The service will attempt to download each URL; any failure marks the whole job as **failed**. |
| `maxConcurrency` *(optional)* | `int` | Maximum number of concurrent image‑processing goroutines for this job. Default value is `3` if omitted or ≤ 0. |

**Example Request**
```json
{
  "storeID": "STORE123",
  "imageURLs": [
    "https://example.com/images/product1.jpg",
    "https://example.com/images/product2.png"
  ],
  "maxConcurrency": 5
}
```

### Validation Rules
1. The JSON must be well‑formed.
2. `storeID` must be a non‑empty string that matches an entry in `storeMaster.csv`.
3. `imageURLs` must be an array containing at least one element.
4. Each URL must be a valid HTTP/HTTPS URL; malformed URLs result in a **400 Bad Request**.
5. Only JPEG and PNG extensions are accepted (case‑insensitive). Unsupported formats cause a **400 Bad Request**.

## Response

### Success Response (202 Accepted)

```
HTTP/1.1 202 Accepted
Content-Type: application/json
```

```json
{
  "jobID": 42,
  "status": "pending",
  "message": "Job created successfully and is being processed."
}
```

| Field | Type | Description |
|-------|------|-------------|
| `jobID` | `integer` | Unique sequential identifier allocated by the service (starting from 1). |
| `status` | `string` | Initial status of the job – always `"pending"` for a newly created job. |
| `message` | `string` | Human‑readable confirmation message. |

### Error Responses

| HTTP Code | Scenario | Response Body (JSON) |
|-----------|----------|----------------------|
| **400 Bad Request** | - Malformed JSON payload  <br> - Missing required fields (`storeID` or `imageURLs`)  <br> - `storeID` not found in store master  <br> - Empty `imageURLs` array  <br> - Invalid URL format  <br> - Unsupported image format | ```json { "error": "Descriptive error message" } ``` |
| **405 Method Not Allowed** | Request used a non‑POST method (e.g., GET) | ```json { "error": "Method not allowed" } ``` |
| **500 Internal Server Error** | Unexpected server error while processing the request (e.g., database/write‑to‑log failure) | ```json { "error": "Internal server error" } ``` |
| **503 Service Unavailable** | Service is temporarily down for maintenance or overload | ```json { "error": "Service unavailable" } ``` |

All error responses are returned in JSON format with a single `error` field describing the problem.

## Next Steps

- Use the returned `jobID` to poll the **Get Job Status** endpoint (`GET /api/status?jobid=<jobID>`) for updates.
- The job will transition through the statuses `"pending"` → `"processing"` → `"completed"`/`"failed"` as images are handled.
- Refer to the **Get All Jobs Details** endpoint (`GET /api/jobs`) for a global view of active and finished jobs.

## Notes

- The endpoint does **not** require any authentication token or API key.
- The service will immediately spawn a goroutine to orchestrate the job, allowing the client to continue using the API.
- Any failure during image download or processing (e.g., network error, unsupported format, malformed image) will cause the entire job to be marked as **failed**, and the failure details will be recorded for later inspection.