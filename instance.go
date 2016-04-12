package glassbox

type Instance interface {
    ProcessSimplePacket(p *SimplePacket)
}
