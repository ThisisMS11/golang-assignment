package models

import (
	"time"
)

// Store represents a store from the store master
type Store struct {
	StoreID   string `json:"store_id"`
	StoreName string `json:"store_name"`
	AreaCode  string `json:"area_code"`
}

// SubmitRequest represents the incoming job request
type SubmitRequest struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

// Visit represents a store visit with images
type Visit struct {
	StoreID   string   `json:"store_id"`
	ImageURLs []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}

// JobResponse represents the response after job submission
type JobResponse struct {
	JobID int `json:"job_id"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// StatusResponse represents the job status response
type StatusResponse struct {
	Status string       `json:"status"`
	JobID  string       `json:"job_id"`
	Errors []StoreError `json:"error,omitempty"`
}

// StoreError represents a store-specific error
type StoreError struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

// ImageResult stores the result of image processing
type ImageResult struct {
	StoreID   string  `json:"store_id"`
	StoreName string  `json:"store_name"`
	AreaCode  string  `json:"area_code"`
	ImageURL  string  `json:"image_url"`
	Perimeter float64 `json:"perimeter"`
	VisitTime string  `json:"visit_time"`
}

// Job represents a processing job
type Job struct {
	ID        int
	Status    string
	Visits    []Visit
	Results   []ImageResult
	Errors    []StoreError
	StartTime time.Time
	EndTime   time.Time
}
