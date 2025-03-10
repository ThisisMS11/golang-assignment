package store

import (
	"encoding/csv"
	"kiranaclub/internal/models"
	"log"
	"os"
	"sync"
)

var (
	storeMaster     = make(map[string]models.Store)
	storeMasterLock sync.RWMutex
)

// GetStore retrieves a store by ID
func GetStore(storeID string) (models.Store, bool) {
	storeMasterLock.RLock()
	defer storeMasterLock.RUnlock()
	store, exists := storeMaster[storeID]
	return store, exists
}

// Count returns the number of stores in the master
func Count() int {
	storeMasterLock.RLock()
	defer storeMasterLock.RUnlock()
	return len(storeMaster)
}

// LoadStoreMaster loads stores from a CSV file
func LoadStoreMaster(filePath string) error {
	// Create store master file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Sample data
		log.Printf("File do not exist hence creating one")
		sampleStores := [][]string{
			{"StoreID", "StoreName", "AreaCode"},
			{"S00339218", "Walmart Central", "A001"},
			{"S01408764", "Target Downtown", "A002"},
			{"S02547931", "Kroger Heights", "A003"},
		}

		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		if err := writer.WriteAll(sampleStores); err != nil {
			return err
		}
	}

	// Read and parse the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row
	if len(records) < 2 {
		return nil
	}

	// Populate the store map
	storeMasterLock.Lock()
	defer storeMasterLock.Unlock()

	for _, record := range records[1:] {
		storeMaster[record[2]] = models.Store{
			StoreID:   record[2],
			StoreName: record[1],
			AreaCode:  record[0],
		}
	}

	return nil
}
