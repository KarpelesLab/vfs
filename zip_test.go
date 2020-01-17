package vfs_test

import (
	"archive/zip"
	"os"
	"testing"

	"github.com/KarpelesLab/vfs/zipfs"
)

func TestZipSmall(t *testing.T) {
	f, err := os.Open("testfiles/small.zip")
	if err != nil {
		t.Errorf("failed to open test zip file: %s", err)
		return
	}

	st, err := f.Stat()
	if err != nil {
		t.Errorf("failed to stat zip file: %s", err)
		return
	}

	zip, err := zip.NewReader(f, st.Size())
	if err != nil {
		t.Errorf("failed to init zip file: %s", err)
		return
	}

	zfs, err := zipfs.New(zip, zipfs.IndexFull)
	if err != nil {
		t.Errorf("failed to init zipfs: %s", err)
		return
	}

	// testing
	st, err = zfs.Stat("a")
	if err != nil {
		t.Errorf("failed to stat \"a\" in zip file: %s", err)
	}
}
