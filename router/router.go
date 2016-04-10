package router

import (
    "net"
    "time"
    "github.com/acondolu/glassbox"
)


type Router struct {}

// TODO Process this packet
func (r Router) Process(orig glassbox.PacketStream, data []byte) {
    // TODO
}

// TODO Check if this magic is supported
func (r Router) IsSupportedMagic(magic [4]byte) bool {
    return false;
}

func (r Router) StartServer() (err error) {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		return
	}
	defer ln.Close()
	for {
		_, err := ln.Accept()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
        stream := new(glassbox.PacketStream)
        //if err = stream.Init(glassbox.STREAM_IN, s, conn, ...); err != nil {
        //    stream.Shutdown()
        //    return
        //}
        go stream.Serve()
	}
}
