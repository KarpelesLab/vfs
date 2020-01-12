package memfs

import (
	"os"
	"time"
)

type memStat struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     interface{}
}
