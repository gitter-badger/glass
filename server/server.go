package server

import (
    "net"
    "time"
    "github.com/acondolu/glassbox"
)


type Server struct {}

// TODO Process this packet
func (s Server) Process(orig glassbox.Stream, magic [4]byte, head, body []byte) {

}

// TODO Check if this magic is supported
func (s Server) IsSupportedMagic(magic [4]byte) bool {
    return false;
}

func (s Server) StartServer() (err error) {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
        stream := new(glassbox.Stream)
        if stream.Init(s, conn) {
            go stream.Serve()
        }
	}
}
