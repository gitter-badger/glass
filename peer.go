package glass

import (
    "crypto/rsa"
)

type Peer struct {
    Addr string
    PublicKey rsa.PublicKey
}

func (Peer) IsTrusted() bool { return false }
func (Peer) Trust() {}
