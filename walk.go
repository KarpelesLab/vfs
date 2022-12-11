package vfs

import (
	"io/fs"
	"path"
)

func walk(fsbase FileSystem, p string, walkFn fs.WalkDirFunc, info fs.DirEntry, err error) error {
	if err != nil {
		return walkFn(p, info, err)
	}
	err = walkFn(p, info, nil)
	if !info.IsDir() {
		return err
	}
	if err == fs.SkipDir {
		return nil
	}
	// note: ReadDir returns results sorted by name
	infos, err := ReadDir(fsbase, p)
	if err != nil {
		return err
	}
	for _, info := range infos {
		name := info.Name()
		if name == "." || name == ".." {
			continue
		}
		if err := walk(fsbase, path.Join(p, info.Name()), walkFn, info, nil); err != nil {
			return err
		}
	}
	return nil
}

// Walk is the equivalent of filepath.Walk
func Walk(fsbase FileSystem, path string, walkFn fs.WalkDirFunc) error {
	info, err := fsbase.Lstat(path)
	return walk(fsbase, path, walkFn, fs.FileInfoToDirEntry(info), err)
}
