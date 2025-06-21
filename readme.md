# RawKit
**Rawkit** is a tiny CGO wrapper around **[LibRaw](https://www.libraw.org/)** that lets Go programs open and manipulate RAW image files **without** needing a C compiler or external dependencies. Rawkit provides a simple yet stable API for working with RAW image formats — letting you do powerful things directly in Go.

---

## ✨ Quick Start

To get started with **rawkit**, install it using Go tooling. You **do not** need to install LibRaw or a C++ compiler — rawkit is fully bundled and statically compiled.

```bash
go get github.com/tlaceby/rawkit@latest
```

### 🚀 Load and Inspect RAW Metadata

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

	fmt.Printf("Image: %dx%d | Channels: %d | Buffer Size: %d\n", img.Width, img.Height, img.Colors, len(img.Buffer))
	fmt.Println("Camera:", img.CameraMake, img.CameraModel)
	fmt.Println("Lens:", img.FocalLength, "mm  f/", img.Aperature)
	fmt.Println("Exposure:", img.ShutterSpeed, "sec  ISO", img.ISO)
	fmt.Println("Artist:", img.Artist)
	fmt.Println("LibRaw Version:", rawkit.LibRawVersionStr())
}
```

### 🧠 Extracted Metadata Includes

* 📸 **Camera Make & Model**: `img.CameraMake`, `img.CameraModel`
* 🎨 **Color Space**: `img.ColorSpace` (e.g., sRGB, Adobe)
* 🧑‍🎨 **Artist Field**: `img.Artist`
* 🌅 **Exposure Settings**:

  * `img.ShutterSpeed` (sec)
  * `img.Aperature` (f/stop)
  * `img.FocalLength` (mm)
  * `img.ISO`
* 🧩 **DNG Info & RAW Count**: `img.DNGVersion`, `img.RawCount`
* 📐 **Image Orientation**: `img.Flip` (0, 3, 5, or 6)
* ✅ **White Balance Applied**: `img.AsShotWBApplied` (0 or 1)

All pixel data is accessible in `img.Buffer`, a `[]uint16` representing unpacked pixels in **channel-major** order. Perfect for custom RAW processing, histogram generation, or astrophotography pipelines.

---

## 📚 Table of Contents

1. [Go API Documentation](./docs.md)
2. [Contributing](#-contributing)
3. [Extending RawKit](#-extending-rawkit)
4. [Building RawKit](#-building--re-building-binaries)

---

## 🤝 Contributing

**Rawkit** was built to power a larger astrophotography application. While its scope is intentionally minimal, I welcome contributions from others who want to extend or improve it.

* **Bug reports / feature requests** — open an issue with steps to reproduce.
* **Pull requests**

  1. Fork → branch → commit (ensure `go test ./...` passes).
  2. Run `make release` — the version bump only applies if tests pass.
  3. Open a PR — GitHub Actions will run the same build pipeline.
* **Code style**

  * Use idiomatic Go (`gofmt`, `go vet`, etc.).
  * Keep C++ code minimal (C++17 or below).
* **Test assets**

  * Don’t commit large RAW files. Use minimal test samples (<1MB) if needed.

---

## 🌱 Extending RawKit

Want to expose new metadata or functionality? Here's the step-by-step workflow:

| Step          | Folder     | Description                                                          |
| ------------- | ---------- | -------------------------------------------------------------------- |
| **1. Bridge** | `wrapper/` | Add a new `extern "C"` function to `libraw_wrapper.cpp/.h`.          |
| **2. Expose** | `rawkit/`  | Declare the function in a CGO preamble, and wrap it in idiomatic Go. |
| **3. Test**   | `tests/`   | Add a Go test and put test files in `tests/testdata/`.               |
| **4. Verify** | repo root  | Run `make release` to build, test, and bump the version.             |

> 💡 **Tip:** Add new fields instead of modifying existing ones to maintain backward compatibility.

---

## 🔨 Building / Re-building Binaries

> 🛠️ Only needed for contributors — precompiled binaries are included for all major platforms.

### 🚧 Development Builds

```bash
make verify
```

Use this during development and testing. It:

* Builds native dependencies
* Regenerates all CGO bindings
* Recompiles the wrapper modules
* Runs `go test ./...` to ensure stability

This is the recommended way to validate changes before finalizing a release.

---

### 🚀 Releasing a New Version

```bash
make release
```

This runs the full pipeline:

| Phase   | Description                                               |
| ------- | --------------------------------------------------------- |
| `clean` | Deletes old `.a` files and Go build cache                 |
| `build` | Compiles LibRaw and the wrapper for your **host OS/arch** |
| `test`  | Runs `go test ./...`                                      |
| `bump`  | Writes a new `VERSION=` to `.env` if tests pass           |

Once development is complete, `make release` finalizes the version bump in `.env` and commits the changes. When pushed to GitHub, this triggers the CI workflow, which builds binaries for:

* 🐧 Linux (amd64 + arm64)
* 🪟 Windows (amd64) ---> Not yet Supported/Tested (TODO)
* 🍎 macOS (amd64 + arm64)

---

### ⚙️ Platform-Specific Scripts

Use these for targeted builds or cross-compilation:

```bash
go generate ./...

# macOS (Intel or Apple Silicon)
bash scripts/build_darwin.sh   vX.Y.Z  [arm64|amd64]

# Linux
bash scripts/build_linux.sh    vX.Y.Z  [amd64|arm64]

# Windows (cross-build, requires mingw-w64)  <-- Not yet included in GitHub Actions. The windows build process is completely untested
bash scripts/build_windows.sh  vX.Y.Z
```

Each script creates:

```
libs/<os_arch>/<version>/libraw_wrapper.a
```

And updates the symlink:

```
libs/<os_arch>/current → <version>
```

This structure ensures consistent versioning and `cgo` compatibility across builds.
