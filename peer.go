package glass

type Peer struct {}

func (Peer) IsTrusted() bool { return false }
func (Peer) Trust() {}
