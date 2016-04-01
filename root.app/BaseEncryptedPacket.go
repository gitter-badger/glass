package main

import (
    "fmt"
    "bytes"
    "encoding/binary"
    "crypto/rsa"
    "crypto/rand"
    "crypto/aes"
    "io"
)

type Key struct {
    n int
    m int
    p int
    d int
    e int
}

/*  BaseEncryptedPacket is the base and most general packet format.
    It consists mainly of the TLS cipher:
    RSA-2048-with-AES-128-CBC-SHA
*/
type BaseEncryptedPacket struct {
    /*
    0xff(1) magic(4) head_size(1) body_size(2)
    partner(16)
    enc_key(32) enc_sig(32)
    iv(16)
    data(...)
    */
    partner [16]byte
    enc_key [32]byte
    enc_sig [32]byte
    iv      [16]byte
    data    [  ]byte
}

func (this *BaseEncryptedPacket) init(priv *rsa.PrivateKey, pub *rsa.PublicKey, payload []byte) (err error) {
    rng := rand.Reader
    // Random key and initialization vector
    key := make([]byte, 16)
    iv := make([]byte, 16)
    _, err = io.ReadFull(rng, key);
    _, err = io.ReadFull(rng, key);
    //if err != nil {
    //    panic("RNG failure")
    //}
    // encrypt payload with aes key,iv
    blocksize := 16 // FIXME?
    // payload = pad payload blocksize
    // size = len(payload) / blocksize
    c, err := aes.NewCipher(key)
    if err != nil {
        return
    }
    out := make([]byte, len(payload))
    c.Encrypt(out, payload)

    out, err = rsa.EncryptPKCS1v15(rng, pub, key)

    return nil
}

func (this *BaseEncryptedPacket) Bytes() []byte {
    var err error
    buf := new(bytes.Buffer)
    // 1) 0xff (1 octet)
    // 2) magic (4 octets)
    buf.WriteString("\xff\x01\x00\x00\x00")
    // 3) Head size (1 octet)
    buf.WriteString("\x2c") // 44*2 bytes use uint8 here!
    // 4) Body size
    //
    // Write the size of the payload in AES blocks
    // (including padding!)
    var size int = len(this.data)
    err = binary.Write(buf, binary.LittleEndian, uint16(size / 16))

    buf.Write(this.partner[:])
    buf.Write(this.enc_key[:])
    buf.Write(this.enc_sig[:])
    buf.Write(this.iv[:])
    buf.Write(this.data)

    return buf.Bytes()
}
