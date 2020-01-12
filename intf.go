package vfs

import (
	"io"
	"os"
)

type FileSystem interface {
	Open(name string) (File, error)
	OpenFile(path string, flag int, perm os.FileMode) (File, error)
}

type File interface {
	io.Closer
	io.Reader
	io.Writer
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}
