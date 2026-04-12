package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CacheEntry struct {
	Command   string `json:"command"`
	Timestamp int64  `json:"timestamp"`
	Duration  int    `json:"duration"`
	Output    string `json:"output"`
}

func generateCacheKey(command string) string {
	hash := sha256.Sum256([]byte(command))
	return fmt.Sprintf("%x", hash)
}

func readCache(cacheDir, key string) (*CacheEntry, error) {
	path := filepath.Join(cacheDir, key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entry CacheEntry
	err = json.Unmarshal(data, &entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func writeCache(cacheDir, key string, entry CacheEntry) error {
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return err
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	tmpPath := filepath.Join(cacheDir, key+".tmp")
	finalPath := filepath.Join(cacheDir, key+".json")

	err = os.WriteFile(tmpPath, data, 0644)
	if err != nil {
		return err
	}

	return os.Rename(tmpPath, finalPath)
}

func isCacheValid(entry *CacheEntry) bool {
	return time.Now().Unix()-entry.Timestamp < int64(entry.Duration)
}

func listCacheEntries(cacheDir string) ([]CacheEntry, error) {
	entries, err := filepath.Glob(filepath.Join(cacheDir, "*.json"))
	if err != nil {
		return nil, err
	}

	var result []CacheEntry
	for _, path := range entries {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var entry CacheEntry
		err = json.Unmarshal(data, &entry)
		if err != nil {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}

func clearCacheEntries(cacheDir, pattern string) (int, error) {
	entries, err := filepath.Glob(filepath.Join(cacheDir, "*.json"))
	if err != nil {
		return 0, err
	}

	count := 0
	for _, path := range entries {
		if pattern == "" {
			os.Remove(path)
			count++
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var entry CacheEntry
		err = json.Unmarshal(data, &entry)
		if err != nil {
			continue
		}

		if matchPattern(pattern, entry.Command) {
			os.Remove(path)
			count++
		}
	}

	return count, nil
}

func matchPattern(pattern, command string) bool {
	if pattern == "" {
		return true
	}

	// Support glob-style * at end
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(command) >= len(prefix) && command[:len(prefix)] == prefix
	}

	return pattern == command
}
