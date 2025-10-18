package key_mutex

type IKeyMutex interface {
	RLock(key string) func()
	Lock(key string) func()
}
