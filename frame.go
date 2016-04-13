package glass

type Frame interface {
    // Frame identifier
    Id() [16]byte
    // Get sender/recipient
    From() [16]byte
    To()   [16]byte
    // Get packet's payload
    Content() []byte
    // Get raw representation
    Bytes() []byte
}

/*
type MsgEncrypted struct {
    magic   [ 8]byte
    partner [ 8]byte
    channel [ 8]byte
    size    [16]byte
    iv      [16]byte
    data    [  ]byte
}
*/
