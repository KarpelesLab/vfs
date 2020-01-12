package memkv

import (
	"os"
	"strings"
	"sync"

	"github.com/KarpelesLab/vfs"
	"github.com/petar/GoLLRB/llrb"
)

type entry struct {
	key  string
	data []byte
}

type entryStart struct{}

type memKv struct {
	data *llrb.LLRB
	lk   sync.RWMutex
}

func New() vfs.Keyval {
	return &memKv{
		data: llrb.New(),
	}
}

func (e *memKv) Get(key string) (vfs.KVEntry, error) {
	e.lk.RLock()
	f := e.data.Get(&entry{key: key})
	e.lk.RUnlock()

	if f == nil {
		return nil, os.ErrNotExist
	}

	return f.(*entry), nil
}

func (e *memKv) Put(key string, value vfs.KVEntry) error {
	data, err := value.Data()
	if err != nil {
		return err
	}

	e.lk.Lock()
	e.data.ReplaceOrInsert(&entry{key: key, data: data})
	e.lk.Unlock()

	return nil
}

func (e *memKv) Delete(key string) error {
	e.lk.Lock()
	e.data.Delete(&entry{key: key})
	e.lk.Unlock()
	return nil
}

func (e memKv) List(prefix string, callback func(key string, value vfs.KVEntry) (bool, error)) error {
	var list []*entry

	// make a copy so callback can use functions that use a lock
	e.lk.RLock()
	if prefix == "" {
		e.data.AscendGreaterOrEqual(entryStart{}, func(i llrb.Item) bool {
			list = append(list, i.(*entry))
			return true
		})
	} else {
		e.data.AscendGreaterOrEqual(&entry{key: prefix}, func(i llrb.Item) bool {
			item := i.(*entry)
			if !strings.HasPrefix(item.key, prefix) {
				return false
			}
			list = append(list, item)
			return true
		})
	}
	e.lk.RUnlock()

	for _, e := range list {
		cont, err := callback(e.key, e)
		if err != nil {
			return err
		}
		if !cont {
			return nil
		}
	}
	return nil
}

func (d *entry) Data() ([]byte, error) {
	return d.data, nil
}

func (d *entry) Less(than llrb.Item) bool {
	switch e := than.(type) {
	case *entry:
		return strings.Compare(d.key, e.key) < 0
	case entryStart:
		return false // always higher than start
	}

	return false // ???
}

func (e entryStart) Less(than llrb.Item) bool {
	return true // start is always less
}
