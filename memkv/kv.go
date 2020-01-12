package memkv

import (
	"os"
	"strings"
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

func (e memKv) List(prefix string, callback func(key string, value vfs.KVEntry) (bool, error)) error {
	cop := make(map[string]entry)

	// make a copy so callback can use functions that use a lock
	e.lk.RLock()
	if prefix == "" {
		for k, v := range e.data {
			cop[k] = v
		}
	} else {
		for k, v := range e.data {
			if strings.HasPrefix(k, prefix) {
				cop[k] = v
			}
		}
	}
	e.lk.RUnlock()

	for k, v := range cop {
		cont, err := callback(k, v)
		if err != nil {
			return err
		}
		if !cont {
			return nil
		}
	}
	return nil
}

func (d entry) Data() ([]byte, error) {
	return []byte(d), nil
}
