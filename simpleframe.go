package glass

import (
    "io"
    "bytes"
    "crypto"
    "crypto/rsa"
    "crypto/aes"
    "crypto/sha256"
    "crypto/cipher"
    //"errors"
)

/*
Structure of the frame:

|-------|-------|
| Name  | Bytes |
|-------|-------|
| 0xff  | 1     |
| 0x01  | 1     |
| ?     | 14    |
| from  | 16    |
| to    | 16    |
| key   | 256   |
| sig   | 256   |
| data  | var   |

*/


const FRAME_SIMPLE = "\xff\x01"

/*
func PKCS5Padding(src []byte, BlockSize int) ([]byte, int) {
    size := len(src)
    padlen := BlockSize - size % BlockSize
    nblocks := (size + padlen) / BlockSize
    padding := bytes.Repeat([]byte{byte(padlen)}, padlen)
    return append(src, padding...), nblocks
}
func PKCS5UnPadding(src []byte) []byte {
    length := len(src)
    padlen := int(src[length-1])
    return src[:(length - padlen)]
}
*/

/*  SimpleFrame is the base and most general packet format.
    It consists mainly of the TLS cipher:
    TLS_RSA_WITH_AES_128_CBC_SHA (with RSA-2048)
*/
type SimpleFrame struct {
    // The identifier of the responsible application
    AppName [16]byte
    // Sender
    From [16]byte
    // Recipient
    To [16]byte
    // The encrypted key and the signature for the payload
    key [256]byte
    sig [256]byte
    // Cleartext
    Content []byte
    // Encrypted version of the cleartext
    enc []byte
}

func (frame *SimpleFrame) Seal(
        priv *rsa.PrivateKey,
        pub *rsa.PublicKey,
    ) (err error) {
    // Generate random key
    var key [8]byte
    if _, err = io.ReadFull(rng, key[:]); err != nil {
        panic("RNG failure")
    }
    // AES-encrypt payload with (key, iv)
    const blocksize = aes.BlockSize
    iv := bytes.Repeat([]byte{0x00}, blocksize)
    //// the payload will be PKCS5-padded on the fly
    size := len(frame.Content)
    padlen := blocksize - size % blocksize
    padding := bytes.Repeat([]byte{byte(padlen)}, padlen)

    //// Total size: orig size + padding + iv + appID
    buf := make([]byte, size + padlen + 16 + 32)

    aes, _ := aes.NewCipher(key[:])
    cbc := cipher.NewCBCEncrypter(aes, iv[:])
    j := size - (size % blocksize)
    if j > 0 {
        cbc.CryptBlocks(buf[:j], frame.Content[:j])
        cbc = cipher.NewCBCEncrypter(aes, buf[j-blocksize:j])
    }
    var buffer bytes.Buffer
    buffer.Write(frame.Content[j:])
    buffer.Write(iv[:])
    buffer.Write(frame.AppName[:])
    buffer.Write(padding)
    cbc.CryptBlocks(buf[j:], buffer.Bytes())
    frame.enc = buf

    // Sign the data
    //// "G" app_id(32) "L" iv(16) "A" partner(16) "S"
    //// SHA256, should be replaced with SHA1? FIXME!
	d := sha256.New()
	d.Reset()
    for _, k := range [][]byte{
            []byte("G"), frame.AppName[:], []byte("L"), iv[:], []byte("A"), frame.To[:], []byte("S"),
        } {
        d.Write(k)
    }
	d.Write(frame.Content)
    // FIXME: is padding missing? is it automatic?
	hashed := d.Sum(nil)
    signature, err := rsa.SignPKCS1v15(rng, priv, crypto.SHA256, hashed[:])
    if err != nil {
        return
    }
    // assert len(signature) == 256
    copy(frame.sig[:], signature)
    return nil
}

func (frame *SimpleFrame) Bytes() []byte {
    buf := new(bytes.Buffer)
    // magic (16 byte)
    buf.WriteString("\xff\x01\x00\x00\x00\x00\x00\x00")
    buf.WriteString("\x00\x00\x00\x00\x00\x00\x00\x00")

    buf.Write(frame.From[:]) // 16 byte
    buf.Write(frame.To  [:]) // 16 byte
    buf.Write(frame.key[:]) // 32 byte
    buf.Write(frame.sig[:]) // 32 byte
    buf.Write(frame.enc)

    return buf.Bytes()
}
//err = binary.Write(buf, binary.LittleEndian, uint16(size / 16))

// Parse a packet from a byte stream
func (frame *SimpleFrame) Read(data []byte) bool {
    if len(data) < 512 + 32 {
        return false
    }
    copy(frame.From   [:], data[      :16    ])
    copy(frame.To     [:], data[      :32    ])
    copy(frame.key    [:], data[32    :32+256])
    copy(frame.sig    [:], data[32+256:32+512])
    frame.enc = data[32+512:]
    return true
}

func (*SimpleFrame) Id() [16]byte { return *new([16]byte) } // FIXME!

func (frame *SimpleFrame) Open(priv *rsa.PrivateKey, pub *rsa.PublicKey) []byte {
    return []byte{}
}
