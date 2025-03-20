package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mohnish226/australian-business-data-api/pkg/config"
	"github.com/mohnish226/australian-business-data-api/pkg/logger"
	"github.com/mohnish226/australian-business-data-api/pkg/services/api/models"
	"github.com/mohnish226/australian-business-data-api/pkg/services/cache"
)

// Service handles API interactions
type Service struct {
	client *http.Client
}

// NewService creates a new API service instance
func NewService() *Service {
	return &Service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// BasicSearch performs a basic search using the datastore_search endpoint
func (s *Service) BasicSearch(query string, filters map[string]string) ([]map[string]interface{}, error) {
	logger.Logger.Printf("Starting basic search with query: %s, filters: %v", query, filters)

	// Check cache first
	if cached, err := cache.GetCache(query, filters); err == nil {
		logger.Logger.Printf("Cache hit for query: %s", query)
		// Safely convert cached data to []map[string]interface{}
		jsonData, err := json.Marshal(cached)
		if err != nil {
			logger.Logger.Printf("Failed to marshal cached data: %v", err)
			return nil, fmt.Errorf("failed to marshal cached data: %v", err)
		}

		var result []map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			logger.Logger.Printf("Failed to unmarshal cached data: %v", err)
			return nil, fmt.Errorf("failed to unmarshal cached data: %v", err)
		}

		return result, nil
	}
	logger.Logger.Printf("Cache miss for query: %s", query)

	// Build request body
	requestBody := map[string]interface{}{
		"resource_id": config.ResourceID,
		"limit":       config.RequestLimit,
		"offset":      0,
		"q":           query,
		"filters":     filters,
	}

	// Convert request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		logger.Logger.Printf("Failed to marshal request body: %v", err)
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create request
	url := fmt.Sprintf("%s%s", config.Host, config.RestPath)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Logger.Printf("Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make request
	logger.Logger.Printf("Making POST request to: %s with body: %s", url, string(jsonBody))
	resp, err := s.client.Do(req)
	if err != nil {
		logger.Logger.Printf("Request failed: %v", err)
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Printf("Failed to read response: %v", err)
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Log raw response for debugging
	logger.Logger.Printf("Raw response: %s", string(body))

	// Parse response
	var searchResp models.SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		logger.Logger.Printf("Failed to parse response: %v", err)
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Extract records from response
	records, ok := searchResp.Result["records"].([]interface{})
	if !ok {
		logger.Logger.Printf("Invalid response format: records not found in result: %+v", searchResp.Result)
		return nil, fmt.Errorf("invalid response format: records not found")
	}

	// Convert records to []map[string]interface{}
	var result []map[string]interface{}
	for _, record := range records {
		// First convert to JSON and back to ensure proper type conversion
		jsonData, err := json.Marshal(record)
		if err != nil {
			logger.Logger.Printf("Failed to marshal record: %v", err)
			continue
		}

		var recordMap map[string]interface{}
		if err := json.Unmarshal(jsonData, &recordMap); err != nil {
			logger.Logger.Printf("Failed to unmarshal record: %v", err)
			continue
		}

		result = append(result, recordMap)
	}

	logger.Logger.Printf("Found %d records for query: %s", len(result), query)

	// Cache the response
	cache.SetCache(query, filters, result)
	logger.Logger.Printf("Cached results for query: %s", query)

	return result, nil
}

// SQLSearch performs a search using the datastore_search_sql endpoint
func (s *Service) SQLSearch(query string) ([]map[string]interface{}, error) {
	logger.Logger.Printf("Starting SQL search with query: %s", query)

	// Check cache first
	if cached, err := cache.GetCache(query, nil); err == nil {
		logger.Logger.Printf("Cache hit for SQL query: %s", query)
		// Safely convert cached data to []map[string]interface{}
		jsonData, err := json.Marshal(cached)
		if err != nil {
			logger.Logger.Printf("Failed to marshal cached data: %v", err)
			return nil, fmt.Errorf("failed to marshal cached data: %v", err)
		}

		var result []map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			logger.Logger.Printf("Failed to unmarshal cached data: %v", err)
			return nil, fmt.Errorf("failed to unmarshal cached data: %v", err)
		}

		return result, nil
	}
	logger.Logger.Printf("Cache miss for SQL query: %s", query)

	// Build request body
	requestBody := map[string]interface{}{
		"resource_id": config.ResourceID,
		"sql":         query,
	}

	// Convert request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		logger.Logger.Printf("Failed to marshal SQL request body: %v", err)
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create request
	url := fmt.Sprintf("%s%s", config.Host, config.SQLPath)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Logger.Printf("Failed to create SQL request: %v", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make request
	logger.Logger.Printf("Making SQL POST request to: %s with body: %s", url, string(jsonBody))
	resp, err := s.client.Do(req)
	if err != nil {
		logger.Logger.Printf("SQL request failed: %v", err)
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Printf("Failed to read SQL response: %v", err)
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Log raw response for debugging
	logger.Logger.Printf("Raw SQL response: %s", string(body))

	// Parse response
	var searchResp models.SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		logger.Logger.Printf("Failed to parse SQL response: %v", err)
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Extract records from response
	records, ok := searchResp.Result["records"].([]interface{})
	if !ok {
		logger.Logger.Printf("Invalid SQL response format: records not found in result: %+v", searchResp.Result)
		return nil, fmt.Errorf("invalid response format: records not found")
	}

	// Convert records to []map[string]interface{}
	var result []map[string]interface{}
	for _, record := range records {
		// First convert to JSON and back to ensure proper type conversion
		jsonData, err := json.Marshal(record)
		if err != nil {
			logger.Logger.Printf("Failed to marshal record: %v", err)
			continue
		}

		var recordMap map[string]interface{}
		if err := json.Unmarshal(jsonData, &recordMap); err != nil {
			logger.Logger.Printf("Failed to unmarshal record: %v", err)
			continue
		}

		result = append(result, recordMap)
	}

	logger.Logger.Printf("Found %d records for SQL query: %s", len(result), query)

	// Cache the response
	cache.SetCache(query, nil, result)
	logger.Logger.Printf("Cached SQL results for query: %s", query)

	return result, nil
}

// GetBusinesses retrieves business data using the specified search method
func (s *Service) GetBusinesses(query string, filters map[string]string, useSQL bool) ([]models.Business, error) {
	logger.Logger.Printf("Getting businesses with query: %s, filters: %v, useSQL: %v", query, filters, useSQL)

	var records []map[string]interface{}
	var err error

	if useSQL {
		records, err = s.SQLSearch(query)
	} else {
		records, err = s.BasicSearch(query, filters)
	}

	if err != nil {
		logger.Logger.Printf("Failed to get businesses: %v", err)
		return nil, err
	}

	// Convert records to Business structs
	var businesses []models.Business
	for _, record := range records {
		jsonData, err := json.Marshal(record)
		if err != nil {
			logger.Logger.Printf("Failed to marshal record: %v", err)
			continue
		}

		var business models.Business
		if err := json.Unmarshal(jsonData, &business); err != nil {
			logger.Logger.Printf("Failed to unmarshal business: %v", err)
			continue
		}

		businesses = append(businesses, business)
	}

	logger.Logger.Printf("Successfully converted %d records to businesses", len(businesses))
	return businesses, nil
}
