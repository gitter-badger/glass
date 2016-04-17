package glass

import (
    "net"
    "errors"
    "bytes"
    "crypto/rand"
    "crypto/tls"
//    "crypto/x509"
    "crypto/aes"
    "crypto/cipher"
    "hash"
    //"crypto/rsa"
    "io"
    "log"
    "fmt"
    "bufio"
    //"time"
)

type StreamDirection int
const (
    STREAM_IN = -1
    STREAM_OUT = 1
)
var client_hello []byte
var server_hello []byte
var strongCipherSuites []uint16
var rng io.Reader

type FrameStream struct {
    FrameHandler func(string, []byte)
    Direction StreamDirection
    Conn net.Conn
    in io.Reader
    out io.Writer
    // FNV-1 counters
    hash_in hash.Hash
    hash_out hash.Hash
    // AES-GCM values
    secret []byte // 16 bytes for AES-128
    nonce_init [12]byte // 96-bits
    // Frames counter
    in_seq_no uint32
    out_seq_no uint32
    // Queue
    queue [][]byte
}

// Initialize global constants
func init() {
    // HELLO constants
    client_hello = []byte("01234567") // First message sent from client
    server_hello = []byte("76543210") // Response to a client hello
    // Choose cipher suites
    strongCipherSuites = []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA}
    // Set up the Random Numbers Generator
    rng = rand.Reader
}

func (stream *FrameStream) Handshake() error {
    // Initialize buffered I/O
    stream.in = bufio.NewReader(stream.Conn)
    stream.out = bufio.NewWriter(stream.Conn)
    // Initialize hashes
    stream.hash_in = NewFNV1()
    stream.hash_out = NewFNV1()
    // Continue according to stream direction
    switch stream.Direction {
    case STREAM_IN:
        return stream.server()
    case STREAM_OUT:
        return stream.client()
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
    this.secret = make([]byte, 16)
    // Generate AES key
    copy(this.secret[0:8], x[0:8])
    copy(this.secret[8:16], y[0:8])
    // Generate 12-byte GCM nonce
    copy(this.nonce_init[0:6], x[8:14])
    copy(this.nonce_init[6:12], y[8:14])
    // Shuffle the bytes. Unnecessary? But fun.
    shuffle(this.secret[:])
    shuffle(this.nonce_init[:])
    // Initialize hashes
    this.hash_in.Write(x[14:16])
    this.hash_out.Write(y[14:16])
}

// Generate a unique nonce for the given n in the current session
func (stream *FrameStream) nonce(seq_no uint32) (b []byte, err error) {
    buf := new(bytes.Buffer)
    write_uint32(buf, seq_no)
    N := buf.Next(4)
    b = make([]byte, 12)
    copy(b, stream.nonce_init[:])
    b[0] ^= N[0]; b[3] ^= N[1]; b[6] ^= N[2]; b[9] ^= N[3]
    return b, nil
}

/*func (stream *FrameStream) Dial(app *App, router *Peer) error {
    conn, err := net.Dial("tcp", "localhost:3001")
    if err != nil { err }
    stream.In(app, conn)
}*/

func (stream *FrameStream) server() (err error) {
    fmt.Println("[Server] Starting server handshake")
    // Phase 1: HELLO
    eight := make([]byte, 8)
    if _, err = stream.Conn.Read(eight); err != nil {
        return errors.New("FrameStream.In: Can't read Client Hello")
    }
    if !bytes.Equal(eight, client_hello) {
        return errors.New("FrameStream.In: Wrong Client Hello")
    }
    if _, err = stream.Conn.Write(server_hello); err != nil {
        return errors.New("FrameStream.In: Can't write Server Hello")
    }
    log.Println("[Server] Simple Hello Phase Over")
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
    // Read client_random
    if _, err = stream.Conn.Read(x); err != nil { return }
    // Send server_random
    if _, err = io.ReadFull(rng, y); err != nil { return }
    if _, err = stream.Conn.Write(y); err != nil { return }
    // Generate AES data
    stream.generateSecret(x, y)

    stream.in_seq_no = 1
    stream.out_seq_no = 0
    return
}

func (stream *FrameStream) client() (err error) {
    fmt.Println("[Client] Starting client handshake.")
    // Phase 1: HELLO
    eight := make([]byte, 8)
    if _, err = stream.Conn.Write(client_hello); err != nil { return }
    if _, err = stream.Conn.Read(eight); err != nil { return }
    if !bytes.Equal(eight, server_hello) { return nil } // FIXME
    log.Println("[Client] Simple Hello Phase Over")

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

    // Phase 3: Exchange random, forget about TLS
    x := make([]byte, 16)
    y := make([]byte, 16)
    // - Send client_random
    if _, err = io.ReadFull(rng, x); err != nil { return }
    if _, err = stream.Conn.Write(x); err != nil { return }
    // - Read server_random
    if _, err = stream.Conn.Read(y); err != nil { return }
    // - Generate AES128-GCM data
    stream.generateSecret(x, y)

    stream.in_seq_no = 0
    stream.out_seq_no = 1
    return nil
}

func (stream *FrameStream) Send(f Frame) error {
    // FIXME: queue
    return stream.write(f)
}

func (s *FrameStream) write(p Frame) (err error) {
    // See: https://gist.github.com/kkirsche/e28da6754c39d5e7ea10
    data := p.Bytes()
    length := len(data)
    if length == 0 || length % 16 != 0 {
        return errors.New("Wrong packet size")
    }
    // Encrypt data with AES128-GCM
    block, err := aes.NewCipher(s.secret)
    if err != nil { return }
    nonce, err := s.nonce(s.out_seq_no)
    fmt.Printf("??? %b", nonce)
    fmt.Printf("???! %b", s.secret)
    if err != nil { return }
    aesgcm, err := cipher.NewGCM(block)
    if err != nil { return }
    ciphertext := aesgcm.Seal(nil, nonce, data, nil)

    // Write header
    //////////////////////////////////////////////////////
    // TO DO HERE:
    // fix the new header with hashing and counting
    // s.hash_out.Write(ciphertext)
    // ^%^&*^%$#$%^&%$#@$%^&^%$#@$%^&%$#@$%^
    //////////////////////////////////////////////////////
    ret := new(bytes.Buffer)
    // Write the length of the payload (2 bytes)
    // Since we are using AES-GCM, the new length
    // of the data will be the old length plus gcmTagSize,
    // which is 16 bytes
    // assert len(data) == length + gcmTagSize
    write_uint16(ret, uint16(length / 16))
    s.out_seq_no += 2
    // Write frame type
    ret.WriteString(p.Type())
    // Write in_seq_no
    write_uint32(ret, s.in_seq_no)
    // Write ciphertext
    _, err = ret.Write(ciphertext)
    // Send the frame
    fmt.Printf("[Client] Sending frame Length=%d\n", len(ciphertext))
    s.Conn.Write(ret.Bytes())
    fmt.Println("[Client] Frame sent")
    return
}

func (stream *FrameStream) Close() (err error) {
    if stream.Conn != nil {
        log.Println("Closing stream connection")
        // Try to close it gently
        // TODO setWriteDeadline
        stream.Conn.Write([]byte("\x00\x00\x00\x00\x00\x00\x00\x00"))
        err = stream.Conn.Close()
    }
    return
}

func (stream *FrameStream) Serve() {
    var nonce []byte
    var eight [8]byte
    var header string
    var data []byte
    var length int
    var err error
    var conn = stream.Conn
    //conn.SetReadDeadline(time.Time(0)) // FIXME
    defer stream.Close()
    for {
        // Compute nonce
        if nonce, err = stream.nonce(stream.in_seq_no); err != nil { break }
        fmt.Printf("!!! %b", nonce)
        fmt.Printf("???! %b", stream.secret)
        // Read eight bytes
        if _, err = conn.Read(eight[:]); err != nil {
            fmt.Println("[Server] Error reading packet header: %e.", err)
        break }
        log.Println("[<-] Incoming packet.")
        // Increment packet count
        stream.in_seq_no += 2
        header = string(eight[:]) // ^ string(stream.hash_in.Sum(nil)) FIXME
        if len(header) != 8 {
            fmt.Printf("[<-] Wrong header length %d.\n", len(header))
            break
        } // FIXME
        fmt.Printf("[<-] Header: %s.\n", header)
        length = int(read_uint16(header[0:2]))
        // Should we close the connection?
        if length == 0 {
            fmt.Println("[Server] Empty frame. Exiting.")
            break
        }
        length = 16 * length + 16
        fmt.Printf("[Server] Incoming packet. Length: %d.\n", length)
        // FIXME Check number of packets
        if read_uint32(header[4:8])  > stream.out_seq_no {
            fmt.Println("[Server] Wrong sequence number. Closing connection")
            break
        }
        // Read *length* blocks of data
        data = make([]byte, length)
        if _, err = conn.Read(data); err != nil { break }
        // Process packet
        go stream.handleFrame(string(header[2:4]), data, nonce)
    }
    log.Println("[Server] Goodbye")
}

func (stream *FrameStream) handleFrame(frame_type string, ciphertext, nonce []byte) {
    var err error
    var block cipher.Block
    var aesgcm cipher.AEAD
    log.Println("[--] Processing incoming packet")
    if block, err = aes.NewCipher(stream.secret); err != nil { return }
    if aesgcm, err = cipher.NewGCM(block); err != nil { return }
    payload, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        // If the message was not encrypted correctly,
        // just ignore it. Does this lead to problems?
        log.Println("Message decryption failed")
        return
    }
    // stream.hash_in.Write(payload) // FIXME
    log.Println("[  ] Incoming packet correctly decrypted")
    stream.FrameHandler(frame_type, payload)
}
