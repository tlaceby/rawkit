package core

// IsRaw reports whether the image is a RAW format.
func (i *Image) IsRaw() bool {
	return i.Type == IMG_TYPE_RAW
}

// Index returns the buffer index for the pixel at coordinates (x, y).
// The returned index points to the R channel; G and B follow at idx+1 and idx+2.
func (img *ImageData) Index(x, y int) int {
	return (y*img.Width + x) * int(img.Channels)
}

// Pixel returns the RGB values for the pixel at coordinates (x, y).
// Coordinates are zero-indexed from the top-left corner.
func (img *ImageData) Pixel(x, y int) (r, g, b uint16) {
	idx := img.Index(x, y)
	return img.Data[idx], img.Data[idx+1], img.Data[idx+2]
}
