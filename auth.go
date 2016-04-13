package glass

import (
    "crypto/rsa"
)

type AuthToken struct {
    // Router data
    RouterIP [16]byte //IPv6
    RouterPublicKey rsa.PublicKey

    PublicKey rsa.PublicKey
    AppID [16]byte
    AppInstanceToken [8]byte
}
