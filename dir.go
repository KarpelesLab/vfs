package vfs

import (
	"fmt"
	"os"
	"strings"
)

func MkdirAll(fs FileSystem, path string, perm os.FileMode) error {
	if len(path) == 0 || path[0] != '/' {
		// path cannot be empty and must be absolute
		return os.ErrInvalid
	}

	cur := ""
	path = path[1:]

	for path != "" {
		cur = cur + "/"
		pos := strings.IndexByte(path, '/')
		if pos == 0 {
			path = path[1:]
			continue
		}
		if pos > 0 {
			cur = cur + path[:pos]
			path = path[pos+1:]
		} else {
			cur = cur + path
			path = ""
		}

		st, err := fs.Lstat(cur)
		if err == nil {
			if !st.IsDir() {
				return fmt.Errorf("%s exists and is %w", cur, ErrNotDirectory)
			}
		} else {
			if err = fs.Mkdir(cur, perm); err != nil {
				return fmt.Errorf("failed to mkdir %s: %w", cur, err)
			}
		}
	}

	return nil
}
