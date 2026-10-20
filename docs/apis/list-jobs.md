# GET /api/jobs

## Overview
The `/api/jobs` endpoint returns a list of all jobs that have been submitted to the service, including their current status, timestamps, and processing statistics. It is useful for monitoring ongoing work and reviewing completed or failed jobs.

## Request

**Method:** `GET`  
**URL:** `/api/jobs`  
**Headers:** None required (no authentication)  
**Query Parameters:** None required. The endpoint returns all jobs; filtering is not currently supported.

### Example Request
```bash
curl -X GET http://localhost:8080/api/jobs
```

## Response

**HTTP Status:** `200 OK`

The response body is a JSON array where each element represents a job object. Jobs are returned in chronological order (newest first).

### Response Schema

Each job object contains the following fields:

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Unique sequential identifier for the job (starting from 1). |
| `status` | string | Current processing status of the job. Possible values: `queued`, `in-progress`, `completed`, `failed`. |
| `createdAt` | string (ISO‑8601 timestamp) | Timestamp when the job was submitted (e.g., `2023-10-01T12:00:00Z`). |
| `updatedAt` | string (ISO‑8601 timestamp) | Timestamp of the most recent status change for the job. |
| `totalImages` | integer | Total number of images associated with the job. |
| `processedImages` | integer | Number of images that have been successfully processed so far. |
| `errorCount` | integer | Number of images that failed processing (or store validation) for this job. |
| `storeIDs` | array of integers | List of store identifiers referenced by the images in this job (validated against `storeMaster.csv`). |
| `imageURLs` *(optional)* | array of strings | The original image URLs submitted with the job (may be omitted for privacy or size reasons). |

### Example Response
```json
[
  {
    "id": 1,
    "status": "completed",
    "createdAt": "2023-10-01T12:00:00Z",
    "updatedAt": "2023-10-01T12:05:30Z",
    "totalImages": 150,
    "processedImages": 150,
    "errorCount": 0,
    "storeIDs": [101, 102]
  },
  {
    "id": 2,
    "status": "in-progress",
    "createdAt": "2023-10-01T12:01:00Z",
    "updatedAt": "2023-10-01T12:03:45Z",
    "totalImages": 200,
    "processedImages": 120,
    "errorCount": 2,
    "storeIDs": [103]
  }
]
```

## Error Responses

| HTTP Code | Description | Body |
|-----------|-------------|------|
| `500 Internal Server Error` | An unexpected error occurred while retrieving the job list (e.g., internal state corruption). | `{"error":"internal server error"}` |
| `503 Service Unavailable` | The service is temporarily unable to process the request (e.g., store master data missing). | `{"error":"service unavailable"}` |

All other status codes follow standard HTTP semantics (e.g., `404 Not Found` is not applicable because the endpoint itself always exists).

## Authentication

No authentication is required for this endpoint. All requests are treated as public.

## Additional Notes

- The endpoint does not support pagination; large numbers of jobs are retained in memory. For production use, consider implementing pagination or a separate administrative interface.
- Job status transitions follow the workflow defined in the job processor: `queued` → `in-progress` → `completed` or `failed`.
- Failed jobs are retained with their error count so that you can audit which images caused processing failures.