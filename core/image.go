package core

// Helper methods for the Image, ImageData & ImageMeta structs

// Whether the image is stored in RAW format
func (i *Image) IsRaw() bool { return i.Type == IMG_TYPE_RAW }

// Given an (x,y) return the proper index into the buffer
func (img *ImageData) Index(x, y int) int {
	return (y*img.Width + x) * img.Channels
}
