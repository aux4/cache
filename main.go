package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "Usage: aux4-cache <cacheDuration> <cacheDir> <configFile> <refresh> <commandJSON> [paramsJSON]\n")
		os.Exit(1)
	}

	cacheDuration, err := strconv.Atoi(os.Args[1])
	if err != nil {
		cacheDuration = 0
	}

	cacheDir := os.Args[2]
	configFilePath := os.Args[3]
	refresh := os.Args[4] == "true"

	var commandTokens []string
	if len(os.Args) > 5 {
		err = json.Unmarshal([]byte(os.Args[5]), &commandTokens)
		if err != nil {
			// Single string, not JSON array
			commandTokens = []string{os.Args[5]}
		}
	}

	var namedParams map[string]interface{}
	if len(os.Args) > 6 {
		json.Unmarshal([]byte(os.Args[6]), &namedParams)
	}
	if namedParams == nil {
		namedParams = make(map[string]interface{})
	}

	if len(commandTokens) == 0 {
		fmt.Fprintf(os.Stderr, "No command specified\n")
		os.Exit(1)
	}

	// Dispatch based on first command token
	switch commandTokens[0] {
	case "clear":
		pattern := ""
		if len(commandTokens) > 1 {
			pattern = commandTokens[1]
		}
		runClear(cacheDir, configFilePath, pattern)
	case "list":
		runList(cacheDir, configFilePath)
	default:
		runExecute(cacheDuration, cacheDir, configFilePath, refresh, commandTokens, namedParams)
	}
}
