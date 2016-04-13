package glass

import (
    "net"
    "log"
)

type App struct {
    // Connection Handlers (future...)
    IncomingConnection func(net.Conn)
    // Frames Handlers
    ProcessSimpleFrame func(*SimpleFrame)
    ProcessTestFrame func(*TestFrame)
}

func (*App) Init(auth AuthToken) { }

func (*App) Dial(Peer) (net.Conn, error) {
    return nil, nil
}
func (*App) Send(Frame) error {
    return nil
}

// Decrypt the packet and give it to the app
func (app *App) processFrame(payload []byte) {
    // Find out frame type
    magic := string(payload[0:2])
    switch magic {
    case FRAME_SIMPLE:
        p := new(SimpleFrame)
        p.Read(payload)
        app.ProcessSimpleFrame(p)
    case FRAME_TEST:
        log.Println("[  ] Received test packet, processing...")
        p := new(TestFrame)
        p.Read(payload)
        app.ProcessTestFrame(p)
    default:
        // If the message type is not supported,
        // just ignore it. Does this lead to problems?
        log.Println("[  ] Unknown incoming packet type: ignoring.")
    }
}

func (app *App) In(conn net.Conn) (fs *FrameStream) {
    if err := fs.Out(app, conn); err != nil {
        return nil
    }
    return
}
func (app *App) Out(conn net.Conn) (fs *FrameStream) {
    if err := fs.Out(app, conn); err != nil {
        return nil
    }
    return
}
