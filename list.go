package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type CacheListItem struct {
	Command   string `json:"command"`
	Age       int64  `json:"age"`
	Duration  int    `json:"duration"`
	Remaining int64  `json:"remaining"`
	Valid     bool   `json:"valid"`
}

func runList(cacheDir string, configFilePath string) {
	var cfg *CacheConfig
	if configFilePath != "" {
		loaded, err := loadConfig(configFilePath)
		if err == nil {
			cfg = loaded
		}
	}

	cacheDir = resolveCacheDir(cfg, cacheDir)

	entries, err := listCacheEntries(cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing cache: %s\n", err)
		os.Exit(1)
	}

	now := time.Now().Unix()
	var items []CacheListItem
	for _, entry := range entries {
		age := now - entry.Timestamp
		remaining := int64(entry.Duration) - age
		if remaining < 0 {
			remaining = 0
		}
		items = append(items, CacheListItem{
			Command:   entry.Command,
			Age:       age,
			Duration:  entry.Duration,
			Remaining: remaining,
			Valid:     remaining > 0,
		})
	}

	if items == nil {
		items = []CacheListItem{}
	}

	output, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
