package glass

const PACKET_TEST = "\x00\x00"

type TestFrame struct {}



func (*TestFrame) Id() [16]byte {
    return *new([16]byte)
}
func (*TestFrame) From() [16]byte {
    return *new([16]byte)
}
func (*TestFrame) To() [16]byte {
    return *new([16]byte)
}
func (*TestFrame) Content() []byte {
    return []byte("16-bytes-string!")
}
func (*TestFrame) Bytes() []byte {
    return []byte(PACKET_TEST + "14bytes string" + "16-bytes-string!")
}

func (*TestFrame) FromBytes(bs []byte) {
    // TODO
}
