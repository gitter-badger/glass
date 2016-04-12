package glassbox

type Instance interface {
    ProcessSimplePacket(p *SimplePacket)
    // Test Packet
    ProcessTestPacket(p *TestPacket)
}
