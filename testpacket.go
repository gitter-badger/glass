package glassbox

const PACKET_TEST = "\x00\x00"

type TestPacket struct {}



func (p *TestPacket) Id() [16]byte {
    var x [16]byte; return x
}
func (p *TestPacket) From() [16]byte {
    var x [16]byte; return x
}
func (p *TestPacket) To() [16]byte {
    var x [16]byte; return x
}
func (p *TestPacket) Content() []byte {
    return []byte("16-bytes-string!")
}
func (p *TestPacket) Bytes() []byte {
    return []byte(PACKET_TEST + "14bytes string" + "16-bytes-string!")
}

func (p *TestPacket) FromBytes(bs []byte) {
    // TODO
}
