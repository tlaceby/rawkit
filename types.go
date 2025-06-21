package rawkit

// LibrawImageFormat represents the final output format of the decoded RAW image. Derived from the `enum LibRaw_image_formats`
type LibrawImageFormat string

const (
	// LIBRAW_IMAGE_JPEG specifies output in standard JPEG format.
	LIBRAW_IMAGE_JPEG LibrawImageFormat = "image/jpeg"

	// LIBRAW_IMAGE_BITMAP specifies output in uncompressed bitmap format.
	LIBRAW_IMAGE_BITMAP LibrawImageFormat = "image/bitmap"

	// LIBRAW_IMAGE_JPEGXL specifies output in JPEG XL format (modern, lossless/lossy hybrid).
	LIBRAW_IMAGE_JPEGXL LibrawImageFormat = "image/jpegxl"

	// LIBRAW_IMAGE_H265 specifies output in H.265 video format (experimental).
	LIBRAW_IMAGE_H265 LibrawImageFormat = "image/h265"
)

// ColorSpace describes the working color space of the image.
type ColorSpace string

const (
	// LIBRAW_COLORSPACE_sRGB represents the standard sRGB color space.
	LIBRAW_COLORSPACE_sRGB ColorSpace = "sRGB"

	// LIBRAW_COLORSPACE_ADOBE represents the Adobe RGB color space.
	LIBRAW_COLORSPACE_ADOBE ColorSpace = "Adobe"

	// LIBRAW_COLORSPACE_UNKNOWN is used when the color space cannot be determined.
	LIBRAW_COLORSPACE_UNKNOWN ColorSpace = "Unknown"
)

// DriveMode represents the camera’s shooting mode (single, continuous low/high).
type DriveMode string

const (
	// LIBRAW_DRIVEMODE_SINGLE_FRAME indicates the camera was in single-shot mode.
	LIBRAW_DRIVEMODE_SINGLE_FRAME DriveMode = "SingleFrame"

	// LIBRAW_DRIVEMODE_CONTINUOUS_LOW indicates a low-speed continuous burst mode.
	LIBRAW_DRIVEMODE_CONTINUOUS_LOW DriveMode = "ContinuousLow"

	// LIBRAW_DRIVEMODE_CONTINUOUS_HIGH indicates a high-speed continuous burst mode.
	LIBRAW_DRIVEMODE_CONTINUOUS_HIGH DriveMode = "ContinuousHigh"
)

// Represents all LibRaw image data after processing. Contains image pixel data, camera and lens info, and colorspace data.
type RawKitImage struct {
	// LibRaw image format (jpeg, jpegxl, bitmap, h265)
	Format LibrawImageFormat
	// stored as packed pixels (16-bit)
	Buffer []uint16
	// Size of visible ("meaningful") part of the image (without the frame).
	Width int
	// Size of visible ("meaningful") part of the image (without the frame).
	Height int
	// (1=CFA/RAW - 3=RGB)
	Colors int
	// Image orientation (0 if does not require rotation; 3 if requires 180-deg rotation; 5 if 90 deg counterclockwise, 6 if 90 deg clockwise).
	Flip int
	// Set to 1 if WB already applied in camera (multishot modes; small raw)
	AsShotWBApplied int
	// optional, depending on raw file depth
	RawBitsPerPixel uint
	// Number of Raw Images in RAW file (>0)
	RawCount uint
	// DNG version (for the DNG format).
	DNGVersion uint
	// What color space the image was encoded in (sRGB, Adobe, Unknown)
	ColorSpace ColorSpace
	// Camera ISO value
	ISO float32
	// Camera exposure shutter speed
	ShutterSpeed float32
	// Camera lens aperature
	Aperature float32
	// Lens focal length
	FocalLength float32
	// Artist information if provided by camera
	Artist string

	// Camera Make (eg: Sony)
	CameraMake string
	// Camer Model (eg: ILCE-6400HHILCE-6400 vN.xxxxxx:xx:xx xxx:xxx:xx)
	CameraModel string
	// (eg: Sony)
	CameraNormalizedMake string
	// eg: (ILCE-6400)
	CameraNormalizedModel string
	// Camera Software (eg: ILCE-6400 v2.002025:02:09 19:19:12 xxxxxxxxxx xxxxxxxxxxxx)
	CameraSoftware string
}
