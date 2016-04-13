package glass

import (
    "net"
)

type App struct {
    Handler Handler
}

func (App) Init(auth string) {
    //
}

func (App) Dial(Peer) (net.Conn, error) {
    return nil, nil
}
func (App) Send(Frame) error {
    return nil
}


type Handler interface {
    // Connection Handler
    IncomingConnection(net.Conn)
    // Frames Handlers
    ProcessSimpleFrame(*SimpleFrame)
    ProcessTestFrame(*TestFrame)
}
