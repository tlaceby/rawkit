package rawkit_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

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
	if img.Image.Width != wantW || img.Image.Height != wantH {
		t.Fatalf("unexpected size %dx%d (want %dx%d)", img.Image.Width, img.Image.Height, wantW, wantH)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err = img.Image.CreateThumbnail(dst); err != nil {
		t.Fatalf("could not create thumbnail: %s", err.Error())
	}

	os.Remove(dst)
}
