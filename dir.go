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

	for {
		cur = cur + "/"
		pos := strings.IndexByte(path, '/')
		if pos == 0 {
			path = path[1:]
			continue
		}
		if pos > 0 {
			cur = cur + path[:pos]
			path = path[pos+1:]
		}

		st, err := fs.Lstat(path)
		if err == nil {
			if !st.IsDir() {
				return fmt.Errorf("%s exists and is %w", path, ErrNotDirectory)
			}
		} else {
			if err = fs.Mkdir(path, perm); err != nil {
				return fmt.Errorf("failed to mkdir %s: %w", path, err)
			}
		}
	}
}
