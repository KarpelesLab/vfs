package boltkv

import (
	"io"
	"os"

	"github.com/boltdb/bolt"
)

type boltVal struct {
	kv  *boltKv
	key string
}

func (bv *boltVal) Data() (res []byte, err error) {
	err = bv.kv.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bv.kv.bucket)
		if b == nil {
			return os.ErrNotExist
		}

		v := b.Get([]byte(bv.key))
		if v == nil {
			return os.ErrNotExist
		}

		res = make([]byte, len(v))
		copy(res, v)

		return nil
	})
	return
}

// io.WriterTo interface
func (bv *boltVal) WriteTo(w io.Writer) (n int64, err error) {
	err = bv.kv.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bv.kv.bucket)
		if b == nil {
			return os.ErrNotExist
		}

		v := b.Get([]byte(bv.key))
		if v == nil {
			return os.ErrNotExist
		}

		n2, err := w.Write(v)
		n = int64(n2)

		return err
	})
	return
}
