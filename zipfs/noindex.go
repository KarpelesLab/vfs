package zipfs

import (
	"archive/zip"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/KarpelesLab/vfs"
)

type noIndexZip struct {
	z *zip.Reader
}

func newZip(z *zip.Reader) (vfs.FileSystem, error) {
	return &noIndexZip{z}, nil
}

func (z *noIndexZip) get(fn string) (*zip.File, error) {
	fn = strings.TrimLeft(fn, "/")
	for _, f := range z.z.File {
		if f.Name == fn {
			return f, nil
		}
	}

	return nil, os.ErrNotExist
}

func (z *noIndexZip) Open(name string) (fs.File, error) {
	f, err := z.get(name)
	if err != nil {
		return nil, err
	}

	return &file{zf: f}, nil
}

func (z *noIndexZip) OpenFile(path string, flag int, perm os.FileMode) (fs.File, error) {
	if flag != os.O_RDONLY {
		return nil, os.ErrPermission
	}

	return z.Open(path)
}

func (z *noIndexZip) Lstat(name string) (os.FileInfo, error) {
	f, err := z.get(name)
	if err != nil {
		return nil, err
	}

	// NewStat(name string, size int64, mode os.FileMode, modTime time.Time, sys interface{})
	st := vfs.NewStat(path.Base(f.Name), int64(f.UncompressedSize64), 0755, f.Modified, f)

	return st, nil
}

func (z *noIndexZip) Stat(path string) (os.FileInfo, error) {
	return z.Lstat(path)
}

func (z *noIndexZip) Mkdir(path string, perm os.FileMode) error {
	return os.ErrPermission
}

func (z *noIndexZip) Remove(path string) error {
	return os.ErrPermission
}
