package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type CacheConfig struct {
	DefaultDuration int            `yaml:"defaultDuration"`
	CacheDir        string         `yaml:"cacheDir"`
	Patterns        []CachePattern `yaml:"patterns"`
}

type CachePattern struct {
	Match    string `yaml:"match"`
	Duration int    `yaml:"duration"`
}

type configFile struct {
	Config struct {
		Cache CacheConfig `yaml:"cache"`
	} `yaml:"config"`
}

func loadConfig(path string) (*CacheConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg configFile
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg.Config.Cache, nil
}

// resolveDuration resolves cache duration with priority:
// 1. Explicit --cacheDuration flag (> 0)
// 2. Matching config pattern
// 3. Config defaultDuration
// 4. Hardcoded fallback (300s)
func resolveDuration(cfg *CacheConfig, command string, flagDuration int) int {
	// Explicit flag always wins
	if flagDuration > 0 {
		return flagDuration
	}

	if cfg != nil {
		// Check patterns in order, return first match
		for _, p := range cfg.Patterns {
			if matchPattern(p.Match, command) {
				return p.Duration
			}
		}

		// Fall back to config default
		if cfg.DefaultDuration > 0 {
			return cfg.DefaultDuration
		}
	}

	return 300
}

func resolveCacheDir(cfg *CacheConfig, flagCacheDir string) string {
	if flagCacheDir != "" && flagCacheDir != ".aux4.cache" {
		return flagCacheDir
	}

	if cfg != nil && cfg.CacheDir != "" {
		return cfg.CacheDir
	}

	if flagCacheDir != "" {
		return flagCacheDir
	}

	return ".aux4.cache"
}
