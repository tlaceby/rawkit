
# RawKit

**RawKit** is a tiny CGO wrapper around **[LibRaw](https://www.libraw.org/)** that lets Go programs open RAW files *without* compiling C. This means LibRaw is not required for installation to use this library. RawKit provides a simple yet stable API which works with RAW image file formats. Allowing you to do powerful things right from the Go language.

---

## ✨ Quick Start
To get started with **rawkit**, simple install it with the built in go tooling. You dont need to install **LibRaw**.
```bash
    go get github.com/tlaceby/rawkit@latest
```

```go
package main

import (
	"fmt"
	"github.com/tlaceby/rawkit"
)

func main() {
	img, err := rawkit.LoadRAW("example.ARW")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%dx%d  %d-channel  %d samples\n",img.Width, img.Height, img.Channels, len(img.Buffer)) // channels --> (1: RAW, 3: RGB)
	fmt.Println("LibRaw Version:", rawkit.LibRawVersionStr())
}
```

---

## 📚 Table of Contents

1. [Core API](#core-api)
2. [Building / Re-building Binaries](#building--re-building-binaries)
3. [Extending RawKit](#extending-rawkit)
4. [Contributing](#contributing)

---

## Core API

| Function             | Description                                                                                                  |
| -------------------- | ------------------------------------------------------------------------------------------------------------ |
| `LoadRAW(path)`      | Opens a RAW file and returns:<br>`*RawKitImage` containing `Buffer []uint16`, `Width`, `Height`, `Channels`. |
| `LibRawVersionNum()` | Numeric version `0xMMmmpp`.                                                                                  |
| `LibRawVersionStr()` | Human string, e.g. `0.22.0-Release`.                                                                         |

All LibRaw allocations are released inside the wrapper; the returned Go slice is fully owned by the caller.

---

## 🔨 Building / Re-building Binaries

### One-shot pipeline

```bash
make release
```

| Phase     | What happens                                              |
| --------- | --------------------------------------------------------- |
| **clean** | Deletes old `.a` files & Go cache                         |
| **build** | Compiles LibRaw + wrapper for the **host** OS/arch        |
| **test**  | Runs `go test ./...`                                      |
| **bump**  | Writes next `VERSION=` to `.env` **only if tests passed** |

### Per-platform scripts

```bash
# macOS (Intel or Apple Silicon)
bash scripts/build_darwin.sh  vX.Y.Z  [arm64|amd64]

# Linux
bash scripts/build_linux.sh   vX.Y.Z  [amd64|arm64]

# Windows (cross-build; needs mingw-w64)
bash scripts/build_windows.sh vX.Y.Z
```

Each script builds into `libs/<os_arch>/<version>/` and then updates the
symlink `libs/<os_arch>/current → <version>`.

---

## 🌱 Extending RawKit

| Step         | Folder     | Action                                                      |
| ------------ | ---------- | ----------------------------------------------------------- |
| **1 Bridge** | `wrapper/` | Add an `extern "C"` function in `libraw_wrapper.cpp/.h`.    |
| **2 Expose** | `rawkit/`  | Declare the symbol in a Go file’s cgo preamble and wrap it. |
| **3 Test**   | `tests/`   | Add/extend `*_test.go`; place samples in `tests/testdata/`. |
| **4 Verify** | repo root  | `make release` — builds, tests, bumps version.              |

> **Tip:** add new fields to existing structs instead of changing old ones to stay backward-compatible.

---

## 🤝 Contributing
My goal for **rawkit** is simple: Use it with other astrophotography application I am writing. I dont need much functionality from ImageRaw as I am writing the rest myself, but if anyone wants to make a pull request or keep this repositiory updated with feature requests, I am happy to continue the work as long as the API remains stable.

* **Bug reports / feature requests** – open an issue with steps to reproduce.
* **Pull requests**

  1. Fork → branch → commit (ensure `go test ./...` is green).
  2. Run `make release` locally; the version bump only lands if tests pass.
  3. Open a PR – GitHub Actions runs the same pipeline.
* **Style** – idiomatic Go (`gofmt`), minimal C++17 in the wrapper.
* **Large RAW assets** – dont! :)

