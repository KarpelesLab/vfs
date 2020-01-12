package vfs

import (
	"io"
	"os"
)

type Filesystem interface {
	Open(name string) (File, error)
}

type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}
