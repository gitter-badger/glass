package main

import (
    "net"
    "time"
    //"errors"
    "bytes"
    "strings"
    "encoding/binary"
)

// Protocol connection headers, both server-to-client and client-to-server
// FIXME naming
const PROTO_CONN_HEADER_C2S = "01234567"
const PROTO_CONN_HEADER_S2C = "76543210"

type Stream struct {
    conn net.Conn
}

func (s Stream) Init(conn  net.Conn) bool {
    var b [8]byte
    n, err := conn.Read(b[:])
    // Initial stream
    if n != 8 || err != nil || strings.Compare(string(b[:]), PROTO_CONN_HEADER_C2S) != 0 {
        s.Shutdown()
        return false
    }
    // Response stream
    n, err = conn.Write([]byte(PROTO_CONN_HEADER_S2C))
    if n != 8 || err != nil {
        s.Shutdown()
        return false
    }
    s.conn = conn
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
    var buf *bytes.Reader
    var headSize uint8
    var bodySize uint16
    for {
        _, err := conn.Read(word)
        if err != nil {
            panic(err)
        }
        if word[0] != '\xff' {
            continue
        }
        magic := word[1:5]
        buf = bytes.NewReader(word[5:5])
        err = binary.Read(buf, binary.LittleEndian, &headSize)
        buf = bytes.NewReader(word[6:7])
        err = binary.Read(buf, binary.LittleEndian, &bodySize)
		head := make([]byte, int(headSize))
		_, err = conn.Read(head)
		if err != nil {
			continue
		}
        body := make([]byte, int(bodySize))
        _, err = conn.Read(body)
		if err != nil {
			continue
		}
        go Process(magic, head, body)
	}
}

func Process(magic, head, body []byte) {

}

func StartServer() (err error) {
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
        s := new(Stream)
        if s.Init(conn) {
            go s.Serve()
        }
	}
}
