package zipfs

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"sync"

	"github.com/KarpelesLab/vfs"
)

type file struct {
	zf *zip.File

	rd   io.ReadCloser
	lk   sync.Mutex
	pos  int64
	rpos int64
}

func (f *file) Close() error {
	f.lk.Lock()
	if f.rd == nil {
		f.lk.Unlock()
		return nil
	}

	err := f.rd.Close()
	f.lk.Unlock()

	return err
}

func (f *file) Read(b []byte) (int, error) {
	f.lk.Lock()
	n, err := f.read(b)
	f.lk.Unlock()
	return n, err
}

func (f *file) read(b []byte) (int, error) {
	if uint64(f.pos) >= f.zf.UncompressedSize64 {
		return 0, io.EOF
	}

	if f.pos < f.rpos && f.rd != nil {
		// close file & resume from start
		f.rd.Close()
		f.rd = nil
	}

	if f.rd == nil {
		rd, err := f.zf.Open()
		if err != nil {
			return 0, err
		}
		f.rd = rd
		f.rpos = 0
	}

	if f.pos > f.rpos {
		// need to move forward in file
		offt := f.pos - f.rpos

		var buf []byte
		if 8192 > offt {
			buf = make([]byte, offt)
		} else {
			buf = make([]byte, 8192)
		}

		for offt > 0 {
			if int64(len(buf)) > offt {
				buf = buf[:offt]
			}

			n, err := f.rd.Read(buf)
			if err != nil {
				return 0, err
			}
			offt -= int64(n)
		}

		f.rpos = f.pos
	}

	n, err := f.rd.Read(b)

	if n > 0 {
		f.pos += int64(n)
		f.rpos += int64(n)
	}

	return n, err
}

func (f *file) ReadAt(b []byte, pos int64) (int, error) {
	f.lk.Lock()
	curpos := f.pos
	f.pos = pos
	n, err := f.read(b)
	f.pos = curpos
	f.lk.Unlock()
	return n, err
}

func (f *file) Readdir(n int) ([]os.FileInfo, error) {
	return nil, os.ErrInvalid
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.pos = offset
		return f.pos, nil
	case io.SeekCurrent:
		f.pos += offset
		return f.pos, nil
	case io.SeekEnd:
		f.pos = int64(f.zf.UncompressedSize64) + offset
		return f.pos, nil
	default:
		return f.pos, os.ErrInvalid
	}
}

func (f *file) Stat() (os.FileInfo, error) {
	st := vfs.NewStat(path.Base(f.zf.Name), int64(f.zf.UncompressedSize64), 0755, f.zf.Modified, f.zf)

	return st, nil
}

func (f *file) Write(b []byte) (int, error) {
	return 0, os.ErrPermission
}

func (f *file) WriteAt(b []byte, pos int64) (int, error) {
	return 0, os.ErrPermission
}
