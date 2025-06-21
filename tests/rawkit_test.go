package rawkit_test

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/tlaceby/rawkit"
)

func TestLibRawVersionMatch(t *testing.T) {
	num := rawkit.LibRawVersionNum()
	str := rawkit.LibRawVersionStr()

	parts := strings.SplitN(str, ".", 3)
	if len(parts) < 3 {
		t.Fatalf("bad version string %q", str)
	}

	maj, _ := strconv.Atoi(parts[0])
	min, _ := strconv.Atoi(parts[1])

	// 0xMMmmpp  (M=major, m=minor, p=patch)
	gotMaj := num >> 16
	gotMin := (num >> 8) & 0xFF

	if gotMaj != maj || gotMin != min {
		t.Fatalf("mismatch: num=0x%x (%d.%d), str=%s", num, gotMaj, gotMin, str)
	}
}

func TestOpeningTreeLarge(t *testing.T) {
	const src = "testdata/tree-large.ARW"
	const dst = "testdata/tmp/tree.jpg"

	if _, err := os.Stat(src); err != nil {
		t.Fatalf("RAW file not found: %v", err)
	}
	img, err := rawkit.LoadRAW(src)
	if err != nil {
		t.Fatalf("rawkit.LoadRAW failed: %v", err)
	}

	wantW, wantH := 4024, 6024 // FIXME: Update with proper (cropped to sensor dimensions) (4000x6000)
	if img.Width != wantW || img.Height != wantH {
		t.Fatalf("unexpected size %dx%d (want %dx%d)", img.Width, img.Height, wantW, wantH)
	}

	rgba := image.NewNRGBA(image.Rect(0, 0, img.Width, img.Height))

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			i := (y*img.Width + x) * img.Colors
			r := uint8(img.Buffer[i+0] >> 8)
			g := uint8(img.Buffer[i+1] >> 8)
			b := uint8(img.Buffer[i+2] >> 8)
			rgba.Set(x, y, color.NRGBA{R: r, G: g, B: b, A: 255})
		}
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := imaging.Save(rgba, dst); err != nil {
		t.Fatalf("imaging.Save: %v", err)
	}
}
