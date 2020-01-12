package prefixkv

import (
	"github.com/KarpelesLab/vfs"
)

type prefixKv struct {
	prefix string
	parent vfs.Keyval
}

func New(parent vfs.Keyval, prefix string) vfs.Keyval {
	return &prefixKv{
		prefix: prefix,
		parent: parent,
	}
}

func (e *prefixKv) Get(key string) (vfs.KVEntry, error) {
	return e.parent.Get(e.prefix + key)
}

func (e *prefixKv) Put(key string, value vfs.KVEntry) error {
	return e.parent.Put(e.prefix+key, value)
}

func (e *prefixKv) Delete(key string) error {
	return e.parent.Delete(e.prefix + key)
}

func (e prefixKv) List(prefix string, callback func(key string, value vfs.KVEntry) (bool, error)) error {
	return e.parent.List(e.prefix+prefix, callback)
}
