package memfs

import (
	"io"
	"os"
	"sync"
	"time"
)

type memFile struct {
	buf     []byte
	mode    os.FileMode
	modTime time.Time
	lk      sync.RWMutex
}

func (m *memFile) ReadAt(b []byte, off int64) (int, error) {
	if off >= int64(len(m.buf)) {
		return 0, io.EOF
	}
	n := copy(b, m.buf[off:])
	return n, nil
}

func (m *memFile) Size() int64 {
	m.lk.RLock()
	sz := len(m.buf)
	m.lk.RUnlock()

	return int64(sz)
}

func (m *memFile) Mode() os.FileMode {
	return m.mode
}

func (m *memFile) ModTime() time.Time {
	return m.modTime
}
