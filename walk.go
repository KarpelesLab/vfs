package vfs

import (
	"os"
	"path"
	"path/filepath"
)

func walk(fs FileSystem, p string, walkFn filepath.WalkFunc, info os.FileInfo, err error) error {
	if err != nil {
		return walkFn(p, info, err)
	}
	err = walkFn(p, info, nil)
	if !info.IsDir() {
		return err
	}
	if err == filepath.SkipDir {
		return nil
	}
	// note: ReadDir returns results sorted by name
	infos, err := ReadDir(fs, p)
	if err != nil {
		return err
	}
	for _, info := range infos {
		name := info.Name()
		if name == "." || name == ".." {
			continue
		}
		if err := walk(fs, path.Join(p, info.Name()), walkFn, info, nil); err != nil {
			return err
		}
	}
	return nil
}

// Walk is the equivalent of filepath.Walk
func Walk(fs FileSystem, path string, walkFn filepath.WalkFunc) error {
	info, err := fs.Lstat(path)
	return walk(fs, path, walkFn, info, err)
}
