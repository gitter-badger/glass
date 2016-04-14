package router

import (
    //"net"
    //"time"
    //"crypto/tls"
    "github.com/acondolu/glass"
)


type Router struct {
    glass.App
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
