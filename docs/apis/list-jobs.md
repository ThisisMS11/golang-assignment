# GET /api/jobs - List All Jobs

## Overview
The `/api/jobs` endpoint allows clients to retrieve a list of all jobs that have been submitted to the Retail Pulse Image Processing Service. It returns an array of job objects, each representing a processing job along with its current status, associated store, image count, and timestamps.

## Request
| Property | Value |
|----------|-------|
| **Method** | `GET` |
| **URL** | `/api/jobs` |
| **Authentication** | Not required (public endpoint) |
| **Query Parameters** | None |

## Response
### HTTP Status
- `200 OK` – Successfully returned the list of jobs (may be an empty array).

### Response Body Schema
The response body is a JSON array containing one or more **Job Objects**.

#### Job Object Fields
- **jobID** (`integer`) – Unique identifier for the job, allocated sequentially starting from 1.
- **status** (`string`) – Current state of the job. Possible values:
  - `"pending"` – Job has been submitted but not yet started.
  - `"processing"` – Job is currently being processed.
  - `"completed"` – All images have been processed successfully.
  - `"failed"` – Job terminated due to an error (e.g., missing store, image download/processing failure).
- **storeID** (`integer`) – Identifier of the retail store associated with the job. Values are validated against the `storeMaster.csv` file.
- **imageCount** (`integer`) – Total number of images included in the job.
- **createdAt** (`string`, RFC3339) – Timestamp when the job was created (ISO‑8601 format, e.g., `"2023-09-15T10:30:00Z"`).
- **completedAt** (`string` or `null`, RFC3339) – Timestamp when the job finished (either successfully or with failure). `null` indicates the job is still in progress.
- **Additional fields** – The service may include extra information such as:
  - `errorMessage` (`string` or `null`) – Description of the failure, present only when `status` is `"failed"`.
  - `processedCount` (`integer`) – Number of images already processed (useful when `status` is `"processing"`).
  - `progress` (`number`) – Percentage progress of the job (0‑100) when `status` is `"processing"`.
  - `storeName` (`string`) – Human‑readable name of the store (enriched from the store master data).

### Example Request
```http
GET /api/jobs HTTP/1.1
Host: localhost:8080
```

### Example Response (200 OK)
```json
[
  {
    "jobID": 1,
    "status": "completed",
    "storeID": 101,
    "imageCount": 50,
    "createdAt": "2023-09-15T10:30:00Z",
    "completedAt": "2023-09-15T10:45:12Z",
    "storeName": "Store‑A"
  },
  {
    "jobID": 2,
    "status": "processing",
    "storeID": 102,
    "imageCount": 120,
    "createdAt": "2023-09-15T11:00:05Z",
    "completedAt": null,
    "processedCount": 94,
    "progress": 78
  },
  {
    "jobID": 3,
    "status": "failed",
    "storeID": 103,
    "imageCount": 30,
    "createdAt": "2023-09-15T11:15:20Z",
    "completedAt": "2023-09-15T11:22:00Z",
    "errorMessage": "Store ID 103 not found in store master data"
  }
]
```

## Error Responses
- **500 Internal Server Error** – Returned when an unexpected error occurs while retrieving jobs (e.g., failure reading the job registry, file I/O, or a panic in the handler). Example body:
  ```json
  {
    "error": "internal server error"
  }
  ```

  *Note:* The service does not return 4xx errors for this endpoint because no client‑supplied parameters are validated. Errors are solely server‑side.

## Authentication
No authentication is required. The endpoint is publicly accessible.

## See Also
- [Submit Job](/apis/submit-job.md) – Create a new processing job.
- [Get Job Status](/apis/status.md) – Retrieve details of a specific job.