package chrootfs

import (
	"os"
	"path"

	"github.com/KarpelesLab/vfs"
)

type chrooter interface {
	Chroot(name string) (vfs.FileSystem, error)
}

type chrootFs struct {
	parent vfs.FileSystem
	root   string
}

func New(fs vfs.FileSystem, name string) (vfs.FileSystem, error) {
	if chfs, ok := fs.(chrooter); ok {
		return chfs.Chroot(name)
	}

	name = path.Clean(name)

	return &chrootFs{
		parent: fs,
		root:   name,
	}, nil
}

func (c *chrootFs) doPath(name string) (string, error) {
	if len(name) == 0 || name[0] != '/' {
		return name, os.ErrInvalid
	}

	name = path.Clean(name)
	name = c.root + name

	return name, nil
}

func (c *chrootFs) Open(name string) (vfs.File, error) {
	name, err := c.doPath(name)
	if err != nil {
		return nil, err
	}

	return c.parent.Open(name)
}

func (c *chrootFs) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	name, err := c.doPath(name)
	if err != nil {
		return nil, err
	}

	return c.parent.OpenFile(name, flag, perm)
}

func (c *chrootFs) Lstat(name string) (os.FileInfo, error) {
	name, err := c.doPath(name)
	if err != nil {
		return nil, err
	}

	return c.parent.Lstat(name)
}

func (c *chrootFs) Stat(name string) (os.FileInfo, error) {
	name, err := c.doPath(name)
	if err != nil {
		return nil, err
	}

	return c.parent.Stat(name)
}

func (c *chrootFs) Mkdir(name string, perm os.FileMode) error {
	name, err := c.doPath(name)
	if err != nil {
		return err
	}

	return c.parent.Mkdir(name, perm)
}

func (c *chrootFs) Remove(name string) error {
	name, err := c.doPath(name)
	if err != nil {
		return err
	}

	return c.parent.Remove(name)
}
