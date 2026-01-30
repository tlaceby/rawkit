package core

// LibrawError represents an error code returned by LibRaw.
// It implements the [error] interface.
type LibrawError int

const (
	LibrawSuccess               LibrawError = 0       // Pretty self explainatory
	LibrawUnspecifiedError      LibrawError = -1      // Unspecified internal error
	LibrawFileUnsupported       LibrawError = -2      // File format not supported
	LibrawRequestForNonexistent LibrawError = -3      // Requested item does not exist
	LibrawOutOfOrderCall        LibrawError = -4      // API called in wrong order
	LibrawNoThumbnail           LibrawError = -5      // No embedded thumbnail found
	LibrawInsufficientMemory    LibrawError = -100007 // Memory allocation failed
	LibrawDataError             LibrawError = -100008 // Data format or corruption error
	LibrawIOError               LibrawError = -100009 // File I/O error
	LibrawCancelledByCallback   LibrawError = -100010 // Processing cancelled by callback
	LibrawBadCrop               LibrawError = -100011 // Invalid crop parameters
	LibrawTooBig                LibrawError = -100012 // Image dimensions too large
	LibrawMempoolOverflow       LibrawError = -100013 // Internal memory pool overflow
)

// Error returns a human-readable description of the LibRaw error.
func (e LibrawError) Error() string {
	switch e {
	case LibrawSuccess:
		return "success"
	case LibrawUnspecifiedError:
		return "unspecified error"
	case LibrawFileUnsupported:
		return "file unsupported"
	case LibrawRequestForNonexistent:
		return "request for nonexistent image"
	case LibrawOutOfOrderCall:
		return "out of order call"
	case LibrawNoThumbnail:
		return "no thumbnail"
	case LibrawInsufficientMemory:
		return "insufficient memory"
	case LibrawDataError:
		return "data error"
	case LibrawIOError:
		return "IO error"
	case LibrawCancelledByCallback:
		return "cancelled by callback"
	case LibrawBadCrop:
		return "bad crop"
	case LibrawTooBig:
		return "too big"
	case LibrawMempoolOverflow:
		return "mempool overflow"
	default:
		return "unknown error"
	}
}
