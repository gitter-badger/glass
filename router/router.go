package router

import (
    "net"
    "time"
    //"crypto/tls"
    "github.com/acondolu/glass"
)


type Router struct {}
func (*Router) Init(string) {}
func (*Router) Dial(glass.Peer) (net.Conn, error) { return nil, nil }
func (*Router) Send(glass.Frame) error { return nil }
func (*Router) IncomingConnection(net.Conn) {}
func (*Router) ProcessSimpleFrame(*glass.SimpleFrame) {}
func (*Router) ProcessTestFrame(*glass.TestFrame) {}

func (r *Router) Serve() (err error) {
    var ln net.Listener
    var conn net.Conn
	ln, err = net.Listen("tcp", ":8081")
	if err != nil {
		return
	}
	defer ln.Close()
	for {
		conn, err = ln.Accept()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
        stream := new(glass.FrameStream)
        if err = stream.In(r, conn); err != nil {
            stream.Shutdown(err)
            return // FIXME
        }
        go stream.Serve()
	}
    return nil
}
