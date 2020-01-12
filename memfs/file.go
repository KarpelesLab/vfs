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

func (m *memFile) WriteAt(b []byte, off int64) (int, error) {
	if off > int64(len(m.buf)) {
		if off <= int64(cap(m.buf)) {
			// already has the capacity, simple resize
			m.buf = m.buf[:off]
		} else {
			// need to add zeroes (pre-allocate space for the new data)
			newBuf := make([]byte, off, off+int64(len(b)))
			copy(newBuf, m.buf)
			m.buf = newBuf
		}
	}

	if off == int64(len(m.buf)) {
		m.buf = append(m.buf, b...)
		return len(b), nil
	}

	if off+int64(len(b)) > int64(len(m.buf)) {
		if off+int64(len(b)) <= int64(cap(m.buf)) {
			// already has the capacity, simple resize
			m.buf = m.buf[:off+int64(len(b))]
		} else {
			// need to resize
			newBuf := make([]byte, off+int64(len(b)))
			copy(newBuf, m.buf[:off])
			m.buf = newBuf
		}
	}
	// fully fits
	n := copy(m.buf[off:], b)
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
