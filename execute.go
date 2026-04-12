package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Known aux4/cache params to exclude when rebuilding the command
var cacheParams = map[string]bool{
	"cacheDuration": true,
	"cacheDir":      true,
	"configFile":    true,
	"command":       true,
	"packageDir":    true,
	"configDir":     true,
	"aux4HomeDir":   true,
	"help":          true,
	"local":         true,
	"refresh":       true,
}

func runExecute(flagDuration int, cacheDir string, configFilePath string, refresh bool, commandTokens []string, namedParams map[string]interface{}) {
	// Load config if provided
	var cfg *CacheConfig
	if configFilePath != "" {
		loaded, err := loadConfig(configFilePath)
		if err == nil {
			cfg = loaded
		}
	}

	// Build the full command string for cache key and execution
	commandArgs := buildCommandArgs(commandTokens, namedParams)
	commandStr := strings.Join(commandArgs, " ")

	// Resolve cache settings from config
	duration := resolveDuration(cfg, commandStr, flagDuration)
	cacheDir = resolveCacheDir(cfg, cacheDir)

	// Check cache (skip on refresh)
	key := generateCacheKey(commandStr)
	if !refresh {
		entry, err := readCache(cacheDir, key)
		if err == nil && isCacheValid(entry) {
			fmt.Print(entry.Output)
			return
		}
	}

	// Cache miss — execute the command
	cmd := exec.Command("aux4", commandArgs...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
		os.Exit(1)
	}

	outputStr := string(output)

	// Cache the result
	newEntry := CacheEntry{
		Command:   commandStr,
		Timestamp: time.Now().Unix(),
		Duration:  duration,
		Output:    outputStr,
	}
	writeCache(cacheDir, key, newEntry)

	fmt.Print(outputStr)
}

func buildCommandArgs(commandTokens []string, namedParams map[string]interface{}) []string {
	var args []string

	// Add command tokens (e.g., "db", "execute")
	args = append(args, commandTokens...)

	// Add named params as --flag value pairs, excluding cache-specific params
	for name, value := range namedParams {
		if cacheParams[name] {
			continue
		}
		args = append(args, fmt.Sprintf("--%s", name), fmt.Sprintf("%v", value))
	}

	return args
}
