package models

import (
	"time"
)

type Store struct {
	StoreID   string `json:"store_id"`
	StoreName string `json:"store_name"`
	AreaCode  string `json:"area_code"`
}

type SubmitRequest struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

type Visit struct {
	StoreID   string   `json:"store_id"`
	ImageURLs []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}

type JobResponse struct {
	JobID int `json:"job_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type StatusResponse struct {
	Status string       `json:"status"`
	JobID  string       `json:"job_id"`
	Errors []StoreError `json:"error,omitempty"`
}

type StoreError struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

type ImageResult struct {
	StoreID   string  `json:"store_id"`
	StoreName string  `json:"store_name"`
	AreaCode  string  `json:"area_code"`
	ImageURL  string  `json:"image_url"`
	Perimeter float64 `json:"perimeter"`
	VisitTime string  `json:"visit_time"`
}

type Job struct {
	ID        int
	Status    string
	Visits    []Visit
	Results   []ImageResult
	Errors    []StoreError
	StartTime time.Time
	EndTime   time.Time
}
