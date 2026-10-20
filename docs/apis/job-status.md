# Job Status API Endpoint

## Overview

The `/api/status` endpoint allows clients to retrieve the current status and details of a specific job. This endpoint is essential for monitoring the progress of image processing jobs submitted via the `/api/submit/` endpoint.

## HTTP Method

`GET`

## Request

### URL

```
GET /api/status?jobid=<job_id>
```

### Query Parameter

| Parameter | Type   | Required | Description                                   |
|-----------|--------|----------|-----------------------------------------------|
| `jobid`   | integer| Yes      | The unique identifier of the job to query.    |

**Example Request**

```http
GET /api/status?jobid=42 HTTP/1.1
Host: localhost:8080
```

### Authentication

No authentication is required. The endpoint is publicly accessible.

### Error Handling for Missing / Invalid `jobid`

| Scenario                              | HTTP Status | Response Body (JSON)                                    |
|---------------------------------------|-------------|--------------------------------------------------------|
| Missing `jobid` query parameter       | `400 Bad Request` | `{ "error": "missing jobid parameter" }`                |
| Invalid `jobid` format (non‑integer) | `400 Bad Request` | `{ "error": "invalid jobid format, must be an integer" }` |
| Job ID does not exist                 | `404 Not Found`   | `{ "error": "job not found" }`                           |
| Unexpected server error               | `500 Internal Server Error` | `{ "error": "internal server error" }`               |

The service validates the `jobid` parameter early in the request handling chain and returns appropriate HTTP status codes with JSON error messages.

## Response Schema

### Success Response (HTTP 200 OK)

```json
{
  "job_id": 42,
  "status": "processing",
  "progress": {
    "total_images": 150,
    "processed_images": 87,
    "percentage": 58
  },
  "results": [
    {
      "image_url": "https://example.com/store1/img1.jpg",
      "perimeter": 1240,
      "processed_at": "2023-10-05T14:32:10Z"
    },
    {
      "image_url": "https://example.com/store1/img2.png",
      "perimeter": 1864,
      "processed_at": "2023-10-05T14:33:05Z"
    }
  ],
  "errors": [
    {
      "image_url": "https://example.com/store1/img5.jpg",
      "error_message": "image not found",
      "occurred_at": "2023-10-05T14:40:22Z"
    }
  ],
  "timestamps": {
    "created_at": "2023-10-05T13:15:00Z",
    "started_at": "2023-10-05T13:16:05Z",
    "updated_at": "2023-10-05T14:45:12Z",
    "completed_at": null
  }
}
```

#### Response Fields

| Field          | Type                             | Description                                                                                    |
|----------------|-----------------------------------|------------------------------------------------------------------------------------------------|
| `job_id`       | integer                           | The unique identifier of the job.                                                             |
| `status`       | string (enum)                     | Possible values: `pending`, `processing`, `completed`, `failed`.                               |
| `progress`     | object                            | Contains counts and percentage of processed images.                                            |
| `results`      | array of result objects           | Each object holds the image URL, calculated perimeter, and timestamp of processing.            |
| `errors`       | array of error objects            | Each object captures the image URL that caused an error, error message, and occurrence time.    |
| `timestamps`   | object with ISO‑8601 strings      | Tracks lifecycle times: `created_at`, `started_at`, `updated_at`, `completed_at` (null if not finished). |

#### Progress Object

| Sub‑field           | Type    | Description                                              |
|---------------------|---------|----------------------------------------------------------|
| `total_images`      | integer | Total number of images in the job.                       |
| `processed_images`  | integer | Number of images that have completed processing.         |
| `percentage`        | integer | Percentage of completion (`(processed_images / total_images) * 100`). |

#### Result Object

| Sub‑field    | Type      | Description                               |
|--------------|-----------|-------------------------------------------|
| `image_url`  | string    | Original URL of the processed image.      |
| `perimeter`  | integer   | Calculated perimeter (`2 * (height + width)`).|
| `processed_at`| string   | ISO‑8601 timestamp when the image was processed.|

#### Error Object

| Sub‑field       | Type      | Description                                          |
|-----------------|-----------|------------------------------------------------------|
| `image_url`     | string    | URL of the problematic image.                       |
| `error_message` | string    | Human‑readable description of the error.             |
| `occurred_at`   | string    | ISO‑8601 timestamp when the error occurred.          |

#### Timestamps Object

| Sub‑field       | Type      | Description                                          |
|-----------------|-----------|------------------------------------------------------|
| `created_at`    | string    | When the job was submitted (ISO‑8601).               |
| `started_at`    | string    | When processing began (ISO‑8601, `null` if not started).|
| `updated_at`    | string    | Last modification time of the job (ISO‑8601).        |
| `completed_at`  | string or null | When processing finished (ISO‑8601) or `null` if still pending/failed.|

## Request/Response Example

### Request

```http
GET /api/status?jobid=7 HTTP/1.1
Host: example.com:8080
```

### Response (Job still processing)

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "job_id": 7,
  "status": "processing",
  "progress": {
    "total_images": 5,
    "processed_images": 2,
    "percentage": 40
  },
  "results": [
    {
      "image_url": "https://pics.example.com/storeA/img1.jpg",
      "perimeter": 3100,
      "processed_at": "2023-10-06T09:12:00Z"
    },
    {
      "image_url": "https://pics.example.com/storeA/img2.png",
      "perimeter": 2540,
      "processed_at": "2023-10-06T09:14:12Z"
    }
  ],
  "errors": [],
  "timestamps": {
    "created_at": "2023-10-06T08:58:00Z",
    "started_at": "2023-10-06T09:00:00Z",
    "updated_at": "2023-10-06T09:15:00Z",
    "completed_at": null
  }
}
```

### Response (Job completed)

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "job_id": 8,
  "status": "completed",
  "progress": {
    "total_images": 10,
    "processed_images": 10,
    "percentage": 100
  },
  "results": [ … ],
  "errors": [ … ],
  "timestamps": {
    "created_at": "2023-09-28T11:00:00Z",
    "started_at": "2023-09-28T11:02:00Z",
    "updated_at": "2023-09-28T11:15:30Z",
    "completed_at": "2023-09-28T11:15:30Z"
  }
}
```

## Flow Diagram

```mermaid
flowchart TD
    %% Client side
    C[Client Application] -->|Sends GET request with jobid| S[Server /api/status]
    
    %% Server side processing
    S -->|Validates jobid| V{Validate jobid}
    V -->|Missing/Invalid| V -->|Returns 400/404 with error JSON|
    
    V -->|Valid jobid| V -->|Retrieves job from in‑memory store|
    V -->|Job not found| V -->|Returns 404 with error JSON|
    
    V -->|Job found| V -->|Builds status response JSON|
    V -->|Returns 200 OK| S -->|C[Client receives response]|
    
    %% Annotations
    C -->|Monitors progress| C
    C -->|Polls periodically| C
```

**Explanation**

1. The client constructs a GET request to `/api/status` with the `jobid` query parameter.
2. The server validates the presence and format of `jobid`. If invalid, a 400 or 404 response is returned immediately.
3. For a valid `jobid`, the server fetches the job data from its in‑memory job store.
4. The server assembles the status payload, including job metadata, progress counters, results, errors, and timestamps.
5. The server responds with HTTP 200 and the JSON payload, which the client parses and displays (or uses for further actions).

## Additional Notes

* The job store is thread‑safe; concurrent status requests for different jobs will not interfere.
* The `progress.percentage` is calculated as integer division (`processed_images * 100 / total_images`).
* If a job has finished (either `completed` or `failed`), the `timestamps.completed_at` field will contain the final timestamp; otherwise it is `null`.
* The endpoint does not require any authentication, aligning with the service's design for open access to job status information.