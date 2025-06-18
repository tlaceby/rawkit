package rawkit

type RawKitImage struct {
	Buffer   []uint16 // stored as packed pixels (16-bit)
	Width    int
	Height   int
	Channels int // 1 = CFA/raw, 3 = RGB, etc.
}
