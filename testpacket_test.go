package glassbox

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
    return make([]byte, 0)
}
