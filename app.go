package glass

import (
    "net"
    "log"
    "errors"
    "fmt"
)

type App struct {
    Token AuthToken

    // Connection storage
    listeners map[*net.Listener]bool
    streams map[*FrameStream]bool

    /*
    MaxWorkers int
    MaxWorkersFrames int
    MaxWorkersFrames int
	MaxQueueFrames int
    */

    didInit bool
    // Exit channel
    shouldClose chan bool
    // Connection State callback
    // ConnState func(net.Conn, ConnState)
    // Connection Handlers
    IncomingConnection func(Peer, net.Conn)
    // Frames Handlers
    ProcessSimpleFrame func(*SimpleFrame)
    ProcessTestFrame func(*TestFrame)
}

func (app *App) init() {
    app.shouldClose = make(chan bool)
    app.streams = make(map[*FrameStream]bool)
    app.listeners = make(map[*net.Listener]bool)
    app.didInit = true
}

func (app *App) Connect() (stream *FrameStream, err error) {
    if !app.didInit { app.init() }
    if app.Token.Router == nil {
        return nil, errors.New("Can't connect app: no router specified.")
    }
    var conn net.Conn
    conn, err = net.Dial("tcp", app.Token.Router.Addr)
    if err != nil { return }
    fmt.Println("Dialed")
    stream = &FrameStream{
        Conn: conn,
        Direction: STREAM_OUT,
        FrameHandler: app.frameHandler,
    }
    if err = stream.Handshake(); err != nil {
        fmt.Printf("[Client] Handshake error: %e.\n", err)
        log.Fatal(err)
        return
    }
    fmt.Println("Dialed. handshake done.")
    app.streams[stream] = true
    return
}

func (app *App) ListenAndServe() (err error) {
    if !app.didInit { app.init() }
    if app.Token.Me == nil {
        return errors.New("Can't connect app: no self specified.")
    }
    var lstn net.Listener
    lstn, err = net.Listen("tcp", app.Token.Me.Addr)
    if err != nil { return err }
    //defer lstn.Close()
    var conn net.Conn
    for {
        conn, err = lstn.Accept()
        if err != nil { continue }
        var stream = &FrameStream{
            Conn: conn,
            Direction: STREAM_IN,
            FrameHandler: app.frameHandler,
        }
        if err = stream.Handshake(); err != nil {
            log.Fatal(err)
        }
        app.streams[stream] = true
        go stream.Serve()
        break // FIXME it now accepts only one stream!
    }
    return nil
}

func (app *App) Close() {
    log.Println("Close called")
    if !app.didInit { return }
    // Close all open streams
    for s := range app.streams {
        if s.Close() == nil {
            delete(app.streams, s)
        }
    }
    // Close all listeners
    for ln, _ := range app.listeners {
        if (*ln).Close() == nil {
            delete(app.listeners, ln)
        }
    }

    if app.shouldClose != nil {
        close(app.shouldClose)
        //app.shouldClose = nil
    }
    return
}

func (app *App) Block() {
    if !app.didInit { return }
    <-app.shouldClose
}

func (*App) Dial(Peer) (net.Conn, error) {
    return nil, nil
}
func (*App) Send(Frame) error {
    return nil
}

// Decrypt the packet and give it to the app
func (app *App) frameHandler(typ string, payload []byte) {
    switch typ {
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
