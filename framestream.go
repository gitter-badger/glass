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
)

type FrameStream struct {
    handler Handler
    conn net.Conn
    // AES-GCM values
    secret [16]byte // AES-128
    nonce_init [12]byte // 96-bits
    // Frames counter
    count_in uint32
    count_out uint32
}

var CLIENT_HELLO []byte
var SERVER_HELLO []byte
var strongCipherSuites []uint16

// Initialize global constants
func init() {
    // HELLO constants
    CLIENT_HELLO = []byte("01234567") // First message sent from client
    SERVER_HELLO = []byte("76543210") // Response to a client hello
    // Choose cipher suites
    strongCipherSuites = []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA}
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
func (this *FrameStream) nonce(n uint32) ([]byte, error) {
    buf := new(bytes.Buffer)
    err := binary.Write(buf, binary.BigEndian, n)
    b := make([]byte, 12)
    copy(b, this.nonce_init[:])
    N := make([]byte, 4)
    if _, err = buf.Read(N); err != nil { return nil, err}
    b[0] ^= N[0]; b[3] ^= N[1]; b[6] ^= N[2]; b[9] ^= N[3]
    return b, nil
}

func (this *FrameStream) In(
        handler Handler, conn net.Conn,
//        cert tls.Certificate,
    ) (err error) {
    this.conn = conn

    // Phase 1: HELLO
    eight := make([]byte, 8)
    if _, err = conn.Read(eight); err != nil { return }
    if !bytes.Equal(eight, CLIENT_HELLO) { return nil } // FIXME
    if _, err = conn.Write(SERVER_HELLO); err != nil { return }

    // Phase 2: (Authenticate & establish forward secrecy)
    // for now, just negotiate a TLS connection
    /* FIXME TODO
    caCertPool := x509.NewCertPool()
    caCertPool.AddCert(cert.Leaf)
    tlsConfig := &tls.Config{
        CipherSuites: strongCipherSuites,
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

    this.handler = handler
    this.count_in = 0
    this.count_out = 1
    return
}

func (this *FrameStream) Out(
        handler Handler, conn net.Conn,
//        mycert *tls.Certificate, mycert2 *x509.Certificate, hercert *x509.Certificate,
    ) (err error) {
    this.conn = conn

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

    this.handler = handler
    this.count_in = 1
    this.count_out = 0
    return
}

func (s *FrameStream) Send(p Frame) error {
    // FIXME: queue
    return s.write(p)
}

func (s *FrameStream) write(p Frame) (err error) {
    // See: https://gist.github.com/kkirsche/e28da6754c39d5e7ea10
    data := p.Bytes()
    size := len(data)
    if size % 16 != 0 || size == 0 {
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

    // Send size of the payload (2 bytes)
    if len(data) != size + 16 {
        log.Printf("%d %d", len(data), size)
        return errors.New("AES-GCM encryption wrong packet size")
        // Should never happen FIXME CHECk
    }
    buf := new(bytes.Buffer)
    err = binary.Write(buf, binary.BigEndian, uint16(1 + size / 16))
    if err != nil { return }
    two := make([]byte, 2)
    if _, err = buf.Read(two); err != nil { return }
    log.Printf("[->] Send length: %s\n", string(two[:]))
    s.conn.Write(two)
    s.count_out += 2
    // 2) Send initialization vector
    //if _, err = io.ReadFull(rand.Reader, iv); err != nil { return }
    //conn.Write(iv)
    // 3) Send AES128-encrypted data
    _, err = s.conn.Write(data)
    log.Println("[->] Frame sent")
    return
}

func (this *FrameStream) Close() {
    this.conn.Write([]byte{'\x00', '\x00'})
}

func (s *FrameStream) Shutdown(cause error) error {
    return s.conn.Close()
}

func (s *FrameStream) Serve() {
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
        log.Printf("[<-] Incoming packet. Size=%d\n", length)
        // Read *length* blocks of data
        data = make([]byte, int(length) * 16)
        _, err = conn.Read(data)
        if err != nil {
            break
        }
        // Process packet
        s.process(data, nonce)
    }
    log.Println("[<-] Goodbye")
    s.Shutdown(err)
}

// Decrypt the packet and give it to the app
func (s *FrameStream) process(ciphertext, nonce []byte) {
    var err error
    var block cipher.Block
    var aesgcm cipher.AEAD
    log.Println("[--] Processing incoming packet")
    if block, err = aes.NewCipher(s.secret[:]); err != nil { return }
	if aesgcm, err = cipher.NewGCM(block); err != nil { return }
    payload, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        // If the message was not encrypted correctly,
        // just ignore it. Does this lead to problems?
        log.Println("Message decryption failed")
        return
    }
    log.Println("[  ] Incoming packet correctly decrypted")
    // Parse the packet type!
    //ciphertext = ciphertext[16:]
    magic := string(payload[0:2])
    if magic == PACKET_SIMPLE {
        p := new(SimpleFrame)
        p.FromBytes(payload)
        s.handler.ProcessSimpleFrame(p)
    } else if (magic == PACKET_TEST) {
        log.Println("[  ] Received test packet, processing...")
        p := new(TestFrame)
        p.FromBytes(payload)
        s.handler.ProcessTestFrame(p)
    } else {
        // If the message type is not supported,
        // just ignore it. Does this lead to problems?
        log.Println("[  ] Unknown incoming packet type: ingnoring.")
    }
}
