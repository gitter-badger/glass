package main

import (
    "net"
    //"errors"
    "bytes"
    "strings"
    "encoding/binary"
)

// First message sent from client
const CLIENT_HELLO = "01234567"
// Sent in response to a client hello
const SERVER_HELLO = "76543210"

// A stream can be incoming or outgoing
type StreamDirection int
const (
        STREAM_IN StreamDirection = iota
        STREAM_OUT
)

type Instance interface {
    IsSupportedMagic(magic [4]byte) bool
    Process(orig Stream, magic [4]byte, head, body []byte)
}

type Stream struct {
    inst Instance
    conn net.Conn
    key [8]byte
    dir StreamDirection
}

func (s Stream) Init(dir StreamDirection, inst Instance, conn net.Conn) bool {
    var b [8]byte
    // Hello phase
    _, err := conn.Read(b[:])
    // Initial stream
    if err != nil || string(b[:]) != CLIENT_HELLO {
        s.Shutdown()
        return false
    }
    // Response stream
    _, err = conn.Write([]byte(SERVER_HELLO))
    if err != nil {
        s.Shutdown()
        return false
    }
    // Initialize TLS? FIXME TODO
    s.dir = dir
    s.conn = conn
    s.inst = inst
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
    var head, body []byte
    var buf *bytes.Reader
    var headSize uint8
    var bodySize uint16
    for {
        // Read eight octets
        _, err := conn.Read(word)
        if err != nil {
            panic(err)
        }
        if word[0] != '\xff' {
            s.Shutdown()
            return
        }
        var magic [4]byte
        copy(magic[:], word[1:5])
        // Check first if this magic is supported
        if !s.inst.IsSupportedMagic(magic) {
            s.Shutdown()
            return
        }
        // Parse head size
        buf = bytes.NewReader(word[5:5])
        err = binary.Read(buf, binary.LittleEndian, &headSize)
        // Parse body size
        buf = bytes.NewReader(word[6:7])
        err = binary.Read(buf, binary.LittleEndian, &bodySize)
        // Read head
		head = make([]byte, int(headSize))
		_, err = conn.Read(head)
		if err != nil {
			panic(err)
		}
        // Read body
        body = make([]byte, int(bodySize))
        _, err = conn.Read(body)
		if err != nil {
			panic(err)
		}
        // Process packet
        go s.inst.Process(s, magic, head, body)
	}
}
