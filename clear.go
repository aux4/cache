package main

import (
	"fmt"
	"os"
)

func runClear(cacheDir string, configFilePath string, pattern string) {
	var cfg *CacheConfig
	if configFilePath != "" {
		loaded, err := loadConfig(configFilePath)
		if err == nil {
			cfg = loaded
		}
	}

	cacheDir = resolveCacheDir(cfg, cacheDir)

	count, err := clearCacheEntries(cacheDir, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error clearing cache: %s\n", err)
		os.Exit(1)
	}

	if pattern == "" {
		fmt.Printf("Cleared %d cached entries\n", count)
	} else {
		fmt.Printf("Cleared %d cached entries matching \"%s\"\n", count, pattern)
	}
}
