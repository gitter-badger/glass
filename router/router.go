package router

import (
    "net"
    "time"
    "crypto/tls"
    "github.com/acondolu/glassbox"
)


type Router struct {}
func (*Router) Init(string) {}
func (*Router) Dial(glassbox.Entity) (net.Conn, error) { return nil, nil }
func (*Router) Send(glassbox.Packet) error { return nil }
func (*Router) IncomingConnection(net.Conn) {}
func (*Router) ProcessSimplePacket(*glassbox.SimplePacket) { }
func (*Router) ProcessTestPacket(*glassbox.TestPacket) { }

func (r *Router) Start(cert *tls.Certificate) (err error) {
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
        stream := new(glassbox.PacketStream)
        if err = stream.In(r, conn, cert); err != nil {
            stream.Shutdown(err)
            return // FIXME
        }
        go stream.Serve()
	}
    return nil
}
