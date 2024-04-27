package command

import "sync"

type KV struct {
	data map[string][]byte
	*sync.RWMutex
}

func NewKV() *KV {
	return &KV{
		data:    map[string][]byte{},
		RWMutex: &sync.RWMutex{},
	}
}

func (kv *KV) Set(key string, value []byte) error {
	kv.Lock()
	defer kv.Unlock()

	kv.data[key] = value

	return nil
}

func (kv *KV) Get(key string) ([]byte, bool) {
	kv.RLock()
	defer kv.RUnlock()
	data, ok := kv.data[key]

	return data, ok
}
