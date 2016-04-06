package glassbox

type Entity struct {}

func (e Entity) IsTrusted() bool {
    return false;
}
func (e Entity) Trust() {}
