package vfs

import "net/http"

type httpFsHandler struct {
	FileSystem
}

func MakeHttpFileSystem(fs FileSystem) http.FileSystem {
	return &httpFsHandler{fs}
}

func (h *httpFsHandler) Open(name string) (http.File, error) {
	return h.FileSystem.Open(name)
}
