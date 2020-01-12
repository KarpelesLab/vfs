package memfs

import (
	"io"
	"os"
	"time"
)

type node interface {
	io.ReaderAt
	io.WriterAt
	Size() int64
	Mode() os.FileMode
	ModTime() time.Time
}
