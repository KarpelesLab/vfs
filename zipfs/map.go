package zipfs

import (
	"archive/zip"
	"os"
	"path"
	"strings"

	"github.com/KarpelesLab/vfs"
	"github.com/KarpelesLab/vfs/vdirfs"
)

type mappedZip struct {
	z *zip.Reader
	i map[string]*zip.File
}

func newZipMap(z *zip.Reader, addVdir bool) (vfs.FileSystem, error) {
	m := &mappedZip{
		z: z,
		i: make(map[string]*zip.File),
	}

	for _, f := range z.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}
		m.i[f.Name] = f
	}

	if addVdir {
		res, err := vdirfs.New(m)
		if err != nil {
			return nil, err
		}

		for f, _ := range m.i {
			res.AddPath(f)
		}

		return res, nil
	}

	return m, nil
}

func (z *mappedZip) get(fn string) (*zip.File, error) {
	fn = strings.TrimLeft(fn, "/")
	if f, ok := z.i[fn]; ok {
		return f, nil
	}

	return nil, os.ErrNotExist
}

func (z *mappedZip) Open(name string) (vfs.File, error) {
	f, err := z.get(name)
	if err != nil {
		return nil, err
	}

	return &file{zf: f}, nil
}

func (z *mappedZip) OpenFile(path string, flag int, perm os.FileMode) (vfs.File, error) {
	if flag != os.O_RDONLY {
		return nil, os.ErrPermission
	}

	return z.Open(path)
}

func (z *mappedZip) Lstat(name string) (os.FileInfo, error) {
	f, err := z.get(name)
	if err != nil {
		return nil, err
	}

	// NewStat(name string, size int64, mode os.FileMode, modTime time.Time, sys interface{})
	st := vfs.NewStat(path.Base(f.Name), int64(f.UncompressedSize64), 0755, f.Modified, f)

	return st, nil
}

func (z *mappedZip) Stat(path string) (os.FileInfo, error) {
	return z.Lstat(path)
}

func (z *mappedZip) Mkdir(path string, perm os.FileMode) error {
	return os.ErrPermission
}

func (z *mappedZip) Remove(path string) error {
	return os.ErrPermission
}
