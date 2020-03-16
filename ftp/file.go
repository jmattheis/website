package ftp

import (
	"io"
	"os"
	"time"
)

type virtualFile struct {
	content    []byte
	readOffset int
}

func (f *virtualFile) Close() error {
	return nil
}

func (f *virtualFile) Read(buffer []byte) (int, error) {
	n := copy(buffer, f.content[f.readOffset:])
	f.readOffset += n

	if n == 0 {
		return 0, io.EOF
	}

	return n, nil
}

func (f *virtualFile) Seek(int64, int) (int64, error) {
	return 0, nil
}

func (f *virtualFile) Write([]byte) (int, error) {
	return 0, nil
}

type virtualFileInfo struct {
	dir     string
	name    string
	size    int64
	mode    os.FileMode
	content []byte
}

func (f virtualFileInfo) Name() string {
	return f.name
}

func (f virtualFileInfo) Size() int64 {
	return f.size
}

func (f virtualFileInfo) Mode() os.FileMode {
	return f.mode
}

func (f virtualFileInfo) IsDir() bool {
	return f.mode.IsDir()
}

func (f virtualFileInfo) ModTime() time.Time {
	return time.Now().UTC()
}

func (f virtualFileInfo) Sys() interface{} {
	return nil
}
