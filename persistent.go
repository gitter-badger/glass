package glass

type Storage interface {
    Store(payload []byte) (key string)
    Delete(key string)
}
