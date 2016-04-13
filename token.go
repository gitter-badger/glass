package glass

import (
    "crypto/rsa"
    "net"
)

type AuthToken struct {
    // Router data
    RouterHost string
    RouterIP net.IP
    RouterPublicKey rsa.PublicKey
    // My data
    PublicKey rsa.PublicKey
    ApplicationID [16]byte
    AppToken [16]byte
    AppSecret [16]byte
}

func (*AuthToken) Read(bs []byte) error {
    return nil
}
