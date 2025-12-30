package core

// Creates a RAW processor which uses Libraw processing to construct an Image
func ProcessImage(filePath string) (*ImageProcessor, error) {
	// Check for file existing
	// Check if openable by Libraw
	// 		Otherwise, try to load non-raw format (PNG, JPEG)
	// Store handle for use by other
	return nil, nil
}

// Whether the file that the processor is working on is actually RAW
func (p *ImageProcessor) IsRawFile() bool {
	return false
}

// Performs important memory cleanup operations. Once called, this instance can no longer be used
func (p *ImageProcessor) Cleanup() error {
	return nil
}

// Perform extraction and optionally extract thumbnail data too. Will extract MetaData and ImageData by default
func (p *ImageProcessor) Extract(extractThumbnail bool) (*Image, error) {
	return nil, nil
}

// Extract the image thumbnail data
func (p *ImageProcessor) ExtractThumbnail() (*Image, error) {
	return nil, nil
}

// Extract the full image data. This operation can take a while depending on size of image
func (p *ImageProcessor) ExtractImageData() (*Image, error) {
	return nil, nil
}

// Extract the metadata if possible. If the file is not a RAW file, this will return an error
func (p *ImageProcessor) ExtractMetadata() (*Image, error) {
	return nil, nil
}
