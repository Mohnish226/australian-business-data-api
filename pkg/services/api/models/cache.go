package models

import "time"

// CacheEntry represents a cached API response
type CacheEntry struct {
	Data       interface{} `json:"data"`
	Timestamp  time.Time   `json:"timestamp"`
	Expiration time.Time   `json:"expiration"`
}
