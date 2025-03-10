package api

import (
	"encoding/json"
	"kiranaclub/internal/job"
	"kiranaclub/internal/models"
	"log"
	"net/http"
	"strconv"
)

// SubmitJobHandler handles job submission requests
func SubmitJobHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("http method : %s",r.Method);
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request models.SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid request format"})
		return
	}

	/* validating request */
	if request.Count != len(request.Visits) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Count doesn't match number of visits"})
		return
	}

	jobID := job.Create(request.Visits)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.JobResponse{JobID: jobID})
}

// GetJobStatusHandler handles job status requests
func GetJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	jobInfo, exists := job.Get(jobID)
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	response := models.StatusResponse{
		Status: jobInfo.Status,
		JobID:  jobIDStr,
	}

	if jobInfo.Status == "failed" {
		response.Errors = jobInfo.Errors
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetJobsInformation handles GET requests to retrieve job information
func GetJobsInformation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	jobs := job.GetAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(jobs) == 0 {
		json.NewEncoder(w).Encode([]models.Job{}) 
		return
	}

	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
