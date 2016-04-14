package glass

import (
    //"crypto/rsa"
//    "net"
)

type AuthToken struct {
    Router *Peer
    // My data
    Me *Peer
    ApplicationID [16]byte
    AppToken [16]byte
    AppSecret [16]byte
}

func (*AuthToken) Read(bs []byte) error {
    return nil
}
