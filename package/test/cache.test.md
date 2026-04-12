# cache

```file:config.yaml
config:
  test:
    host: localhost
    port: 5432
```

## execute with caching

```beforeEach
aux4 cache clear
```

```afterAll
aux4 cache clear
rm -f config.yaml
```

### should execute command and cache the result

```execute
aux4 cache config get test/host --file config.yaml
```

```expect
localhost
```

### should return cached result on second call

```execute
aux4 cache config get test/host --file config.yaml && aux4 cache config get test/host --file config.yaml
```

```expect
localhostlocalhost
```

### should cache with custom duration

```execute
aux4 cache --cacheDuration 10 config get test/port --file config.yaml && aux4 cache list
```

```expect:partial
5432[
**"duration": 10,**
]
```

### should refresh cache when --refresh is true

```execute
aux4 cache config get test/host --file config.yaml > /dev/null && aux4 cache --refresh true config get test/host --file config.yaml > /dev/null && aux4 cache list
```

```expect:partial
**"valid": true**
```

## list

```beforeEach
aux4 cache clear
```

### should return empty array when no cache entries

```execute
aux4 cache list
```

```expect
[]
```

### should list cached entries

```execute
aux4 cache config get test/host --file config.yaml > /dev/null && aux4 cache list
```

```expect:partial
**"command": "config get test/host --file config.yaml"**
```

## clear

### should clear all cached entries

```execute
aux4 cache config get test/host --file config.yaml > /dev/null && aux4 cache clear
```

```expect
Cleared 1 cached entries
```

### should clear entries matching pattern

```execute
aux4 cache config get test/host --file config.yaml > /dev/null && aux4 cache config get test/port --file config.yaml > /dev/null && aux4 cache clear "config get test/host *"
```

```expect
Cleared 1 cached entries matching "config get test/host *"
```
