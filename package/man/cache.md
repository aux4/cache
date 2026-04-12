#### Description

The `cache` command wraps any aux4 command with a caching layer. When executed, it checks if a valid cached result exists for the given command. On a cache hit, the cached output is returned instantly. On a cache miss, the command runs via `aux4`, the output is stored in the cache directory, and the result is printed.

The cache key is the SHA256 hash of the full reconstructed command string (command tokens + named parameters). Cache duration can be set per-invocation via `--cacheDuration`, or configured per-command-pattern in a `config.yaml` file.

Subcommands:

- **clear** — Remove cached entries (all or by pattern)
- **list** — List cached entries with age and validity

#### Usage

```bash
aux4 cache [--cacheDuration <seconds>] [--cacheDir <dir>] [--configFile <file>] [--refresh <true|false>] <command...>
```

--cacheDuration  Cache duration in seconds. 0 means auto-resolve from config or default 300s
--cacheDir       Directory for cache files (default: .aux4.cache)
--configFile     Path to config.yaml with cache patterns
--refresh        Force re-execute and refresh the cached entry (default: false)

#### Example

```bash
aux4 cache config get dev/host --file config.yaml
```

```text
localhost
```

Running the same command again returns the cached result without executing `aux4 config get`:

```bash
aux4 cache config get dev/host --file config.yaml
```

```text
localhost
```
