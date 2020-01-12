package vfs

import "errors"

var (
	ErrNotDirectory = errors.New("not a directory")
	ErrIsDirectory  = errors.New("is a directory")
	ErrNotEmpty     = errors.New("not empty")
)
