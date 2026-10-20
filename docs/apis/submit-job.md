# Retail Pulse Image Processing Service - Submit Job API

## Endpoint Overview

| Property | Value |
|----------|-------|
| **URL** | `/api/submit/` |
| **Method** | `POST` |
| **Authentication** | None |
| **Description** | Accepts a set of image URLs and a store identifier to create a new batch processing job. The service will download the images, calculate their perimeters (2 × [height + width]), and track the job status until it completes or fails. |

---

## Request

### HTTP Request

```
POST /api/submit/
Content-Type: application/json
```

### Request Body Schema

```json
{
  "image_urls": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.png",
    // … more URLs …
  ],
  "store_id": "Store123"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `image_urls` | `[]string` | **Yes** | An ordered list of publicly accessible URLs for the images to be processed. Only JPEG and PNG formats are supported. |
| `store_id` | `string` | **Yes** | Identifier of the retail store that the images belong to. The value must exist in `storeMaster.csv`. |

### Example Request

```bash
curl -X POST http://localhost:8080/api/submit/ \
  -H "Content-Type: application/json" \
  -d '{
    "image_urls": [
      "https://example.com/store1/product1.jpg",
      "https://example.com/store1/product2.png"
    ],
    "store_id": "Store1"
  }'
```

---

## Response

### Successful Response

**HTTP Status:** `202 Accepted`

**Response Body**

```json
{
  "job_id": 1,
  "status": "created",
  "message": "Job created successfully"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `job_id` | `integer` | Unique sequential identifier for the job (starting at 1). |
| `status` | `string` | Initial status of the job – `"created"` indicates the job is queued for processing. |
| `message` | `string` | Human‑readable status description. |

### Error Responses

| HTTP Status | Example Body | Description |
|-------------|--------------|-------------|
| `400 Bad Request` | `{ "error": "invalid request payload", "details": "image_urls must be a non‑empty array" }` | The request body is malformed or missing required fields. |
| `400 Bad Request` | `{ "error": "validation failed", "details": "store_id 'XYZ' not found in store master" }` | The supplied `store_id` does not exist in `storeMaster.csv`. |
| `500 Internal Server Error** | `{ "error": "internal server error", "details": "failed to create job: <error details>" }` | An unexpected server‑side error occurred while creating the job. |

#### Example Error Payloads

```json
// Missing or empty image_urls
{
  "error": "validation failed",
  "details": "image_urls must be a non‑empty array"
}
```

```json
// Unknown store_id
{
  "error": "validation failed",
  "details": "store_id 'InvalidStore' not found in store master"
}
```

```json
// Generic server error
{
  "error": "internal server error",
  "details": "failed to create job: unable to acquire lock"
}
```

---

## Workflow Diagram

```mermaid
flowchart TD
    %% Client
    C[Client Application] -->|POST /api/submit/| S[API Gateway]
    
    %% API Layer
    S --> V[Validate Request]
    V -->|Valid| J[Create Job Record]
    V -->|Invalid| C[< Error Response>]
    
    %% Job Creation
    J -->|Generate ID<br>Store Job| T[Job Tracker (in‑memory)]
    T -->|Publish to Queue| P[Processor Pool]
    
    %% Processing
    P -->|Download Image| D[Image Downloader]
    D -->|Calculate Perimeter| C2[Perimeter Calculator]
    C2 -->|Store Result| R[Result Store]
    R -->|All Images Done| T -->|Update Job Status| J2[Job Completion]
    
    %% Response
    J2 -->|Update Status<br>Success/Failure| C[Client < Job Status>]
    
    %% Error Handling
    V -->|Store Not Found| C[< Error Response>]
    P -->|Download/Process Error| T -->|Mark Job Failed| C[Client < Error Response>]
    
    %% Styles (optional)
    classDef client fill:#e1f5fe,stroke:#01579b,color:#01579b;
    classDef api fill:#fff3e0,stroke:#e65100,color:#e65100;
    classDef proc fill:#f3e5f5,stroke:#4a148c,color:#4a148c;
    
    class C client;
    class S,V,J,T,P,D,C2,R,J2 api;
    class D,C2,R proc;
```

**Diagram Legend**

- **Client** – External service or UI that invokes the endpoint.  
- **API Gateway** – Router/handler that receives the request (`/api/submit/`).  
- **Validate Request** – Checks JSON structure, non‑empty `image_urls`, and that `store_id` exists in `storeMaster.csv`.  
- **Create Job Record** – Generates a new sequential `job_id` and stores an initial job entry in the in‑memory tracker.  
- **Job Tracker** – Central coordinator that holds job state (`created`, `processing`, `completed`, `failed`).  
- **Processor Pool** – Goroutine workers that download images, compute perimeters, and store results.  
- **Image Downloader** – Fetches the image from the supplied URL (fails if non‑HTTP/HTTPS, unsupported format, or network error).  
- **Perimeter Calculator** – Simulates GPU processing (random sleep) and computes `2 × (height + width)`.  
- **Result Store** – Persists the calculated perimeter for each image in the job’s result set.  
- **Job Completion** – Marks the job as `completed` (all images succeeded) or `failed` (any image error or store‑validation failure) and updates the tracker.  

The flow ensures that the client receives an immediate acknowledgment (`202 Accepted`) while processing occurs concurrently in the background. Errors are captured and reflected in the job status so the client can query the status endpoint later.

---

## Additional Notes

- The `storeMaster.csv` file is loaded at application start and used for validation. It contains a list of valid `store_id`s.
- Images are only processed if the format is JPEG or PNG. Unsupported formats cause the job to be marked as `failed`.
- Job IDs are monotonically increasing and never reused, even if a job fails.
- The service is designed for high concurrency; multiple clients can submit jobs simultaneously, each being assigned a unique `job_id`.

--- 

*Generated automatically from repository documentation.*