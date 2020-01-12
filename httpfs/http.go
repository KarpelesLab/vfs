package httpfs

import (
	"net/http"

	"github.com/KarpelesLab/vfs"
)

type httpFsHandler struct {
	vfs.FileSystem
}

func MakeHttpFileSystem(fs vfs.FileSystem) http.FileSystem {
	return &httpFsHandler{fs}
}

func (h *httpFsHandler) Open(name string) (http.File, error) {
	return h.FileSystem.Open(name)
}
