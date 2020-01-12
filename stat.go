package vfs

import (
	"os"
	"time"
)

type vStat struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     interface{}
}

func NewStat(name string, size int64, mode os.FileMode, modTime time.Time, sys interface{}) os.FileInfo {
	return &vStat{
		name:    name,
		size:    size,
		mode:    mode,
		modTime: modTime,
		sys:     sys,
	}
}

func (v *vStat) Name() string {
	return v.name
}

func (v *vStat) Size() int64 {
	return v.size
}

func (v *vStat) IsDir() bool {
	return v.mode.IsDir()
}

func (v *vStat) Mode() os.FileMode {
	return v.mode
}

func (v *vStat) ModTime() time.Time {
	return v.modTime
}

func (v *vStat) Sys() interface{} {
	return v.sys
}
