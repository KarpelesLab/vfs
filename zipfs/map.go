package zipfs

import (
	"archive/zip"
	"errors"

	"github.com/KarpelesLab/vfs"
)

func newZipMap(z *zip.Reader) (vfs.FileSystem, error) {
	return nil, errors.New("TODO")
}
