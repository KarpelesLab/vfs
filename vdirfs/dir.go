package vdirfs

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/KarpelesLab/vfs"
)

type dir struct {
	fs       *FS
	path     string
	children map[string]*dir
}

func (d *dir) getDir(name string, create bool) (*dir, error) {
	pos := strings.IndexByte(name, '/')
	if pos == -1 {
		sd, ok := d.children[name]
		if ok {
			return sd, nil
		}
		if !create {
			return nil, os.ErrNotExist
		}
		sd = &dir{
			path:     path.Join(d.path, name),
			children: make(map[string]*dir),
		}
		d.children[name] = sd
		return sd, nil
	}

	next := name[pos+1:]
	name = name[:pos]

	sd, ok := d.children[name]
	if ok {
		return sd.getDir(next, create)
	}

	if !create {
		return nil, os.ErrNotExist
	}

	sd = &dir{
		fs:       d.fs,
		path:     path.Join(d.path, name),
		children: make(map[string]*dir),
	}
	d.children[name] = sd

	return sd.getDir(next, create)
}

func (d *dir) Close() error {
	return nil
}

func (d *dir) Readdir(n int) ([]os.FileInfo, error) {
	var res []os.FileInfo
	now := time.Now()
	for n, sd := range d.children {
		if sd == nil {
			// need to stat
			st, err := d.fs.parent.Lstat(path.Join(d.path, n))
			if err != nil {
				return nil, err
			}
			res = append(res, st)
		} else {
			st := vfs.NewStat(path.Base(sd.path), 0, 0755, now, sd)
			res = append(res, st)
		}
	}

	return res, nil
}

func (d *dir) Read(b []byte) (int, error) {
	return 0, vfs.ErrIsDirectory
}

func (d *dir) ReadAt(b []byte, off int64) (int, error) {
	return 0, vfs.ErrIsDirectory
}

func (d *dir) Write(b []byte) (int, error) {
	return 0, vfs.ErrIsDirectory
}

func (d *dir) WriteAt(b []byte, off int64) (int, error) {
	return 0, vfs.ErrIsDirectory
}

func (d *dir) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (d *dir) Stat() (os.FileInfo, error) {
	return vfs.NewStat(path.Base(d.path), 0, 0755, time.Now(), d), nil
}
