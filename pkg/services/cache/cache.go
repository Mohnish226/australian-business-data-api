package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mohnish226/australian-business-data-api/pkg/config"
	"github.com/mohnish226/australian-business-data-api/pkg/services/api/models"
)

// initCacheDir creates the cache directory if it doesn't exist
func initCacheDir() error {
	if err := os.MkdirAll(config.CacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}
	return nil
}

// GetCache retrieves cached data for a given query and filters
func GetCache(query string, filters map[string]string) (interface{}, error) {
	if err := initCacheDir(); err != nil {
		return nil, err
	}

	cacheKey := generateCacheKey(query, filters)
	cacheFile := filepath.Join(config.CacheDir, cacheKey+".json")

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var entry models.CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}

	if time.Now().After(entry.Expiration) {
		os.Remove(cacheFile)
		return nil, fmt.Errorf("cache expired")
	}

	return entry.Data, nil
}

// SetCache stores data in the cache with expiration
func SetCache(query string, filters map[string]string, data interface{}) error {
	if err := initCacheDir(); err != nil {
		return err
	}

	cacheKey := generateCacheKey(query, filters)
	cacheFile := filepath.Join(config.CacheDir, cacheKey+".json")

	entry := models.CacheEntry{
		Data:       data,
		Timestamp:  time.Now(),
		Expiration: time.Now().Add(config.CacheExpiration),
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, jsonData, 0644)
}

// RemoveExpiredCache removes all expired cache entries
func RemoveExpiredCache() error {
	if err := initCacheDir(); err != nil {
		return err
	}

	files, err := os.ReadDir(config.CacheDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		cacheFile := filepath.Join(config.CacheDir, file.Name())
		data, err := os.ReadFile(cacheFile)
		if err != nil {
			continue
		}

		var entry models.CacheEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}

		if time.Now().After(entry.Expiration) {
			os.Remove(cacheFile)
		}
	}

	return nil
}

// CleanCache removes all cache entries
func CleanCache() error {
	if err := initCacheDir(); err != nil {
		return err
	}

	files, err := os.ReadDir(config.CacheDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			os.Remove(filepath.Join(config.CacheDir, file.Name()))
		}
	}

	return nil
}

// generateCacheKey creates a unique key for cache entries
func generateCacheKey(query string, filters map[string]string) string {
	key := query
	for k, v := range filters {
		key += fmt.Sprintf("_%s_%s", k, v)
	}
	return key
}
