package glass

const FRAME_TEST = "\x00\x00"

type TestFrame struct {
    Id [16]byte
    From [16]byte
    To [16]byte
    Content []byte
}

func (*TestFrame) Bytes() []byte {
    return []byte(FRAME_TEST + "14bytes string" + "16-bytes-string!")
}

func (frame *TestFrame) Read(bs []byte) bool {
    frame.Content = []byte{}
    return true
}

func (frame *TestFrame) Type() string {
    return FRAME_TEST
}
