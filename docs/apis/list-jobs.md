# List Jobs API

## Overview
This endpoint returns a list of all jobs (including ongoing, completed, and failed) currently tracked by the service. It provides a quick snapshot of the processing pipeline for monitoring and auditing purposes.

## Endpoint
- **URL**: `/api/jobs`
- **Method**: `GET`
- **Authentication**: None required

## Request
No parameters are required for this endpoint. Query parameters are not used.

**Sample request**:
```
GET /api/jobs HTTP/1.1
Host: localhost:8080
```

## Response

### Success Response
**Status Code**: `200 OK`  
**Content-Type**: `application/json`

**Schema**: An array of job objects. Each job object contains the following fields:
- `job_id` (integer): Unique identifier of the job, allocated sequentially starting from 1.
- `status` (string): Current processing state. Possible values: `pending`, `processing`, `completed`, `failed`.
- `created_at` (string, ISO 8601 timestamp): Timestamp when the job was created.
- `updated_at` (string, ISO 8601 timestamp): Timestamp of the last status change.
- `store_id` (integer): Identifier of the retail store associated with the job.
- `image_count` (integer, optional): Number of images the job processes (if known).
- `error_message` (string, optional): Error description if the job has failed (null otherwise).

**Example response**:
```json
[
  {
    "job_id": 1,
    "status": "completed",
    "created_at": "2023-09-01T12:34:56Z",
    "updated_at": "2023-09-01T12:45:00Z",
    "store_id": 42,
    "image_count": 150,
    "error_message": null
  },
  {
    "job_id": 2,
    "status": "processing",
    "created_at": "2023-09-01T13:10:00Z",
    "updated_at": "2023-09-01T13:15:30Z",
    "store_id": 7,
    "image_count": 200,
    "error_message": null
  }
]
```

### Error Responses
| Status Code | Description                                 | Example Body                     |
|-------------|---------------------------------------------|----------------------------------|
| `500`       | Internal server error – unable to retrieve jobs list | `{"error":"internal server error"}` |
| `503`       | Service unavailable (e.g., job store not initialized) | `{"error":"service unavailable"}` |

## Flow Diagram
The following Mermaid flowchart illustrates the flow of retrieving all jobs from the server to the client.

```mermaid
flowchart TD
    %% Client
    C[Client Application] -->|GET /api/jobs| S[Go Server]

    %% Server side
    S -->|Read jobs map (thread-safe)| M[In‑memory Job Storage]
    M -->|Iterate over all entries| A[Assemble Job Objects]
    A -->|Apply field formatting| F[JSON Marshal]
    F -->|Return response| C

    %% Optional error path
    S -.->|Panic / DB error| E[Error Handler]
    E -->|Return 5xx| C
```

**Explanation**:
1. The client sends a `GET` request to `/api/jobs`.
2. The server’s handler reads the internal `jobs` map (protected by a mutex) to ensure a consistent snapshot.
3. Each job entry is transformed into a structured object with the fields defined above.
4. The array of objects is serialized to JSON and returned with a `200 OK` status.
5. If an unexpected error occurs while accessing the job storage, the server returns a 5xx error response.

## Notes
- No authentication token is required.
- The list includes all jobs regardless of their current status (`pending`, `processing`, `completed`, `failed`).
- Job IDs are monotonic and unique across the lifetime of the process.
- Store master data is validated before a job is added; a job with an invalid `store_id` will be marked as `failed`.