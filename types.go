package rawkit

type LibrawImageFormat string
type ColorSpace string
type DriveMode string

const (
	LIBRAW_IMAGE_JPEG   LibrawImageFormat = "image/jpeg"
	LIBRAW_IMAGE_BITMAP LibrawImageFormat = "image/bitmap"
	LIBRAW_IMAGE_JPEGXL LibrawImageFormat = "image/jpegxl"
	LIBRAW_IMAGE_H265   LibrawImageFormat = "image/h265"
)

const (
	LIBRAW_COLORSPACE_sRGB    ColorSpace = "sRGB"
	LIBRAW_COLORSPACE_ADOBE   ColorSpace = "Adobe"
	LIBRAW_COLORSPACE_UNKNOWN ColorSpace = "Unknown"
)

const (
	LIBRAW_DRIVEMODE_SINGLE_FRAME    DriveMode = "SingleFrame"
	LIBRAW_DRIVEMODE_CONTINUOUS_LOW  DriveMode = "ContinuousLow"
	LIBRAW_DRIVEMODE_CONTINUOUS_HIGH DriveMode = "ContinuousHigh"
)

type RawKitImage struct {
	Format LibrawImageFormat
	Buffer []uint16 // stored as packed pixels (16-bit)
	Width  int      // Size of visible ("meaningful") part of the image (without the frame).
	Height int

	Colors int // (1=CFA/RAW -- 3=RGB)
	Flip   int // Image orientation (0 if does not require rotation; 3 if requires 180-deg rotation; 5 if 90 deg counterclockwise, 6 if 90 deg clockwise).

	AsShotWBApplied int  // Set to 1 if WB already applied in camera (multishot modes; small raw)
	RawBitsPerPixel uint // optional, depending on raw file depth
	RawCount        uint // raw
	DNGVersion      uint
	// ColorSpace      ColorSpace // TODO:

	ISO          float32
	ShutterSpeed float32
	Aperature    float32
	FocalLength  float32
	Artist       string

	CameraMake            string
	CameraModel           string
	CameraNormalizedMake  string
	CameraNormalizedModel string
	CameraSoftware        string
}
