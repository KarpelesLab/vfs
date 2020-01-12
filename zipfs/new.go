package zipfs

import (
	"archive/zip"

	"github.com/KarpelesLab/vfs"
)

func New(z *zip.Reader, idx Index) (vfs.FileSystem, error) {
	switch idx {
	case IndexMap:
		return newZipMap(z)
	case IndexFull:
		return newZipFull(z)
	default:
		return newZip(z)
	}
}
