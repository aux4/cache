#### Description

The `list` subcommand displays all cached entries in the cache directory as a JSON array. Each entry shows the original command, its age in seconds, the configured TTL duration, remaining time, and whether it is still valid.

This is useful for inspecting what is currently cached and identifying stale entries.

#### Usage

```bash
aux4 cache list [--cacheDir <dir>] [--configFile <file>]
```

--cacheDir    Directory for cache files (default: .aux4.cache)
--configFile  Path to config.yaml for cache directory override

#### Example

```bash
aux4 cache list
```

```json
[
  {
    "command": "config get dev/host --file config.yaml",
    "age": 45,
    "duration": 300,
    "remaining": 255,
    "valid": true
  },
  {
    "command": "curl request https://api.example.com/data",
    "age": 650,
    "duration": 600,
    "remaining": 0,
    "valid": false
  }
]
```
