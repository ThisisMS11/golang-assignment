package job

import (
	"log"
	"sync"
	"time"
	"kiranaclub/internal/models"
	"kiranaclub/internal/processor"
	"kiranaclub/internal/store"
)

var (
	jobs       = make(map[int]*models.Job)
	jobsMutex  sync.RWMutex
	nextJobID  = 1
	jobIDMutex sync.Mutex
)

// Create creates a new job and starts processing it
func Create(visits []models.Visit) int {
	// Generate job ID
	jobIDMutex.Lock()
	jobID := nextJobID
	nextJobID++
	jobIDMutex.Unlock()

	// Create job
	job := &models.Job{
		ID:      jobID,
		Status:  "created",
		Visits:  visits,
		Results: []models.ImageResult{},
		Errors:  []models.StoreError{},
	}

	// Store job
	jobsMutex.Lock()
	jobs[jobID] = job
	jobsMutex.Unlock()

	// Process job in background
	go processJob(job)

	return jobID
}

// Get retrieves a job by ID
func Get(jobID int) (*models.Job, bool) {
	jobsMutex.RLock()
	defer jobsMutex.RUnlock()
	job, exists := jobs[jobID]
	return job, exists
}

// Get all jobs
func GetAll() ([]models.Job) {
	jobsMutex.RLock()
	defer jobsMutex.RUnlock()
	jobsSlice := make([]models.Job, 0, len(jobs))
	for _, job := range jobs {
		jobsSlice = append(jobsSlice, *job)
	}
	return jobsSlice
}

// processJob handles the processing of a job
func processJob(job *models.Job) {
	job.Status = "ongoing"
	job.StartTime = time.Now()

	log.Printf("Made this job ongoing");

	var wg sync.WaitGroup
	resultChan := make(chan models.ImageResult)
	errorChan := make(chan models.StoreError)

	log.Printf("Store Visits: %+v", job.Visits)

	// Process each visit
	for _, visit := range job.Visits {
		// Check if store exists in master
		log.Printf("Checking store with StoreID: %s", visit.StoreID)
		storeInfo, exists := store.GetStore(visit.StoreID)

		log.Printf("storeInfo : %+v", storeInfo)
		if !exists {
			// Add error if store doesn't exist
			errorChan <- models.StoreError{
				StoreID: visit.StoreID,
				Error:   "store not found",
			}
			log.Printf("Store not found for StoreID: %s", visit.StoreID)
			continue
		}

		log.Printf("Store Info for StoreID %s: %+v", visit.StoreID, storeInfo)

		// Process each image
		for _, imgURL := range visit.ImageURLs {
			wg.Add(1)
			go func(url, storeID, visitTime string, storeInfo models.Store) {
				defer wg.Done()
				perimeter, err := processor.ProcessImage(url)
				if err != nil {
					errorChan <- models.StoreError{
						StoreID: storeID,
						Error:   err.Error(),
					}
					return
				}

				resultChan <- models.ImageResult{
					StoreID:   storeID,
					StoreName: storeInfo.StoreName,
					AreaCode:  storeInfo.AreaCode,
					ImageURL:  url,
					Perimeter: perimeter,
					VisitTime: visitTime,
				}
			}(imgURL, visit.StoreID, visit.VisitTime, storeInfo)
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