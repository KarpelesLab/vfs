package vfs

import (
	"fmt"
	"io/fs"
	"strings"
)

func MkdirAll(fsbase FileSystem, path string, perm fs.FileMode) error {
	if len(path) == 0 {
		// path cannot be empty and must be absolute
		return fs.ErrInvalid
	}

	cur := ""

	for path != "" {
		if cur != "" {
			cur = cur + "/"
		}

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

		st, err := fsbase.Lstat(cur)
		if err == nil {
			if !st.IsDir() {
				return fmt.Errorf("%s exists and is %w", cur, ErrNotDirectory)
			}
		} else {
			if err = fsbase.Mkdir(cur, perm); err != nil {
				return fmt.Errorf("failed to mkdir %s: %w", cur, err)
			}
		}
	}

	return nil
}

func ReadDir(fsbase FileSystem, path string) ([]fs.DirEntry, error) {
	return fs.ReadDir(fsbase, path)
}
