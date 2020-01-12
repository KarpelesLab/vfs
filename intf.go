package vfs

import (
	"io"
	"os"
)

type FileSystem interface {
	Open(name string) (File, error)
	OpenFile(path string, flag int, perm os.FileMode) (File, error)
	// Lstat returns the os.FileInfo for the given path, without
	// following symlinks.
	Lstat(path string) (os.FileInfo, error)
	// Stat returns the os.FileInfo for the given path, following
	// symlinks.
	Stat(path string) (os.FileInfo, error)
	// Mkdir creates a directory at the given path. If the directory
	// already exists or its parent directory does not exist or
	// the permissions don't allow it, an error will be returned. See
	// also the shorthand function MkdirAll.
	Mkdir(path string, perm os.FileMode) error
	// Remove removes the item at the given path. If the path does
	// not exists or the permissions don't allow removing it or it's
	// a non-empty directory, an error will be returned. See also
	// the shorthand function RemoveAll.
	Remove(path string) error
}

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Writer
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}
