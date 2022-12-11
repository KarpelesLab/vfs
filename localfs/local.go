package localfs

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/KarpelesLab/vfs"
)

type localFS struct {
	root string
}

type localFile os.File

// New creates a new local filesystem with root as root point. Note that
// the root argument format depends on filesystem.
func New(p string) (vfs.FileSystem, error) {
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
		return nil, vfs.ErrNotDirectory
	}

	return &localFS{p}, nil
}

func (l *localFS) doPath(p string) string {
	p = path.Clean(p)
	if l.root == "" {
		return filepath.FromSlash(p)
	}
	return filepath.Join(l.root, filepath.FromSlash(p))
}

func (l *localFS) Open(name string) (vfs.File, error) {
	f, err := os.Open(l.doPath(name))
	if err != nil {
		return nil, err
	}

	return (*localFile)(f), nil
}

func (l *localFS) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	f, err := os.OpenFile(l.doPath(name), flag, perm)
	if err != nil {
		return nil, err
	}

	return (*localFile)(f), nil
}

func (l *localFS) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(l.doPath(name))
}

func (l *localFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(l.doPath(name))
}

func (l *localFS) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(l.doPath(name), perm)
}

func (l *localFS) Remove(name string) error {
	return os.Remove(l.doPath(name))
}

func (l *localFS) Chroot(name string) (vfs.FileSystem, error) {
	p := l.doPath(name)
	return New(p)
}

func (f *localFile) Close() error {
	return (*os.File)(f).Close()
}

func (f *localFile) Read(p []byte) (int, error) {
	return (*os.File)(f).Read(p)
}

func (f *localFile) ReadAt(p []byte, offset int64) (int, error) {
	return (*os.File)(f).ReadAt(p, offset)
}

func (f *localFile) Write(p []byte) (int, error) {
	return (*os.File)(f).Write(p)
}

func (f *localFile) WriteAt(p []byte, offset int64) (int, error) {
	return (*os.File)(f).WriteAt(p, offset)
}

func (f *localFile) ReadDir(n int) ([]fs.DirEntry, error) {
	return (*os.File)(f).ReadDir(n)
}

func (f *localFile) Seek(offset int64, whence int) (int64, error) {
	return (*os.File)(f).Seek(offset, whence)
}

func (f *localFile) Stat() (os.FileInfo, error) {
	return (*os.File)(f).Stat()
}
