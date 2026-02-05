# Contributing

## Install libraw

**Linux**
```bash
sudo apt install libraw-dev
```

**macOS**
```bash
brew install libraw
```

**Windows**

Install libraw and set `LIBRAW_PATH` to where you installed it.

## Build

```bash
go build ./...
```

## Testing

```bash
go test -v ./...
```

## Adding functions

1. Add declaration to `rk/core/helpers.h`
2. Add implementation to `rk/core/helpers.c`
3. Expose it in `rk/rk.go`