
# rawkit

A Go library for reading and processing RAW image files. Wraps [LibRaw](https://www.libraw.org/) via cgo with fallback support for standard image formats (JPEG, PNG).

```
go get github.com/tlaceby/rawkit@latest
```

## Supported Formats

**RAW formats:** ARW (Sony), CR2/CR3 (Canon), NEF (Nikon), DNG (Adobe), ORF (Olympus), RAF (Fujifilm), RW2 (Panasonic)

**Standard formats:** JPEG, PNG

## Requirements

LibRaw must be installed on your system.

```bash
# macOS
brew install libraw

# Ubuntu/Debian
sudo apt install libraw-dev

# Fedora
sudo dnf install LibRaw-devel
```

## Packages

| Package | Description |
|---------|-------------|
| `core` | Image reading, metadata extraction, pixel data |
| `histogram` | Histogram generation and analysis |
| `processing` | Image adjustments and transformations |

For detailed API documentation, see [pkg.go.dev/github.com/tlaceby/rawkit](https://pkg.go.dev/github.com/tlaceby/rawkit) or run:

```bash
go doc github.com/tlaceby/rawkit/core
go doc github.com/tlaceby/rawkit/histogram
go doc github.com/tlaceby/rawkit/processing
```

## Examples

### Read an Image

```go
package main

import (
    "fmt"
    "log"

    "github.com/tlaceby/rawkit/core"
)

func main() {
    img, err := core.ReadAll("/path/to/photo.arw")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Dimensions: %dx%d\n", img.Data.Width, img.Data.Height)
    fmt.Printf("Bit Depth: %d\n", img.Data.BitDepth)

    if img.Meta != nil {
        fmt.Printf("Camera: %s %s\n", img.Meta.CameraMake, img.Meta.CameraModel)
        fmt.Printf("ISO: %d, f/%.1f, %.4fs\n", img.Meta.ISO, img.Meta.Aperture, img.Meta.SS)
    }
}
```

### Metadata Only (Fast)

```go
meta, err := core.ReadRAWMetadata("/path/to/photo.arw")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Camera: %s %s\n", meta.CameraMake, meta.CameraModel)
```

### Pixel Data Only

```go
data, err := core.ReadImageData("/path/to/photo.arw")
if err != nil {
    log.Fatal(err)
}

// Access pixel at (100, 200)
r, g, b := data.Pixel(100, 200)
```

### Pixel Data Format

All pixel data is `[]uint16` in RGB layout:

```go
// Using helper methods
r, g, b := data.Pixel(x, y)

// Direct buffer access
idx := data.Index(x, y)
r, g, b := data.Data[idx], data.Data[idx+1], data.Data[idx+2]
```

## License

MIT
