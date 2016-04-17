package router

import (
    //"net"
    //"time"
    //"crypto/tls"
    "github.com/acondolu/glass"
)


type Router struct {
    glass.App
    streams map[[16]byte]*FrameStream
}

/* FIXME!!!
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
        var stream *glass.FrameStream
        if stream = r.In(conn); stream == nil {
            stream.Shutdown(nil) // FIXME
            return // FIXME
        }
        go stream.Serve()
	}
    return nil
}
*/

type frame struct {
  To [16]byte
  Datagram []byte
}



func (this *Router) route(f *frame) {
  stream := this.streams[f.To]
  if stream != nil {
    // Handle the datagram routing
    stream.Queue(f.Datagram)
  } else {
    // No persistency for now, discard...
    // FIXME in Alpha-2
  }
}
