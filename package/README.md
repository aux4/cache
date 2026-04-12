# aux4/cache

Cache the output of any aux4 command with a configurable TTL. On subsequent identical calls within the TTL, the cached result is returned instantly instead of re-executing the command.

## Installation

```bash
aux4 aux4 pkger install aux4/cache
```

## Quick Start

```bash
# Run a command with caching (default 300s TTL)
aux4 cache db execute --query "SELECT * FROM users"

# Second call returns cached result instantly
aux4 cache db execute --query "SELECT * FROM users"

# Set a custom TTL
aux4 cache --cacheDuration 60 curl request https://api.example.com/data
```

## Configuration

Create a `config.yaml` to define per-command cache duration patterns:

```yaml
config:
  cache:
    defaultDuration: 300
    cacheDir: .aux4.cache
    patterns:
      - match: "db execute *"
        duration: 300
      - match: "db stream *"
        duration: 120
      - match: "curl *"
        duration: 600
```

Use it with `--configFile`:

```bash
aux4 cache --configFile config.yaml db execute --query "SELECT 1"
```

### Duration Resolution Priority

1. Explicit `--cacheDuration` flag (if > 0)
2. First matching pattern from `config.yaml`
3. `defaultDuration` from `config.yaml`
4. Hardcoded fallback: 300 seconds

### Pattern Matching

Patterns use prefix matching with `*` as a wildcard suffix:

- `"db execute *"` matches any command starting with `db execute`
- `"curl *"` matches any command starting with `curl`
- `"config get test/host"` matches that exact command

## Commands

### `aux4 cache <command>`

Execute a command with caching. If a cached result exists and is within the TTL, it is returned without re-executing. Otherwise, the command runs, its output is cached, and the result is printed.

```bash
aux4 cache [--cacheDuration <seconds>] [--cacheDir <dir>] [--configFile <file>] <command...>
```

| Flag | Default | Description |
|------|---------|-------------|
| `--cacheDuration` | `0` (auto) | Cache duration in seconds. `0` means resolve from config or use 300s default |
| `--cacheDir` | `.aux4.cache` | Directory for cache files |
| `--configFile` | (none) | Path to config.yaml for pattern-based durations |
| `--refresh` | `false` | Force re-execute and refresh the cached entry |

### `aux4 cache clear [pattern]`

Remove cached entries. Without a pattern, all entries are removed. With a pattern, only matching entries are cleared.

```bash
# Clear all cached entries
aux4 cache clear

# Clear entries matching a pattern
aux4 cache clear "db execute *"
```

### `aux4 cache list`

List all cached entries with their age, TTL, and validity status. Output is JSON.

```bash
aux4 cache list
```

```json
[
  {
    "command": "db execute --query SELECT 1",
    "age": 45,
    "duration": 300,
    "remaining": 255,
    "valid": true
  }
]
```

## Cache Storage

Cache files are stored as JSON in the cache directory (default `.aux4.cache/`). Each entry is a file named by the SHA256 hash of the command string:

```
.aux4.cache/
  a1b2c3d4...json
  e5f6g7h8...json
```

Writes are atomic (temp file + rename) to prevent partial reads.
