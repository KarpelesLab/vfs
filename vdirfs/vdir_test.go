package vdirfs_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/KarpelesLab/vfs"
	"github.com/KarpelesLab/vfs/memfs"
	"github.com/KarpelesLab/vfs/vdirfs"
)

func putFile(fs vfs.FileSystem, name, data string) error {
	err := vfs.MkdirAll(fs, path.Dir(name), 0755)
	if err != nil {
		return fmt.Errorf("mkdirall failed: %w", err)
	}

	fp, err := fs.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0755)
	return err
	_, err = fp.Write([]byte(data))
	fp.Close()
	return err
}

func TestVdir(t *testing.T) {
	// populate stuff
	m := memfs.New()

	files := []string{
		"a/b/c/hello",
		"a/b/d/hello",
		"a/b/e/hello",
	}

	for _, f := range files {
		err := putFile(m, f, "hello world")
		if err != nil {
			t.Errorf("failed to create file %s: %s", f, err)
			// give up as next tests are bound to fail
			return
		}
	}

	// generate index
	n, err := vdirfs.New(m)
	if err != nil {
		t.Errorf("failed to instanciate vdirfs: %s", err)
		return
	}
	for _, f := range files {
		n.AddPath(f)
	}

	// list files in a/b
	list, err := vfs.ReadDir(n, "a/b")
	if err != nil {
		t.Errorf("list dir failed: %s", err)
		return
	}

	if len(list) != 3 {
		t.Errorf("expected 3 values, got %d", len(list))
		return
	}

	successList := []string{"c", "d", "e"}
	for i, f := range list {
		if f.Name() != successList[i] {
			t.Errorf("expected entry %d to be named %s, received %s", i, successList[i], f.Name())
		}
	}

	// let's stat
	if _, err := n.Stat("a"); err != nil {
		t.Errorf("failed to stat a: %s", err)
	}
}
