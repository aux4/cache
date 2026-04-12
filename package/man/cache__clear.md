#### Description

The `clear` subcommand removes cached entries from the cache directory. Without a pattern argument, all cached entries are removed. With a pattern, only entries whose command string matches the pattern are deleted.

Patterns support `*` as a suffix wildcard. For example, `"db execute *"` matches any cached command starting with `db execute`.

#### Usage

```bash
aux4 cache clear [pattern] [--cacheDir <dir>] [--configFile <file>]
```

--cacheDir    Directory for cache files (default: .aux4.cache)
--configFile  Path to config.yaml for cache directory override

#### Example

Clear all cached entries:

```bash
aux4 cache clear
```

```text
Cleared 5 cached entries
```

Clear only entries matching a pattern:

```bash
aux4 cache clear "db execute *"
```

```text
Cleared 2 cached entries matching "db execute *"
```
