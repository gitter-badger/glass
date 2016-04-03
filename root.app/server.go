package main

import (
    "net"
    "time"
    //"errors"
    "bytes"
    "strings"
    "encoding/binary"
)

// First message sent from client
const CLIENT_HELLO = "01234567"
// Sent in response to a client hello
const SERVER_HELLO = "76543210"

type Server struct {

}

type Stream struct {
    serv Server
    conn net.Conn
}

func (s Stream) Init(serv Server, conn net.Conn) bool {
    var b [8]byte
    // Hello phase
    _, err := conn.Read(b[:])
    // Initial stream
    if err != nil || strings.Compare(string(b[:]), CLIENT_HELLO) != 0 {
        s.Shutdown()
        return false
    }
    // Response stream
    _, err = conn.Write([]byte(SERVER_HELLO))
    if err != nil {
        s.Shutdown()
        return false
    }
    s.conn = conn
    s.serv = serv
    return true
}

func (s Stream) Write(p Packet) error {
    _, err := s.conn.Write(p.Bytes())
    return err
}

func (s Stream) Shutdown() error {
    return s.conn.Close()
}

func (s Stream) Serve() {
    word := make([]byte, 8)
    conn := s.conn
    var magic, head, body []byte
    var buf *bytes.Reader
    var headSize uint8
    var bodySize uint16
    for {
        _, err := conn.Read(word)
        if err != nil {
            panic(err)
        }
        if word[0] != '\xff' {
            panic(err)
        }
        magic = word[1:5]
        // TODO Check if supported
        buf = bytes.NewReader(word[5:5])
        err = binary.Read(buf, binary.LittleEndian, &headSize)
        buf = bytes.NewReader(word[6:7])
        err = binary.Read(buf, binary.LittleEndian, &bodySize)
		head = make([]byte, int(headSize))
		_, err = conn.Read(head)
		if err != nil {
			panic(err)
		}
        body = make([]byte, int(bodySize))
        _, err = conn.Read(body)
		if err != nil {
			panic(err)
		}
        go s.serv.Process(magic, head, body)
	}
}

// TODO
func (s Server) Process(magic, head, body []byte) {}

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
        stream := new(Stream)
        if stream.Init(s, conn) {
            go stream.Serve()
        }
	}
}
