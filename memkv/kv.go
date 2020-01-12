package memkv

import (
	"os"
	"sync"

	"github.com/KarpelesLab/vfs"
)

type entry []byte
type memKv struct {
	data map[string]entry
	lk   sync.RWMutex
}

func New() vfs.Keyval {
	return &memKv{
		data: make(map[string]entry),
	}
}

func (e *memKv) Get(key string) (vfs.KVEntry, error) {
	e.lk.RLock()
	f, ok := e.data[key]
	e.lk.RUnlock()

	if !ok {
		return nil, os.ErrNotExist
	}

	return f, nil
}

func (e *memKv) Put(key string, value vfs.KVEntry) error {
	data, err := value.Data()
	if err != nil {
		return err
	}

	e.lk.Lock()
	e.data[key] = entry(data)
	e.lk.Unlock()

	return nil
}

func (e *memKv) Delete(key string) error {
	e.lk.Lock()
	delete(e.data, key)
	e.lk.Unlock()
	return nil
}

func (d entry) Data() ([]byte, error) {
	return []byte(d), nil
}
