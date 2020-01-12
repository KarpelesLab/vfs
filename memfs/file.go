package memfs

import (
	"io"
	"os"
	"time"
)

type memFile struct {
	buf     []byte
	mode    os.FileMode
	modTime time.Time
}

func (m *memFile) ReadAt(b []byte, off int64) (int, error) {
	if off >= int64(len(m.buf)) {
		return 0, io.EOF
	}
	n := copy(b, m.buf[off:])
	return n, nil
}
