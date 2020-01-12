package memfs

import (
	"io"
	"os"

	"github.com/KarpelesLab/vfs"
)

type memOpen struct {
	node   node
	flag   int
	name   string
	offset int64
}

func (m *memOpen) Close() error {
	return nil
}

func (m *memOpen) Read(b []byte) (int, error) {
	n, err := m.node.ReadAt(b, m.offset)
	if n > 0 {
		m.offset += int64(n)
	}
	return n, err
}

func (m *memOpen) ReadAt(b []byte, offset int64) (int, error) {
	return m.node.ReadAt(b, offset)
}

func (m *memOpen) Write(b []byte) (int, error) {
	n, err := m.node.WriteAt(b, m.offset)
	if n > 0 {
		m.offset += int64(n)
	}
	return n, err
}

func (m *memOpen) WriteAt(b []byte, offset int64) (int, error) {
	return m.node.WriteAt(b, offset)
}

func (m *memOpen) Readdir(n int) ([]os.FileInfo, error) {
	dir, ok := m.node.(*memDir)
	if !ok {
		return nil, vfs.ErrNotDirectory
	}

	return dir.Readdir()
}

func (m *memOpen) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		m.offset = offset
		return m.offset, nil
	case io.SeekCurrent:
		m.offset += offset
		return m.offset, nil
	case io.SeekEnd:
		m.offset = m.node.Size() + offset
		return m.offset, nil
	default:
		return m.offset, os.ErrInvalid
	}
}

func (m *memOpen) Stat() (os.FileInfo, error) {
	return vfs.NewStat(m.name, m.node.Size(), m.node.Mode(), m.node.ModTime(), m.node), nil
}
