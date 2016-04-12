package glassbox

import (
    "net"
)

type App interface {
    Init(auth string)

    Dial(Entity) (net.Conn, error)
    Send(Packet) error

    // Connection Handler
    IncomingConnection(net.Conn)
    // Packet Handlers
    ProcessSimplePacket(*SimplePacket)
    ProcessTestPacket(*TestPacket)
}
