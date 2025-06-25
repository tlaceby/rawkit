package rawkit

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"slices"
	"strings"

	"github.com/disintegration/imaging"
)

// Creates a Go standard image.NRGBA from the *RawKitImage.Image
func (rk *RawKitImage) RGBAImage() *image.NRGBA {
	return rk.Image.RGBAImage()
}

// Creates a go standard *image.NRGBA from a rawkit.Image
func (img *Image) RGBAImage() *image.NRGBA {
	rgba := image.NewNRGBA(image.Rect(0, 0, img.Width, img.Height))

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			offset := (y*img.Width + x) * img.Colors
			r := uint8(img.Buffer[offset+0] >> 8)
			g := uint8(img.Buffer[offset+1] >> 8)
			b := uint8(img.Buffer[offset+2] >> 8)
			rgba.Set(x, y, color.NRGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return rgba
}

// CreateThumbnail creates a thumbnail from a rawkit.Image.
// It accepts options like crop, quality, and kind using helper functions.
// Returns an error if the thumbnail could not be created or if any parameters are invalid.
func (img *Image) CreateThumbnail(path string, opts ...ThumbnailOpt) error {
	conf := &ThumbnailOptions{
		Quality: 90,
		Crop:    image.Rect(0, 0, img.Width, img.Height),
	}

	for _, opt := range opts {
		opt(conf)
	}

	if conf.Quality < 1 || conf.Quality > 100 {
		return fmt.Errorf("Expected thumbnail.quality between 1 and 100. Recieved %d instead", conf.Quality)
	}

	if !conf.Crop.In(image.Rect(0, 0, img.Width, img.Height)) {
		return fmt.Errorf("invalid crop specified (x0=%d, y0=%d, x1=%d, y1=%d). Crop is not contained in the Image", conf.Crop.Min.X, conf.Crop.Min.Y, conf.Crop.Max.X, conf.Crop.Max.Y)
	}

	lowerExt := strings.ToLower(filepath.Ext(path))
	if len(lowerExt) == 0 {
		path += ".jpg"
	} else if !isValidExtension(lowerExt) {
		return fmt.Errorf("invalid thumbnail extension %s provided", lowerExt)
	}

	return imaging.Save(img.RGBAImage(), path, imaging.JPEGQuality(conf.Quality))
}

// WithThumbnailQuality sets the JPEG/PNG quality (1–100) for CreateThumbnail.
func WithThumbnailQuality(q int) ThumbnailOpt {
	return func(to *ThumbnailOptions) { to.Quality = q }
}

// WithThumbnailCrop sets the crop rectangle for CreateThumbnail.
// The crop must be fully contained within the image, or CreateThumbnail will return an error.
func WithThumbnailCrop(c image.Rectangle) ThumbnailOpt {
	return func(to *ThumbnailOptions) { to.Crop = c }
}

// -------------------------------=--------------------------------------------
// ------------------------------ HELPER METHODS ------------------------------
// -------------------------------=--------------------------------------------

func isValidExtension(ext string) bool {
	exts := []string{".png", ".jpeg", ".jpg", ".tiff", ".tif"}
	return slices.Contains(exts, ext)
}
