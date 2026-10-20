# GET /api/status

**Endpoint:** `/api/status`

**Method:** `GET`

## Query Parameters
| Parameter | Required | Description |
|-----------|----------|-------------|
| `jobid`   | Yes      | The unique identifier of the job whose status you want to retrieve. Must be a positive integer. |

## Request Details
- **Headers:** No authentication is required. Any standard HTTP headers are ignored.
- **URL Example:** `GET /api/status?jobid=42`

## Authentication
None – the endpoint is publicly accessible.

## Response Codes
- `200 OK` – Job status retrieved successfully. Returns a JSON body (see below).
- `400 Bad Request` – `jobid` is missing, malformed, or not a valid integer.
- `404 Not Found` – No job exists with the supplied `jobid`.
- `500 Internal Server Error` – An unexpected server error occurred while processing the request.

## Response Body (JSON)

```json
{
  "status": "pending|processing|completed|failed",
  "job_details": {
    "id": 42,
    "store_id": "ST001",
    "images": [
      {
        "url": "https://example.com/image1.jpg",
        "height": 1080,
        "width": 1920,
        "perimeter": 6000
      },
      // … additional image objects …
    ],
    "progress": 12,
    "total": 20,
    "created_at": "2023-09-15T10:30:00Z",
    "updated_at": "2023-09-15T10:45:23Z"
  },
  "errors": [
    "Failed to download image https://example.com/image3.png",
    "Store ID ST002 not found in master data"
  ]
}
```

### Field Descriptions
- **`status`** – Current processing state of the job. Possible values:
  - `pending`   – Job created but not yet started.
  - `processing` – Images are being downloaded/processed.
  - `completed` – All images processed successfully.
  - `failed`    – Job terminated due to errors.
- **`job_details`** – Object containing full job information.
  - **`id`** – Unique job identifier (sequential integer).
  - **`store_id`** – Identifier of the store associated with the job (validated against `storeMaster.csv`).
  - **`images`** – Array of image objects being processed.
    - **`url`** – Original image URL.
    - **`height`** – Image height in pixels.
    - **`width`** – Image width in pixels.
    - **`perimeter`** – Calculated perimeter (2 × [height + width]).
  - **`progress`** – Number of images already processed.
  - **`total`** – Total number of images in the job.
  - **`created_at`** – ISO‑8601 timestamp when the job was created.
  - **`updated_at`** – ISO‑8601 timestamp of the last status change.
- **`errors`** – List of error messages encountered during processing. Empty array (`[]`) when no errors have occurred.

## Example Success Response (Status 200)

```json
{
  "status": "processing",
  "job_details": {
    "id": 7,
    "store_id": "ST005",
    "images": [
      {
        "url": "https://example.com/store5/001.jpg",
        "height": 800,
        "width": 1200,
        "perimeter": 4000
      },
      {
        "url": "https://example.com/store5/002.png",
        "height": 1024,
        "width": 768,
        "perimeter": 3584
      }
    ],
    "progress": 1,
    "total": 5,
    "created_at": "2023-10-01T08:15:42Z",
    "updated_at": "2023-10-01T08:16:10Z"
  },
  "errors": []
}
```

## Example Error Responses

- **Missing `jobid` (400 Bad Request)**
  ```json
  { "error": "missing required query parameter 'jobid'" }
  ```

- **Job not found (404 Not Found)**
  ```json
  { "error": "job with id 99 not found" }
  ```

- **Server error (500 Internal Server Error)**
  ```json
  { "error": "internal server error" }
  ```

## Notes
- The service maintains thread‑safe job state using mutexes, so concurrent requests are handled correctly.
- No authentication or authorization is required for this endpoint.
- The `jobid` values are allocated sequentially from 1 upward; invalid or out‑of‑range identifiers will result in a 404.
- If a job fails, the `errors` array will contain detailed messages describing why processing could not complete (e.g., download failures, missing store IDs, unsupported image formats).

## See Also
- [Submit Job](/apis/submit-job.md) – Creates a new processing job.
- [List All Jobs](/apis/list-jobs.md) – Retrieves information about all jobs.