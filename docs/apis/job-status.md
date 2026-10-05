# GET /api/status

## Overview
The `/api/status` endpoint allows clients to retrieve the current status and detailed information of a specific job submitted for image processing. No authentication is required to access this endpoint.

## Request
| Attribute | Details |
|-----------|---------|
| **Method** | `GET` |
| **URL** | `/api/status` |
| **Query Parameters** | `jobid` (required) – the unique integer identifier of the job. |

### Example Request
```http
GET /api/status?jobid=123
```

## Response
The endpoint returns a JSON object describing the job. The structure varies based on the HTTP status code.

### Successful Response (200 OK)
```json
{
  "jobID": 123,
  "status": "processing",
  "progress": 45,
  "results": [
    {
      "imageID": "img_001",
      "storeID": 42,
      "perimeter": 1234.5,
      "processedAt": "2023-09-15T10:30:00Z"
    }
  ],
  "errors": [],
  "timestamps": {
    "createdAt": "2023-09-15T09:00:00Z",
    "startedAt": "2023-09-15T09:05:00Z",
    "updatedAt": "2023-09-15T10:30:00Z"
  }
}
```

#### Fields
- **jobID** – The unique integer identifier of the job.
- **status** – Current state of the job (e.g., `"pending"`, `"processing"`, `"completed"`, `"failed"`).
- **progress** – Completion percentage (0‑100) indicating how much of the job’s work has been performed.
- **results** – Array of objects, each representing the processing outcome for an individual image (contains image identifier, store identifier, computed perimeter, and timestamp).
- **errors** – Array of error messages encountered during processing; empty if no errors have occurred.
- **timestamps** – Object with `createdAt`, `startedAt`, and `updatedAt` fields, all in RFC3339 format.

### Error Responses
| HTTP Code | Description | Example Body |
|-----------|-------------|--------------|
| **400 Bad Request** | `jobid` parameter missing or malformed. | `{"error":"jobid query parameter is required"}` |
| **404 Not Found** | No job exists with the supplied `jobid`. | `{"error":"job not found"}` |
| **500 Internal Server Error** | Unexpected server error. | `{"error":"internal server error"}` |

#### Specific Error Cases
- **Missing `jobid`** – The query string does not contain a `jobid` parameter. The server responds with HTTP 400 and an appropriate error message.
- **Invalid `jobid` format** – The `jobid` value cannot be parsed as an integer (e.g., non‑numeric characters). The server responds with HTTP 400 and an error indicating the expected format.
- **Job not found** – A numeric `jobid` is provided but no corresponding job exists in the system. The server responds with HTTP 404.
- **Internal server error** – An unexpected condition prevents the request from being fulfilled (e.g., database or memory failure). The server responds with HTTP 500.

## Example Requests
```http
GET /api/status?jobid=123
GET /api/status?jobid=0      // edge case, job IDs start at 1
```

## Notes
- No authentication tokens or headers are required; the endpoint is publicly accessible.
- Job IDs are allocated sequentially from 1 upwards and are unique for the lifetime of the service.
- Job state is stored in memory; long‑running or completed jobs may be cleared based on application configuration.