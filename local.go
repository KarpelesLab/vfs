package vfs

import (
	"os"
	"path"
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

func (l *LocalFS) doPath(p string) string {
	p = path.Clean(p)
	if l.root == "" {
		return filepath.FromSlash(p)
	}
	return filepath.Join(l.root, filepath.FromSlash(p))
}

func (l *LocalFS) Open(name string) (File, error) {
	f, err := os.Open(l.doPath(name))
	if err != nil {
		return nil, err
	}

	return (*LocalFile)(f), nil
}

func (l *LocalFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	f, err := os.OpenFile(l.doPath(name), flag, perm)
	if err != nil {
		return nil, err
	}

	return (*LocalFile)(f), nil
}

func (l *LocalFS) Chroot(name string) (FileSystem, error) {
	p := l.doPath(name)

	p, err := filepath.Abs(p)
	if err != nil {
		return nil, err
	}

	p, err = filepath.EvalSymlinks(p)
	if err != nil {
		return nil, err
	}

	// Check if p is an existing directory
	st, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return nil, ErrNotDirectory
	}

	return &LocalFS{p}, nil
}

func (f *LocalFile) Close() error {
	return (*os.File)(f).Close()
}

func (f *LocalFile) Read(p []byte) (int, error) {
	return (*os.File)(f).Read(p)
}

func (f *LocalFile) Write(p []byte) (int, error) {
	return (*os.File)(f).Write(p)
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
