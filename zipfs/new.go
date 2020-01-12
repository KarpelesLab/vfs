package zipfs

import (
	"archive/zip"

	"github.com/KarpelesLab/vfs"
)

func New(z *zip.Reader, idx Index) (vfs.FileSystem, error) {
	switch idx {
	case IndexMap:
		return newZipMap(z, false)
	case IndexFull:
		return newZipMap(z, true)
	default:
		return newZip(z)
	}
}
