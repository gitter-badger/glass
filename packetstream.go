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
    "fmt"
)

type Instance interface {
    //IsSupportedMagic(magic [4]byte) bool
    Process(orig *PacketStream, data []byte)
}

type PacketStream struct {
    inst Instance
    conn net.Conn
    // AES-GCM values
    secret [16]byte // AES-128
    nonce_init [12]byte // 96-bits
    // Packets counter
    count_in uint32
    count_out uint32
}

func shuffle(x []byte) {
    // FIXME TODO shuffle bits, not bytes!
    // See http://programming.sirrida.de/bit_perm.html (Shuffle)
    size := len(x)
    y := make([]byte, size)
    size = size >> 1;
    for i := range x {
        if i % 2 == 0 {
            y[i >> 1] = x[i]
        } else {
            y[(i >> 1) + size] = x[i]
        }
    }
    copy(x, y)
}

// Generate a AES key/nonce pair
func (this *PacketStream) generateSecret(x, y []byte) {
    // Generate AES key
    copy(this.secret[:8], x[:8])
    copy(this.secret[8:], y[:8])
    // Generate AES-GCM nonce
    copy(this.nonce_init[:6], x[8:])
    copy(this.nonce_init[6:], y[8:])
    // shuffle
    shuffle(this.secret[:])
    shuffle(this.nonce_init[:])
}

func (this *PacketStream) nonce(n uint32) ([]byte, error) {
    buf := new(bytes.Buffer)
    err := binary.Write(buf, binary.BigEndian, n)
    b := make([]byte, 12)
    copy(b, this.nonce_init[:])
    N := make([]byte, 4)
    if _, err = buf.Read(N); err != nil { return nil, err}
    b[0] ^= N[0]; b[3] ^= N[1]; b[6] ^= N[2]; b[9] ^= N[3]
    return b, nil
}

func (this *PacketStream) In(inst Instance, conn net.Conn, cert *tls.Certificate) (err error){
    var b [8]byte
    var tlsConfig *tls.Config
    var tlsConn *tls.Conn
    this.conn = conn
    // constants
    // First message sent from client
    CLIENT_HELLO := []byte("01234567")
    // Sent in response to a client hello
    SERVER_HELLO := []byte("76543210")

    fmt.Println(">hellok")
    // Phase 1: HELLO
    if _, err = conn.Read(b[:]); err != nil { return }
    if !bytes.Equal(b[:], CLIENT_HELLO) { return nil } // FIXME
    if _, err = conn.Write(SERVER_HELLO); err != nil { return }
    fmt.Println(">/hellok")
    // Phase 2: Negotiate TLS connection
    tlsConfig = &tls.Config{
        Certificates: []tls.Certificate{*cert},
        ClientAuth: tls.VerifyClientCertIfGiven,
        ServerName: "example.com",
    }
    tlsConn = tls.Server(conn, tlsConfig)
    tlsConn.Handshake()
    // Phase 3: Agree on a secret, forget about TLS
    x := make([]byte, 14)
    y := make([]byte, 14)
    if _, err = conn.Read(x); err != nil { return }
    //// Send my secret
    if _, err = io.ReadFull(rand.Reader, y); err != nil { return }
    if _, err = conn.Write(y); err != nil { return }
    // Generate AES128 data
    this.generateSecret(x, y)
    // Return successfully
    this.inst = inst
    this.count_in = 0
    this.count_out = 1
    return nil
}

func (this *PacketStream) Out(inst Instance, conn net.Conn, cert *tls.Certificate) (err error){
    var b [8]byte
    var tlsConfig *tls.Config
    var tlsConn *tls.Conn
    this.conn = conn
    // constants
    // First message sent from client
    CLIENT_HELLO := []byte("01234567")
    // Sent in response to a client hello
    SERVER_HELLO := []byte("76543210")
    // Start hello phase

    // Phase 1: HELLO
    if _, err = conn.Write(SERVER_HELLO); err != nil { return }
    if _, err = conn.Read(b[:]); err != nil { return }
    if !bytes.Equal(b[:], CLIENT_HELLO) { return nil } // FIXME
    fmt.Println("Hello phase over")
    // Phase 2: Negotiate TLS connection
    tlsConfig = &tls.Config{
        Certificates: []tls.Certificate{*cert},
        ServerName: "example.com",
    }
    tlsConn = tls.Client(conn, tlsConfig)
    tlsConn.Handshake()
    // Phase 3: Agree on a secret, forget about TLS
    x := make([]byte, 14)
    y := make([]byte, 14)
    //// Send my secret
    if _, err = io.ReadFull(rand.Reader, x); err != nil { return }
    if _, err = conn.Write(x); err != nil { return }
    //// Read client's secret
    if _, err = conn.Read(y); err != nil { return }
    //// Generate AES128 data
    this.generateSecret(x, y)
    // Return successfully
    this.inst = inst
    this.count_in = 1
    this.count_out = 0
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
    var nonce []byte
    if block, err = aes.NewCipher(s.secret[:]); err != nil { return }
	if nonce, err = s.nonce(s.count_out); err != nil { return }
	if aesgcm, err = cipher.NewGCM(block); err != nil { return }
	data = aesgcm.Seal(nil, nonce, data, nil)

    conn := s.conn
    // 1) Send size
    size := len(data)
    if size % 16 != 0 {
        return errors.New("Wrong packet size after GCM")
        // Should never happen FIXME
    }
    buf := new(bytes.Buffer)
    err = binary.Write(buf, binary.BigEndian, uint16(size / 16))
    if err != nil { return }
    two := make([]byte, 2)
    if _, err = buf.Read(two); err != nil { return }
    conn.Write(two)
    s.count_out += 2
    // 2) Send initialization vector
    //if _, err = io.ReadFull(rand.Reader, iv); err != nil { return }
    //conn.Write(iv)
    // 3) Send AES128-encrypted data
    _, err = conn.Write(data)
    return
}

func (this *PacketStream) Close() {
    this.conn.Write([]byte{'\x00', '\x00'})
}

func (s *PacketStream) Shutdown(cause error) error {
    return s.conn.Close()
}

func (s *PacketStream) Serve() {
    var nonce []byte
    var two [2]byte
    var data []byte
    var buf *bytes.Reader
    var length uint16
    var err error
    //var secret = s.secret
    conn := s.conn
    defer s.Shutdown(nil)
    for {
        // Compute nonce
        if nonce, err = s.nonce(s.count_in); err != nil { break }
        // Read packet size
        if _, err = conn.Read(two[:]); err != nil { break }
        // Increment packet count
        s.count_in += 2
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
            break
        }
        // Process packet
        go Process(s, data, nonce)
    }
    s.Shutdown(err)
}

// Decrypt the packet and give it to the instance
func Process(s *PacketStream, ciphertext, nonce []byte) {
    var err error
    var block cipher.Block
    var aesgcm cipher.AEAD

    if block, err = aes.NewCipher(s.secret[:]); err != nil { return }
	if aesgcm, err = cipher.NewGCM(block); err != nil { return }
    _, err = aesgcm.Open(ciphertext[0:], nonce, ciphertext, nil)
    if err == nil {
        s.inst.Process(s, ciphertext)
    }
}
