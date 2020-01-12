package memfs

import (
	"os"
	"time"
)

type node interface {
	ReadAt(b []byte, off int64) (n int, err error)
	Size() int64
	Mode() os.FileMode
	ModTime() time.Time
}
