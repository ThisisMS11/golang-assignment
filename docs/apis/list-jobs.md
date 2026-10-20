# List Jobs API

## Overview
The **List Jobs** endpoint retrieves information about all jobs submitted to the Retail Pulse Image Processing Service. It returns a list of jobs with their identifiers, current statuses, and timestamps.

- **Endpoint**: `GET /api/jobs`
- **Description**: Returns an array of job objects representing all executed and ongoing jobs.
- **Authentication**: None required.

## Request
- **Method**: `GET`
- **URL**: `/api/jobs`
- **Query Parameters**: None
- **Headers**: None required

## Response

### Successful Response (200 OK)
The handler returns a JSON array of job objects. Each job object contains the following fields:

| Field | Type | Description |
|-------|------|-------------|
| `jobID` | integer | Unique identifier for the job (auto‑incremented). |
| `status` | string | Current status of the job. Possible values: `pending`, `processing`, `completed`, `failed`. |
| `createdAt` | string (ISO‑8601) | Timestamp when the job was submitted. |
| `updatedAt` | string (ISO‑8601) | Timestamp of the last status change. |

**Example**:
```json
[
  {
    "jobID": 1,
    "status": "completed",
    "createdAt": "2023-09-20T12:34:56Z",
    "updatedAt": "2023-09-20T12:36:02Z"
  },
  {
    "jobID": 2,
    "status": "processing",
    "createdAt": "2023-09-20T13:10:15Z",
    "updatedAt": "2023-09-20T13:10:20Z"
  }
]
```

> **Note**: The service may optionally include extra metadata (e.g., `storeCount`, `imageCount`) in each job object, but the core fields above are always present.

### Error Response (500 Internal Server Error)
If an unexpected error occurs while retrieving the job list, the service returns:

- **Status Code**: `500 Internal Server Error`
- **Content‑Type**: `application/json`

**Example**:
```json
{
  "error": "internal server error"
}
```

## Flow Diagram
The following diagram illustrates the flow when a client requests the list of jobs:

```mermaid
flowchart TD
    A[Client sends GET /api/jobs] --> B[HTTP Router]
    B --> C[Handler: listJobs]
    C --> D[Fetch all jobs from tracker]
    D --> E[Build job objects (ID, status, timestamps)]
    E --> F[Serialize to JSON]
    F --> G[Return HTTP 200 with JSON body]
    
    H[Error while fetching] --> I[Log error]
    I --> G[Return HTTP 500 with error JSON]
```

1. **Client** issues a `GET` request to `/api/jobs`.
2. The **Router** directs the request to the `listJobs` handler.
3. The handler queries the in‑memory job tracker for all jobs.
4. Each job is transformed into a JSON object containing `jobID`, `status`, `createdAt`, and `updatedAt`.
5. The array of objects is serialized and sent back with HTTP status `200 OK`.
6. Should any internal error occur (e.g., database or tracker failure), an HTTP `500` response with a generic error message is returned.

## Error Codes
| Code | Meaning |
|------|----------|
| `200` | Successfully returned the list of jobs. |
| `500` | Internal server error while processing the request. |

## See Also
- [Submit Job API](/docs/apis/submit-job.md) – create a new processing job.  
- [Job Status API](/docs/apis/job-status.md) – retrieve details of a specific job.