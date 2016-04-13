package glass

type Frame interface {
    Read([]byte) bool
    Bytes() []byte
}
