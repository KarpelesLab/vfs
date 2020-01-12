package boltkv

import (
	"bytes"
	"os"

	"github.com/KarpelesLab/vfs"
	"github.com/boltdb/bolt"
)

type boltKv struct {
	db     *bolt.DB
	bucket []byte
}

func New(db *bolt.DB, bucket []byte) vfs.Keyval {
	res := &boltKv{
		db:     db,
		bucket: bucket,
	}

	return res
}

func (d *boltKv) Get(key string) (res vfs.KVEntry, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		if b == nil {
			return os.ErrNotExist
		}

		if b.Get([]byte(key)) == nil {
			return os.ErrNotExist
		}

		res = &boltVal{
			kv:  d,
			key: key,
		}
		return nil
	})
	return
}

func (d *boltKv) Put(key string, value vfs.KVEntry) error {
	val, err := value.Data()
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(d.bucket)
		if err != nil {
			return err
		}

		return b.Put([]byte(key), val)
	})
}

func (d *boltKv) Delete(key string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		if b == nil {
			return nil
		}

		return b.Delete([]byte(key))
	})
}

func (d *boltKv) List(prefix string, callback func(key string, value vfs.KVEntry) (cont bool, err error)) error {
	var list []string
	pfx := []byte(prefix)

	// generate list
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		if b == nil {
			// empty list
			return nil
		}

		c := b.Cursor()
		k, _ := c.Seek(pfx)
		if k == nil {
			return nil
		}

		for bytes.HasPrefix(k, pfx) {
			list = append(list, string(k))

			k, _ = c.Next()
			if k == nil {
				break
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, k := range list {
		cont, err := callback(k, &boltVal{kv: d, key: k})
		if err != nil {
			return err
		}
		if !cont {
			break
		}
	}
	return nil
}
