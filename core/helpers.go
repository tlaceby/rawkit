package core

import (
	"image"
	"image/color"
	"image/png"

	"github.com/disintegration/imaging"
)

// reports whether the image is a RAW format.
func (i *Image) IsRaw() bool {
	return i.Type == IMG_TYPE_RAW
}

// returns the buffer index for the pixel at coordinates (x, y).
// The returned index points to the R channel; G and B follow at idx+1 and idx+2.
func (img *ImageData) Index(x, y int) int {
	return (y*img.Width + x) * int(img.Channels)
}

// returns the RGB values for the pixel at coordinates (x, y).
// Coordinates are zero-indexed from the top-left corner.
func (img *ImageData) Pixel(x, y int) (r, g, b uint16) {
	idx := img.Index(x, y)
	return img.Data[idx], img.Data[idx+1], img.Data[idx+2]
}

// image.Image methods

// returns the color model for the image based on it's channels
func (img *ImageData) ColorModel() color.Model {
	if img.Channels == LIBRAW_CHANNELS_RGBA {
		return color.NRGBA64Model
	}

	return color.RGBA64Model
}

// returns the domain for which At can return non-zero color.
func (img *ImageData) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.Width, img.Height)
}

// returns the color of the pixel at (x, y).
func (img *ImageData) At(x, y int) color.Color {
	if x < 0 || x >= img.Width || y < 0 || y >= img.Height {
		return color.RGBA64{}
	}

	idx := img.Index(x, y)

	if img.Channels == LIBRAW_CHANNELS_RGBA {
		return color.NRGBA64{
			R: img.Data[idx],
			G: img.Data[idx+1],
			B: img.Data[idx+2],
			A: img.Data[idx+3],
		}
	}

	return color.RGBA64{
		R: img.Data[idx],
		G: img.Data[idx+1],
		B: img.Data[idx+2],
		A: 0xffff,
	}
}

type EncodeOption func(*encodeConfig)
type encodeConfig struct {
	jpegQuality       int
	pngCompression    png.CompressionLevel
	hasJPEGQuality    bool
	hasPNGCompression bool
}

// sets the JPEG encoding quality (1-100, where 100 is best).
func WithJPEGQuality(quality int) EncodeOption {
	return func(cfg *encodeConfig) {
		if quality < 1 {
			quality = 1
		}

		if quality > 100 {
			quality = 100
		}

		cfg.jpegQuality = quality
		cfg.hasJPEGQuality = true
	}
}

// sets the PNG compression level
func WithPNGCompression(level png.CompressionLevel) EncodeOption {
	return func(cfg *encodeConfig) {
		cfg.pngCompression = level
		cfg.hasPNGCompression = true
	}
}

const (
	PNGCompressionDefault = png.DefaultCompression
	PNGCompressionNone    = png.NoCompression
	PNGCompressionFast    = png.BestSpeed
	PNGCompressionBest    = png.BestCompression
)

// saves the image to file with the specified filename.
// The format is determined from the filename extension: "jpg" (or "jpeg"), "png", "gif", "tif" (or "tiff") and "bmp" are supported.
func (img *ImageData) Save(outputPath string, opts ...EncodeOption) error {
	cfg := &encodeConfig{
		jpegQuality:    90,
		pngCompression: png.DefaultCompression,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	var imagingOpts []imaging.EncodeOption
	if cfg.hasJPEGQuality {
		imagingOpts = append(imagingOpts, imaging.JPEGQuality(cfg.jpegQuality))
	}

	if cfg.hasPNGCompression {
		imagingOpts = append(imagingOpts, imaging.PNGCompressionLevel(cfg.pngCompression))
	}

	return imaging.Save(img, outputPath, imagingOpts...)
}

// constructs an ImageData from a standard image.Image. Colorspace defaults to sRGB and BitDepth to 16.
func ImageDataFromImage(img image.Image) *ImageData {
	bounds := img.Bounds()
	newWidth := bounds.Dx()
	newHeight := bounds.Dy()

	channels := LIBRAW_CHANNELS_RGB
	switch img.ColorModel() {
	case color.RGBAModel, color.RGBA64Model, color.NRGBAModel, color.NRGBA64Model, color.AlphaModel, color.Alpha16Model:
		channels = LIBRAW_CHANNELS_RGBA
	}

	newData := make([]uint16, newWidth*newHeight*int(channels))
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			idx := (y*newWidth + x) * int(channels)

			newData[idx] = uint16(r)
			newData[idx+1] = uint16(g)
			newData[idx+2] = uint16(b)

			if channels == LIBRAW_CHANNELS_RGBA {
				newData[idx+3] = uint16(a)
			}
		}
	}

	return &ImageData{
		Width:      newWidth,
		Height:     newHeight,
		Channels:   channels,
		Colorspace: LIBRAW_COLORSPACE_sRGB,
		Data:       newData,
	}
}
