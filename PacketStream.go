package glassbox

import (
    "net"
    "errors"
    "bytes"
    "encoding/binary"
    "crypto/rand"
    "crypto/tls"
    "crypto/aes"
    "crypto/cipher"
    "io"
)

// A stream can be incoming or outgoing
type StreamDirection int
const (
        STREAM_IN StreamDirection = iota
        STREAM_OUT
)

type Instance interface {
    //IsSupportedMagic(magic [4]byte) bool
    Process(orig *PacketStream, data []byte, iv [16]byte)
}

type PacketStream struct {
    inst Instance
    conn net.Conn

    secret [16]byte
    // FIXME TODO add here a counter for the exchanged blocks.
    // in GCM, 96 bit nonces can only produce 2^32 different blocks
    // (at 16 bytes per block)
}

func GenerateSecret(s1, s2, secret []byte) {
    // FIXME
    copy(secret, s1)
    copy(secret[8:], s2)
}

func (this *PacketStream) Init(
        dir StreamDirection,
        inst Instance,
        conn net.Conn,
        cert tls.Certificate,
    ) (err error) {
    var b [8]byte
    var b2 [8]byte
    var tlsConfig *tls.Config
    var tlsConn *tls.Conn
    this.conn = conn
    // constants
    // First message sent from client
    CLIENT_HELLO := []byte("01234567")
    // Sent in response to a client hello
    SERVER_HELLO := []byte("76543210")
    // Start hello phase
    if dir == STREAM_IN {
        // Phase 1: HELLO
        _, err = conn.Read(b[:])
        if err != nil { return }
        if !bytes.Equal(b[:], CLIENT_HELLO) { return nil } // FIXME
        if _, err = conn.Write(SERVER_HELLO); err != nil { return }
        // Phase 2: Negotiate TLS connection
        tlsConfig = &tls.Config{
            Certificates: []tls.Certificate{cert},
            ClientAuth: tls.VerifyClientCertIfGiven,
            ServerName: "example.com",
        }
        tlsConn = tls.Server(conn, tlsConfig)
        tlsConn.Handshake()
        // Phase 3: Agree on a secret, forget about TLS
        if _, err = conn.Read(b[:]); err != nil { return }
        //// Send my secret
        if _, err = io.ReadFull(rand.Reader, b2[:]); err != nil { return }
        if _, err = conn.Write(b2[:]); err != nil { return }
        // Generate AES128 key
        GenerateSecret(b[:], b2[:], this.secret[:])
    } else if dir == STREAM_OUT {
        // Phase 1: HELLO
        if _, err = conn.Write(SERVER_HELLO); err != nil { return }
        if _, err = conn.Read(b[:]); err != nil { return }
        if !bytes.Equal(b[:], CLIENT_HELLO) { return nil } // FIXME
        // Phase 2: Negotiate TLS connection
        tlsConfig = &tls.Config{
            Certificates: []tls.Certificate{cert},
            ServerName: "example.com",
        }
        tlsConn = tls.Client(conn, tlsConfig)
        tlsConn.Handshake()
        // Phase 3: Agree on a secret, forget about TLS
        //// Send my secret
        if _, err = io.ReadFull(rand.Reader, b2[:]); err != nil { return }
        if _, err = conn.Write(b2[:]); err != nil { return }
        //// Read client's secret
        if _, err = conn.Read(b[:]); err != nil { return }
        //// Generate AES128 key
        GenerateSecret(b2[:], b[:], this.secret[:])
    }
    this.inst = inst
    return nil
}

func (s *PacketStream) Write(p Packet) (err error) {
    // See: https://gist.github.com/kkirsche/e28da6754c39d5e7ea10
    data := p.Bytes()
    if len(data) % 16 != 0 {
        return errors.New("Wrong packet size")
    }
    // Encrypt data with AES128-GCM
    var block cipher.Block
    var aesgcm cipher.AEAD
    if block, err = aes.NewCipher(s.secret[:]); err != nil { return }
	nonce := make([]byte, 12) // 96-bits nonce
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil { return }
	if aesgcm, err = cipher.NewGCM(block); err != nil { return }
	data = aesgcm.Seal(nil, nonce, data, nil)

    conn := s.conn
    // 1) Send size
    size := len(data)
    if len(data) % 16 != 0 {
        return errors.New("Wrong packet size after GCM")
    }
    buf := new(bytes.Buffer)
    err = binary.Write(buf, binary.BigEndian, uint16(size / 16))
    if err != nil { return }
    two := make([]byte, 2)
    if _, err = buf.Read(two); err != nil { return }
    conn.Write(two)
    // 2) Send initialization vector
    //if _, err = io.ReadFull(rand.Reader, iv); err != nil { return }
    //conn.Write(iv)
    // 2) Send nonce
    conn.Write(nonce)
    // 3) Send AES128-encrypted data
    _, err = conn.Write(data)
    return
}

func (s *PacketStream) Shutdown(cause error) error {
    return s.conn.Close()
}

func (s *PacketStream) Serve() {
    var iv [16]byte
    var two [2]byte
    var data []byte
    var buf *bytes.Reader
    var length uint16
    var err error
    //var secret = s.secret
    conn := s.conn
    defer s.Shutdown(nil)
    for {
        // Read iv
        _, err = conn.Read(iv[:])
        if err != nil {
            s.Shutdown(err)
            break
        }
        // Read packet size
        _, err = conn.Read(two[:])
        if err != nil {
            s.Shutdown(err)
            break
        }
        // Convert length from bytes
        buf = bytes.NewReader(two[:])
        err = binary.Read(buf, binary.BigEndian, &length)
        // Should we close the connection?
        if length == 0 {
            break
        }
        // Read *length* blocks of data
        data = make([]byte, int(length) * 16)
        _, err = conn.Read(data)
        if err != nil {
            s.Shutdown(err)
            break
        }
        // Process packet
        go s.inst.Process(s, data, iv)
    }
}
