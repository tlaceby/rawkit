package rawkit

type RawKitImage struct {
	buffer   []uint16 // stored as packed pixels (16-bit)
	width    int
	height   int
	channels int // 1 = CFA/raw, 3 = RGB, etc.
}
