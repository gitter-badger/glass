package glassbox

import (
//    "fmt"
    "bytes"
    //"encoding/binary"
    "crypto/cipher"
    "crypto"
    "crypto/rsa"
    "crypto/rand"
    "crypto/aes"
    "crypto/sha256"
    "io"
    //"errors"
)

/*
func PKCS5Padding(src []byte, BlockSize int) ([]byte, int) {
    size := len(src)
    padlen := BlockSize - size % BlockSize
    nblocks := (size + padlen) / BlockSize
    padding := bytes.Repeat([]byte{byte(padlen)}, padlen)
    return append(src, padding...), nblocks
}
*/


func PKCS5UnPadding(src []byte) []byte {
    length := len(src)
    padlen := int(src[length-1])
    return src[:(length - padlen)]
}

/*  BaseEncryptedPacket is the base and most general packet format.
    It consists mainly of the TLS cipher:
    TLS_RSA_WITH_AES_128_CBC_SHA (with RSA-2048)
*/
type BaseEncryptedPacket struct {
    /*
    magic(16)
    to/from(16)
    enc_key(32) enc_sig(32)
    iv(16)
    data(...)
    */

    // The initialization vector for the AES encryption step,
    // and this packet's unique identifier
    iv      [16]byte
    // The identifier of the responsible application
    appID   [32]byte
    // Sender
    from [16]byte
    // Recipient
    to [ 16]byte
    // The encrypted key and the signature for the payload
    key     [256]byte
    sig     [256]byte
    // Cleartext
    payload []byte
    // Encrypted version of the cleartext
    enc     []byte
    // Packet direction
    incoming bool
}

func (this *BaseEncryptedPacket) Encrypt(
        appID [32]byte,
        priv *rsa.PrivateKey, pub *rsa.PublicKey,payload []byte) (err error) {
    // Set the packet as outgoing
    this.incoming = false
    // Set appID
    this.appID = appID
    //this.payload = payload
    rng := rand.Reader
    // Generate random key and initialization vector
    var key [8]byte
    var iv [16]byte
    if _, err = io.ReadFull(rng, key[:]); err != nil {
        panic("RNG failure")
    }
    if _, err = io.ReadFull(rng, iv[:]); err != nil {
        panic("RNG failure")
    }
    this.iv = iv
    // AES-encrypt payload with (key, iv)
    const blocksize = aes.BlockSize
    //// the payload will be PKCS5-padded on the fly
    size := len(payload)
    padlen := blocksize - size % blocksize
    padding := bytes.Repeat([]byte{byte(padlen)}, padlen)

    //// Total size: orig size + padding + iv + appID
    buf := make([]byte, size + padlen + 16 + 32)

    aes, _ := aes.NewCipher(key[:])
    cbc := cipher.NewCBCEncrypter(aes, iv[:])
    j := size - (size % blocksize)
    if j > 0 {
        cbc.CryptBlocks(buf[:j], payload[:j])
        cbc = cipher.NewCBCEncrypter(aes, buf[j-blocksize:j])
    }
    var buffer bytes.Buffer
    buffer.Write(payload[j:])
    buffer.Write(iv[:])
    buffer.Write(appID[:])
    buffer.Write(padding)
    cbc.CryptBlocks(buf[j:], buffer.Bytes())
    this.enc = buf

    // Sign the data
    //// "G" app_id(32) "L" iv(16) "A" partner(16) "S"
    //// SHA256, should be replaced with SHA1? FIXME!
	d := sha256.New()
	d.Reset()
    for _, k := range [][]byte{
            []byte("G"), appID[:], []byte("L"), iv[:], []byte("A"), this.to[:], []byte("S"),
        } {
        d.Write(k)
    }
	d.Write(payload)
    // FIXME: is padding missing? is it automatic?
	hashed := d.Sum(nil)
    signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed[:])
    if err != nil {
        return
    }
    // assert len(signature) == 256
    copy(this.sig[:], signature)
    return nil
}

func (this *BaseEncryptedPacket) Bytes() ([]byte, error) {
    buf := new(bytes.Buffer)
    // magic (16 octets)
    buf.WriteString("\xff\x01\x00\x00\x00\x00\x00\x00")
    buf.WriteString("\xff\x01\x00\x00\x00\x00\x00\x00")

    buf.Write(this.to[:])
    buf.Write(this.key[:])
    buf.Write(this.sig[:])
    buf.Write(this.iv[:])
    buf.Write(this.enc)

    return buf.Bytes(), nil
}
//err = binary.Write(buf, binary.LittleEndian, uint16(size / 16))

// Parse a packet from a byte stream
func (this *BaseEncryptedPacket) FromBytes(data []byte) error {
    copy(this.from[:], data[      :16    ])
    copy(this.key    [:], data[16    :16+256])
    copy(this.sig    [:], data[16+256:16+512])
    copy(this.iv     [:], data[16+512:32+512])
    this.enc = data[32+512:]
    // Decrypt here!
    return nil
}

func (this *BaseEncryptedPacket) From() [16]byte { return this.from }
func (this *BaseEncryptedPacket) To() [16]byte { return this.to }
func (this *BaseEncryptedPacket) Id() [16]byte { return this.iv }

func (this *BaseEncryptedPacket) Decrypt(priv *rsa.PrivateKey, pub *rsa.PublicKey) []byte {
    return []byte{}
}
