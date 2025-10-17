package key_mutex

type IKeyMutex interface {
	Lock(key string) func()
}
