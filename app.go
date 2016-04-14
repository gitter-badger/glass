package glass

import (
    "net"
    "log"
    "errors"
)

type App struct {
    Token AuthToken
    stream FrameStream
    // Exit channel
    shouldClose chan bool
    // Connection State callback
    //ConnState func(net.Conn, ConnState)
    // Connection Handlers (future...)
    IncomingConnection func(Peer, net.Conn)
    // Frames Handlers
    ProcessSimpleFrame func(*SimpleFrame)
    ProcessTestFrame func(*TestFrame)
}

func (app *App) init() {
    app.shouldClose = make(chan bool)
}

func (app *App) Connect() (stream *FrameStream, err error) {
    if app.Token.Router == nil {
        return nil, errors.New("Can't connect app: no router specified.")
    }
    var conn net.Conn
    conn, err = net.Dial("tcp", app.Token.Router.Addr)
    if err != nil { return }
    stream = &FrameStream{
        Conn: conn,
        Direction: STREAM_OUT,
        FrameHandler: app.frameHandler,
    }
    if err = stream.Handshake(); err != nil { return }
    app.stream = *stream
    return
}

func (app *App) ListenAndServe() error {
    if app.Token.Me == nil {
        return errors.New("Can't connect app: no self specified.")
    }
    l, err := net.Listen("tcp", app.Token.Me.Addr)
    if err != nil {
        return err
    }
    defer l.Close()
    for {
        conn, err := l.Accept()
        if err != nil { continue }
        // FIXME multiple connections
        var stream = FrameStream{
            Conn: conn,
            Direction: STREAM_IN,
            FrameHandler: app.frameHandler,
        }
        if err = stream.Handshake(); err != nil { continue }
        app.stream = stream
        go stream.Serve()
        break // FIXME it now accepts only one stream!
    }
    return nil
}

func (app *App) Close() (err error) {
    log.Println("Close called")
    err = app.stream.Close()
    if app.shouldClose != nil {
        close(app.shouldClose)
        //app.shouldClose = nil
    }
    return
}

func (app *App) Block() {
    if app.shouldClose == nil { app.init() }
    <-app.shouldClose
}

func (*App) Dial(Peer) (net.Conn, error) {
    return nil, nil
}
func (*App) Send(Frame) error {
    return nil
}

// Decrypt the packet and give it to the app
func (app *App) frameHandler(payload []byte) {
    // Find out frame type
    magic := string(payload[0:2])
    switch magic {
    case FRAME_SIMPLE:
        p := new(SimpleFrame)
        p.Read(payload)
        if app.ProcessSimpleFrame != nil {
            app.ProcessSimpleFrame(p)
        }
    case FRAME_TEST:
        log.Println("[  ] Received test packet, processing...")
        p := new(TestFrame)
        p.Read(payload)
        if app.ProcessTestFrame != nil {
            app.ProcessTestFrame(p)
        }
    default:
        // If the message type is not supported,
        // just ignore it. Does this lead to problems?
        log.Println("[  ] Unknown incoming packet type: ignoring.")
    }
}

/*func (app *App) In(conn net.Conn) (*FrameStream, error) {
    stream := &FrameStream{App: app}
    if err := stream.In(conn); err != nil {
        return nil, err
    }
    return stream, nil
}

func (app *App) Out(conn net.Conn) (*FrameStream, error) {
    stream := &FrameStream{App: app}
    if err := stream.Out(conn); err != nil {
        return nil, err
    }
    return stream, nil
}*/
