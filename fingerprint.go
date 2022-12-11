package vfs

import (
	"bytes"
	"io/fs"
	"strconv"
)

// Fingerprintable is a type of object returned by Sys() that allows obtaining
// a string identifying a file at a current status. In case the file is
// modified the fingerprint should be different, and its value can be used in
// HTTP etag for example. If the backend storage allows it, the value should be
// a hash (ie. CRC32 for zipfs), and if not, size and timestamp may do.
type Fingerprintable interface {
	Fingerprint() (string, error)
}

// Fingerprint will return a fingerprint for a given file, silently falling
// back on using size+timestamp in case of backend failure.
func Fingerprint(f fs.FileInfo) string {
	if fp, ok := f.Sys().(Fingerprintable); ok {
		if res, err := fp.Fingerprint(); err == nil && res != "" {
			return res
		}
	}

	// fallback
	t := f.ModTime()

	buf := &bytes.Buffer{}
	buf.WriteString(strconv.FormatInt(f.Size(), 36))
	buf.WriteByte('.')
	buf.WriteString(strconv.FormatInt(t.Unix(), 36))
	buf.WriteByte('.')
	buf.WriteString(strconv.FormatInt(int64(t.Nanosecond()), 36))

	return buf.String()
}
