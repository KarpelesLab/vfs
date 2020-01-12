package memfs

type node interface {
	ReadAt(b []byte, off int64) (n int, err error)
}
