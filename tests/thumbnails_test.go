package rawkit_test

import (
	"image"
	"os"
	"testing"

	"github.com/tlaceby/rawkit"
)

func TestThumbnailTiffPath(t *testing.T) {
	const src = "testdata/tree-large.ARW"
	const dst = "testdata/tmp/tree-tiff.tif"

	img, err := rawkit.LoadRAW(src)
	if err != nil {
		t.Fatalf("rawkit.LoadRAW failed: %v", err)
	}

	if err = img.Image.CreateThumbnail(dst); err != nil {
		t.Fatalf("failed to generate tif thumbnail at path %s wioth type .tif-> error: %s", dst, err.Error())
		os.Remove(dst)
	}

	os.Remove(dst)
}

func TestThumbnailNoPath(t *testing.T) {
	const src = "testdata/tree-large.ARW"
	const dst = "testdata/tmp/tree-jpeg"

	img, err := rawkit.LoadRAW(src)
	if err != nil {
		t.Fatalf("rawkit.LoadRAW failed: %v", err)
	}

	if err = img.Image.CreateThumbnail(dst); err != nil {
		t.Fatalf("failed to generate tif thumbnail at path %s -> error: %s", dst, err.Error())
		os.Remove(dst)
	}

	os.Remove(dst + ".jpg")
}

func TestThumbnailPngPath(t *testing.T) {
	const src = "testdata/tree-large.ARW"
	const dst = "testdata/tmp/tree-png.png"
	img, err := rawkit.LoadRAW(src)
	if err != nil {
		t.Fatalf("rawkit.LoadRAW failed: %v", err)
	}

	if err = img.Image.CreateThumbnail(dst); err != nil {
		t.Fatalf("thumbnail options test failed: path=%s and error=%s", dst, err.Error())
		os.Remove(dst)
	}

	os.Remove(dst)
}

func TestThumbnailInvalidPath(t *testing.T) {
	const src = "testdata/tree-large.ARW"
	const dst = "testdata/tmp/tree-bad.foo"
	img, err := rawkit.LoadRAW(src)
	if err != nil {
		t.Fatalf("rawkit.LoadRAW failed: %v", err)
	}

	if err = img.Image.CreateThumbnail(dst); err == nil {
		t.Fatalf("thumbnail options test failed: path=%s", dst)
		os.Remove(dst)
	}

	os.Remove(dst)
}

func TestThumbnailJpgPath(t *testing.T) {
	const src = "testdata/tree-large.ARW"
	const dst = "testdata/tmp/tree-jpg.jpg"
	img, err := rawkit.LoadRAW(src)
	if err != nil {
		t.Fatalf("rawkit.LoadRAW failed: %v", err)
	}

	if err = img.Image.CreateThumbnail(dst); err != nil {
		t.Fatalf("thumbnail options test failed: path=%s and error=%s", dst, err.Error())
		os.Remove(dst)
	}

	os.Remove(dst)
}

func TestThumbnailWithInvalidOutputPath(t *testing.T) {
	const src = "testdata/tree-large.ARW"
	const dst = "testdata/tmp/does-not-exist/thumbnail.jpg"

	img, err := rawkit.LoadRAW(src)
	if err != nil {
		t.Fatalf("rawkit.LoadRAW failed: %v", err)
	}

	if err = img.Image.CreateThumbnail(dst); err == nil {
		t.Fatalf("thumbnail options test failed: passed path of %s and kind .tiff", dst)
		os.Remove(dst)
	}

	os.Remove(dst)
}

func TestThumbnailWithInvalidCrop(t *testing.T) {
	const src = "testdata/tree-large.ARW"
	const dst = "testdata/tmp/invalid-crop.jpg"

	img, err := rawkit.LoadRAW(src)
	if err != nil {
		t.Fatalf("rawkit.LoadRAW failed: %v", err)
	}

	if err = img.Image.CreateThumbnail(dst, rawkit.WithThumbnailCrop(image.Rect(0, 0, img.Image.Width+5, img.Image.Height+5))); err == nil {
		t.Fatalf("thumbnail options test failed: passed path of %s and kind .tiff", dst)
		os.Remove(dst)
	}
}

func TestThumbnailWithValidCrop(t *testing.T) {
	const src = "testdata/tree-large.ARW"
	const dst = "testdata/tmp/cropped-thumbnail.png"

	img, err := rawkit.LoadRAW(src)
	if err != nil {
		t.Fatalf("rawkit.LoadRAW failed: %v", err)
	}

	if err = img.Image.CreateThumbnail(dst, rawkit.WithThumbnailCrop(image.Rect(500, 500, img.Image.Width-1000, img.Image.Width-2000))); err != nil {
		t.Fatalf("thumbnail crop option test failed. This crop should pass for this image --> error: %s", err.Error())
		os.Remove(dst)
	}

	os.Remove(dst)
}
