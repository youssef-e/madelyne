package testerfile

import (
	"io"
	"os"
)

type FileOpener interface {
	Open(name string) (io.ReadCloser, error)
}

func New() FileOpener {
	return &fileOpener{}
}

type fileOpener struct {
}

func (f fileOpener) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}
