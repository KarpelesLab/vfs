package vfs_test

import (
	"testing"

	"github.com/KarpelesLab/vfs"
	"github.com/KarpelesLab/vfs/memfs"
)

func TestMkdirAll(t *testing.T) {
	fs := memfs.New()

	err := vfs.MkdirAll(fs, "/a/b/C", 0755)
	if err != nil {
		t.Errorf("failed test: %s", err)
	}
}
