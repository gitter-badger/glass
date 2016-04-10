package router

import (
    "net"
    "time"
    "crypto/tls"
    "github.com/acondolu/glassbox"
)


type Router struct {}

// TODO Process this packet
func (r Router) Process(orig *glassbox.PacketStream, data []byte) {
    // TODO
}

// TODO Check if this magic is supported
func (r Router) IsSupportedMagic(magic [4]byte) bool {
    return false;
}

func (r Router) Start(cert tls.Certificate) (err error) {
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
