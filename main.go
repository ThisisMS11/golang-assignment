package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
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

var (
	jobs        = make(map[int]*Job)
	jobsMutex   sync.RWMutex
	nextJobID   = 1
	jobIDMutex  sync.Mutex
	storeMaster = make(map[string]Store)
)

// Load store data from a JSON file
func loadStoreMaster() error {
	// Create store master file if it doesn't exist
	if _, err := os.Stat("store_master.json"); os.IsNotExist(err) {
		// Sample data
		sampleStores := []Store{
			{StoreID: "S00339218", StoreName: "Walmart Central", AreaCode: "A001"},
			{StoreID: "S01408764", StoreName: "Target Downtown", AreaCode: "A002"},
			{StoreID: "S02547931", StoreName: "Kroger Heights", AreaCode: "A003"},
		}

		data, err := json.MarshalIndent(sampleStores, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile("store_master.json", data, 0644); err != nil {
			return err
		}
	}

	// Read and parse the file
	data, err := os.ReadFile("store_master.json")
	if err != nil {
		return err
	}

	var stores []Store
	if err := json.Unmarshal(data, &stores); err != nil {
		return err
	}

	// Populate the map
	for _, store := range stores {
		storeMaster[store.StoreID] = store
	}

	return nil
}

// Download and process an image
func processImage(imageURL, storeID, visitTime string) (float64, error) {
	// Parse URL
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return 0, fmt.Errorf("invalid URL: %v", err)
	}

	// Make HTTP request
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return 0, fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to download image: HTTP %d", resp.StatusCode)
	}

	// Decode image
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to decode image: %v", err)
	}

	// Calculate perimeter
	bounds := img.Bounds()
	height := bounds.Max.Y - bounds.Min.Y
	width := bounds.Max.X - bounds.Min.X
	perimeter := 2.0 * float64(height+width)

	// Random sleep (0.1 to 0.4 seconds) to simulate GPU processing
	sleepTime := 100 + rand.Intn(300)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	return perimeter, nil
}

// Process a job in the background
func processJob(job *Job) {
	job.Status = "ongoing"
	job.StartTime = time.Now()

	var wg sync.WaitGroup
	resultChan := make(chan ImageResult)
	errorChan := make(chan StoreError)

	// Process each visit
	for _, visit := range job.Visits {
		// Check if store exists in master
		store, exists := storeMaster[visit.StoreID]
		if !exists {
			// Add error if store doesn't exist
			errorChan <- StoreError{
				StoreID: visit.StoreID,
				Error:   "store not found",
			}
			continue
		}

		// Process each image
		for _, imgURL := range visit.ImageURLs {
			wg.Add(1)
			go func(url, storeID, visitTime string, store Store) {
				defer wg.Done()
				perimeter, err := processImage(url, storeID, visitTime)
				if err != nil {
					errorChan <- StoreError{
						StoreID: storeID,
						Error:   err.Error(),
					}
					return
				}

				resultChan <- ImageResult{
					StoreID:   storeID,
					StoreName: store.StoreName,
					AreaCode:  store.AreaCode,
					ImageURL:  url,
					Perimeter: perimeter,
					VisitTime: visitTime,
				}
			}(imgURL, visit.StoreID, visit.VisitTime, store)
		}
	}

	// Close channels when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results and errors
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else {
				job.Results = append(job.Results, result)
			}
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
			} else {
				job.Errors = append(job.Errors, err)
			}
		}

		if resultChan == nil && errorChan == nil {
			break
		}
	}

	// Update job status
	job.EndTime = time.Now()
	if len(job.Errors) > 0 {
		job.Status = "failed"
	} else {
		job.Status = "completed"
	}

	// Log job completion
	log.Printf("Job %d completed with status: %s, processed %d images, errors: %d",
		job.ID, job.Status, len(job.Results), len(job.Errors))
}

// Handler for submit job endpoint
func submitJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format"})
		return
	}

	// Validate request
	if request.Count != len(request.Visits) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Count doesn't match number of visits"})
		return
	}

	// Create a new job
	jobIDMutex.Lock()
	jobID := nextJobID
	nextJobID++
	jobIDMutex.Unlock()

	job := &Job{
		ID:      jobID,
		Status:  "created",
		Visits:  request.Visits,
		Results: []ImageResult{},
		Errors:  []StoreError{},
	}

	// Store job
	jobsMutex.Lock()
	jobs[jobID] = job
	jobsMutex.Unlock()

	// Process job in the background
	go processJob(job)

	// Return job ID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JobResponse{JobID: jobID})
}

// Handler for job status endpoint
func getJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get job ID from query parameter
	jobIDStr := r.URL.Query().Get("jobid")
	if jobIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	// Get job
	jobsMutex.RLock()
	job, exists := jobs[jobID]
	jobsMutex.RUnlock()

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	// Prepare response
	response := StatusResponse{
		Status: job.Status,
		JobID:  jobIDStr,
	}

	if job.Status == "failed" {
		response.Errors = job.Errors
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Load store master data
	if err := loadStoreMaster(); err != nil {
		log.Fatalf("Failed to load store master: %v", err)
	}
	log.Printf("Loaded %d stores from store master", len(storeMaster))

	// Set up HTTP handlers
	http.HandleFunc("/api/submit/", submitJobHandler)
	http.HandleFunc("/api/status", getJobStatusHandler)

	// Start server
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
