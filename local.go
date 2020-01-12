package vfs

import (
	"os"
	"path/filepath"
)

type LocalFS struct {
	root string
}

type LocalFile os.File

// NewLocal creates a new local filesystem with root as root point. Note that
// the root argument format depends on filesystem.
func NewLocal(root string) (FileSystem, error) {
	return &LocalFS{root}, nil
}

func (l *LocalFS) Open(name string) (File, error) {
	p := filepath.Join(l.root, filepath.FromSlash(name))

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	return (*LocalFile)(f), nil
}

func (f *LocalFile) Close() error {
	return (*os.File)(f).Close()
}

func (f *LocalFile) Read(p []byte) (int, error) {
	return (*os.File)(f).Read(p)
}

func (f *LocalFile) Readdir(n int) ([]os.FileInfo, error) {
	return (*os.File)(f).Readdir(n)
}

func (f *LocalFile) Seek(offset int64, whence int) (int64, error) {
	return (*os.File)(f).Seek(offset, whence)
}

func (f *LocalFile) Stat() (os.FileInfo, error) {
	return (*os.File)(f).Stat()
}
