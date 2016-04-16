package glass

import (
    "net"
    "errors"
    "bytes"
    "encoding/binary"
    "crypto/rand"
    "crypto/tls"
//    "crypto/x509"
    "crypto/aes"
    "crypto/cipher"
    //"crypto/rsa"
    "io"
    "log"
    //"time"
)

type StreamDirection int
const (
    STREAM_IN = 0
    STREAM_OUT = 1
)

type FrameStream struct {
    FrameHandler func([]byte)
    Direction StreamDirection
    Conn net.Conn
    // AES-GCM values
    secret [16]byte // AES-128
    nonce_init [12]byte // 96-bits
    // Frames counter
    count_in uint32
    count_out uint32
    // Queue
    queue [][]byte
}

var CLIENT_HELLO []byte
var SERVER_HELLO []byte
var strongCipherSuites []uint16
var RNG io.Reader

// Initialize global constants
func init() {
    // HELLO constants
    CLIENT_HELLO = []byte("01234567") // First message sent from client
    SERVER_HELLO = []byte("76543210") // Response to a client hello
    // Choose cipher suites
    strongCipherSuites = []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA}
    //
    RNG = rand.Reader
}

func (stream *FrameStream) Handshake() error {
    switch stream.Direction {
    case STREAM_IN:
        return stream.in()
    case STREAM_OUT:
        return stream.out()
    default:
    }
    return errors.New("FrameStream.Handshake: Invalid stream direction.")
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

// Generate a AES128 key with a 96-bit nonce for GCM mode
func (this *FrameStream) generateSecret(x, y []byte) {
    // Generate AES key
    copy(this.secret[0:8], x[0:8])
    copy(this.secret[8:16], y[0:8])
    // Generate 12-byte GCM nonce
    copy(this.nonce_init[0:6], x[8:14])
    copy(this.nonce_init[6:12], y[8:14])
    // Shuffle the bytes. Unnecessary? But fun.
    shuffle(this.secret[:])
    shuffle(this.nonce_init[:])
}

// Generate a unique nonce for the given n in the current session
func (stream *FrameStream) nonce(n uint32) ([]byte, error) {
    buf := new(bytes.Buffer)
    err := binary.Write(buf, binary.BigEndian, n)
    b := make([]byte, 12)
    copy(b, stream.nonce_init[:])
    N := make([]byte, 4)
    if _, err = buf.Read(N); err != nil { return nil, err}
    b[0] ^= N[0]; b[3] ^= N[1]; b[6] ^= N[2]; b[9] ^= N[3]
    return b, nil
}

/*func (stream *FrameStream) Dial(app *App, router *Peer) error {
    conn, err := net.Dial("tcp", "localhost:3001")
    if err != nil { err }
    stream.In(app, conn)
}*/

func (this *FrameStream) in() (err error) {
    var conn = this.Conn
    // Phase 1: HELLO
    eight := make([]byte, 8)
    if _, err = conn.Read(eight); err != nil {
        return errors.New("FrameStream.In: Can't read Client Hello")
    }
    if !bytes.Equal(eight, CLIENT_HELLO) {
        return errors.New("FrameStream.In: Wrong Client Hello")
    }
    if _, err = conn.Write(SERVER_HELLO); err != nil {
        return errors.New("FrameStream.In: Can't write Server Hello")
    }

    // Phase 2: (Authenticate & establish forward secrecy)
    // for now, just negotiate a TLS connection
    /* FIXME TODO
    caCertPool := x509.NewCertPool()
    caCertPool.AddCert(cert.Leaf)
    tlsConfig := &tls.Config{
        //CipherSuites: strongCipherSuites,
        Certificates: []tls.Certificate{cert},
        ClientCAs: caCertPool,
        ClientAuth: tls.VerifyClientCertIfGiven,
        SessionTicketsDisabled: true,
    }
    tlsConn := tls.Server(conn, tlsConfig)
    tlsConn.Handshake()
    */

    // Phase 3: Agree on a secret, forget about TLS
    x := make([]byte, 16)
    y := make([]byte, 16)
    if _, err = conn.Read(x); err != nil { return }
    // - Send my secret
    if _, err = io.ReadFull(rand.Reader, y); err != nil { return }
    if _, err = conn.Write(y); err != nil { return }
    // Generate AES data
    this.generateSecret(x, y)

    this.count_in = 1
    this.count_out = 0
    return
}

func (this *FrameStream) out() (err error) {
    var conn = this.Conn
    // Phase 1: HELLO
    eight := make([]byte, 8)
    if _, err = conn.Write(CLIENT_HELLO); err != nil { return }
    if _, err = conn.Read(eight); err != nil { return }
    if !bytes.Equal(eight, SERVER_HELLO) { return nil } // FIXME
    log.Println("[->] Simple Hello Phase Over")

    // Phase 2: Negotiate TLS connection
    // This is temporary: in the future, the session
    // key will be negotiated directly with
    // RSA-2048-AES-...etc, without a TLS handshake.
    /* Phase currently disabled, until I understand better
       the go libraries... FIXME! TODO!
    caCertPool := x509.NewCertPool()
    caCertPool.AddCert(hercert)
    tlsConfig := &tls.Config{
        CipherSuites: strongCipherSuites,
        SessionTicketsDisabled: true,
        RootCAs: caCertPool,
    }
    if mycert != nil { // Client certificate is optional
        tlsConfig.Certificates = []tls.Certificate{*mycert}
    }
    tlsConn := tls.Client(conn, tlsConfig)
    tlsConn.Handshake()
    // Verify server identity
    var peerCertificates = tlsConn.ConnectionState().PeerCertificates
    */

    // Phase 3: Agree on a secret, forget about TLS
    x := make([]byte, 16)
    y := make([]byte, 16)
    // - Send my secret
    if _, err = io.ReadFull(rand.Reader, x); err != nil { return }
    if _, err = conn.Write(x); err != nil { return }
    // - Read client's secret
    if _, err = conn.Read(y); err != nil { return }
    // - Generate AES128-GCM data
    this.generateSecret(x, y)

    this.count_in = 0
    this.count_out = 1
    return nil
}

func (s *FrameStream) Send(p Frame) error {
    // FIXME: queue
    return s.write(p)
}

func (s *FrameStream) write(p Frame) (err error) {
    // See: https://gist.github.com/kkirsche/e28da6754c39d5e7ea10
    data := p.Bytes()
    length := len(data)
    if length == 0 || length % 16 != 0 {
        return errors.New("Wrong packet size")
    }
    // Encrypt data with AES128-GCM
    if block, err := aes.NewCipher(s.secret[:]); err != nil { return }
    if nonce, err := s.nonce(s.count_out); err != nil { return }
    if aesgcm, err := cipher.NewGCM(block); err != nil { return }
    ciphertext := aesgcm.Seal(nil, nonce, data, nil)

    ret := new(bytes.Buffer)
    // Write the length of the payload (2 bytes)
    // Since we are using AES-GCM, the new length
    // of the data will be the old length plus gcmTagSize,
    // which is 16 bytes
    // assert len(data) == length + gcmTagSize
    writeLength(ret, uint16(1 + length / 16))
    s.count_out += 2
    // Write ciphertext
    _, err = ret.Write(ciphertext)
    // Send the frame
    s.out.Write(ret.Bytes())
    log.Println("[->] Frame sent")
    return
}

func (stream *FrameStream) Close() (err error) {
    if stream.Conn != nil {
        log.Println("Closing stream connection")
        // Try to close it gently
        // TODO setWriteDeadline
        stream.Conn.Write([]byte{'\x00', '\x00'})
        err = stream.Conn.Close()
    }
    return
}

func (s *FrameStream) Serve() {
    var nonce []byte
    var two [2]byte
    var data []byte
    var buf *bytes.Reader
    var length uint16
    var err error
    //var secret = s.secret
    var conn = s.Conn
    //conn.SetReadDeadline(time.Time(0)) // FIXME
    defer s.Close()
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
        log.Printf("[<-] Incoming packet. Size=%d\n", length)
        // Read *length* blocks of data
        data = make([]byte, int(length) * 16)
        _, err = conn.Read(data)
        if err != nil {
            break
        }
        // Process packet
        s.processFrame(data, nonce)
    }
    log.Println("[<-] Goodbye")
}

func (stream *FrameStream) processFrame(ciphertext, nonce []byte) {
    var err error
    var block cipher.Block
    var aesgcm cipher.AEAD
    log.Println("[--] Processing incoming packet")
    if block, err = aes.NewCipher(stream.secret[:]); err != nil { return }
    if aesgcm, err = cipher.NewGCM(block); err != nil { return }
    payload, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        // If the message was not encrypted correctly,
        // just ignore it. Does this lead to problems?
        log.Println("Message decryption failed")
        return
    }
    log.Println("[  ] Incoming packet correctly decrypted")
    stream.FrameHandler(payload)
}
