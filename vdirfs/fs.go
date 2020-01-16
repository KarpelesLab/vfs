package vdirfs

import (
	"os"
	"path"
	"strings"

	"github.com/KarpelesLab/vfs"
)

type FS struct {
	parent vfs.FileSystem
	root   *dir
}

func New(fs vfs.FileSystem) (*FS, error) {
	f := &FS{
		parent: fs,
	}

	f.root = &dir{
		path:     "",
		children: make(map[string]*dir),
		fs:       f,
	}
	return f, nil
}

func (f *FS) AddPath(name string) error {
	name = strings.TrimLeft(name, "/")

	d, err := f.root.getDir(path.Dir(name), true)
	if err != nil {
		return err
	}
	d.children[path.Base(name)] = nil
	return nil
}

func (f *FS) Open(name string) (vfs.File, error) {
	d, err := f.root.getDir(name, false)
	if err != nil {
		return nil, err
	}

	if d != nil {
		return d, nil
	}

	return f.parent.Open(name)
}

func (f *FS) OpenFile(path string, flag int, perm os.FileMode) (vfs.File, error) {
	if flag != os.O_RDONLY {
		return f.parent.OpenFile(path, flag, perm)
	}

	d, err := f.root.getDir(path, false)
	if err != nil {
		return nil, err
	}
	if d != nil {
		return d, nil
	}

	return f.parent.OpenFile(path, flag, perm)
}

func (f *FS) Lstat(path string) (os.FileInfo, error) {
	if d, err := f.root.getDir(path, false); err == nil && d != nil {
		return d.Stat()
	}
	return f.parent.Lstat(path)
}

func (f *FS) Stat(path string) (os.FileInfo, error) {
	if d, err := f.root.getDir(path, false); err == nil && d != nil {
		return d.Stat()
	}
	return f.parent.Stat(path)
}

func (f *FS) Mkdir(path string, perm os.FileMode) error {
	return f.parent.Mkdir(path, perm)
}

func (f *FS) Remove(path string) error {
	return f.parent.Remove(path)
}
