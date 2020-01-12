package memfs

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/KarpelesLab/vfs"
)

type memDir struct {
	children map[string]node
	mode     os.FileMode
	modTime  time.Time
}

func New() (vfs.FileSystem, error) {
	// make a root
	root := &memDir{
		children: make(map[string]node),
		mode:     os.ModeDir | 0755,
		modTime:  time.Now(),
	}

	return root, nil
}

func (m *memDir) access(name string) (node, error) {
	name = path.Clean(name)
	for len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	pos := strings.IndexByte(name, '/')
	if pos == -1 {
		v, ok := m.children[name]
		if !ok {
			return nil, os.ErrNotExist
		}
		return v, nil
	}

	sub := name[:pos]
	name = name[pos+1:]

	v, ok := m.children[sub]
	if !ok {
		return nil, os.ErrNotExist
	}

	vD, ok := v.(*memDir)
	if !ok {
		return nil, vfs.ErrNotDirectory
	}

	return vD.access(name)
}

func (m *memDir) Open(name string) (vfs.File, error) {
	a, err := m.access(name)
	if err != nil {
		return nil, err
	}

	return &memOpen{node: a, flag: os.O_RDONLY, name: path.Base(name)}, nil
}

func (m *memDir) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	if flag&os.O_CREATE != os.O_CREATE {
		// request for an existing file, easy.
		a, err := m.access(name)
		if err != nil {
			return nil, err
		}
		return &memOpen{node: a, flag: flag, name: path.Base(name)}, nil
	}

	a, err := m.access(path.Dir(name))
	if err != nil {
		return nil, err
	}
	name = path.Base(name)
	if name == "/" || name == "." || name == ".." {
		return nil, os.ErrInvalid
	}
	dir, ok := a.(*memDir)
	if !ok {
		return nil, vfs.ErrNotDirectory
	}

	_, hasOld := dir.children[name]

	if flag&os.O_EXCL == os.O_EXCL {
		// fail if exists
		if hasOld {
			return nil, os.ErrExist
		}
	}

	// create file
	newFile := &memFile{
		modTime: time.Now(),
	}
	dir.children[name] = newFile

	return &memOpen{node: newFile, flag: flag, name: name}, nil
}

func (m *memDir) Lstat(name string) (os.FileInfo, error) {
	// no symlinks
	return m.Stat(name)
}

func (m *memDir) Stat(name string) (os.FileInfo, error) {
	a, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	return a.Stat()
}

func (m *memDir) Mkdir(name string, perm os.FileMode) error {
	a, err := m.access(path.Dir(name))
	if err != nil {
		return err
	}
	name = path.Base(name)
	if name == "/" || name == "." || name == ".." {
		return os.ErrInvalid
	}

	dir, ok := a.(*memDir)
	if !ok {
		return vfs.ErrNotDirectory
	}

	_, hasOld := dir.children[name]
	if hasOld {
		return os.ErrExist
	}

	newDir := &memDir{
		children: make(map[string]node),
		mode:     os.ModeDir | perm,
		modTime:  time.Now(),
	}

	dir.children[name] = newDir
	return nil
}

func (m *memDir) Remove(name string) error {
	a, err := m.access(path.Dir(name))
	if err != nil {
		return err
	}
	name = path.Base(name)
	if name == "/" || name == "." || name == ".." {
		return os.ErrInvalid
	}

	dir, ok := a.(*memDir)
	if !ok {
		return vfs.ErrNotDirectory
	}

	if obj, ok := dir.children[name]; ok {
		if dir, ok := obj.(*memDir); ok {
			if len(dir.children) != 0 {
				return vfs.ErrNotEmpty
			}
		}
	}

	delete(dir.children, name)
	return nil
}

func (m *memDir) ReadAt(b []byte, off int64) (int, error) {
	return 0, vfs.ErrIsDirectory
}