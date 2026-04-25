package ftp

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type virtualFile struct {
	info       *virtualFileInfo
	content    []byte
	readOffset int
}

func (f *virtualFile) Close() error                                { return nil }
func (f *virtualFile) Name() string                                { return filepath.Join(f.info.dir, f.info.name) }
func (f *virtualFile) Readdir(count int) ([]os.FileInfo, error)    { return nil, os.ErrInvalid }
func (f *virtualFile) Readdirnames(count int) ([]string, error)    { return nil, os.ErrInvalid }
func (f *virtualFile) Stat() (os.FileInfo, error)                  { return f.info, nil }
func (f *virtualFile) Sync() error                                 { return nil }
func (f *virtualFile) Truncate(size int64) error                   { return errReadOnly }
func (f *virtualFile) Write([]byte) (int, error)                   { return 0, errReadOnly }
func (f *virtualFile) WriteAt(b []byte, off int64) (int, error)    { return 0, errReadOnly }
func (f *virtualFile) WriteString(s string) (int, error)           { return 0, errReadOnly }

func (f *virtualFile) Read(buffer []byte) (int, error) {
	if f.readOffset >= len(f.content) {
		return 0, io.EOF
	}
	n := copy(buffer, f.content[f.readOffset:])
	f.readOffset += n
	return n, nil
}

func (f *virtualFile) ReadAt(b []byte, off int64) (int, error) {
	if int(off) >= len(f.content) {
		return 0, io.EOF
	}
	n := copy(b, f.content[off:])
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}

func (f *virtualFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.readOffset = int(offset)
	case io.SeekCurrent:
		f.readOffset += int(offset)
	case io.SeekEnd:
		f.readOffset = len(f.content) + int(offset)
	}
	return int64(f.readOffset), nil
}

type virtualDir struct {
	info     *virtualFileInfo
	files    []os.FileInfo
	readDone bool
}

func (d *virtualDir) Close() error                                { return nil }
func (d *virtualDir) Name() string                                { return filepath.Join(d.info.dir, d.info.name) }
func (d *virtualDir) Read([]byte) (int, error)                    { return 0, os.ErrInvalid }
func (d *virtualDir) ReadAt(b []byte, off int64) (int, error)     { return 0, os.ErrInvalid }
func (d *virtualDir) Seek(int64, int) (int64, error)              { return 0, nil }
func (d *virtualDir) Write([]byte) (int, error)                   { return 0, errReadOnly }
func (d *virtualDir) WriteAt(b []byte, off int64) (int, error)    { return 0, errReadOnly }
func (d *virtualDir) WriteString(s string) (int, error)           { return 0, errReadOnly }
func (d *virtualDir) Stat() (os.FileInfo, error)                  { return d.info, nil }
func (d *virtualDir) Sync() error                                 { return nil }
func (d *virtualDir) Truncate(size int64) error                   { return errReadOnly }

func (d *virtualDir) Readdir(count int) ([]os.FileInfo, error) {
	if d.readDone {
		return nil, io.EOF
	}
	d.readDone = true
	return d.files, nil
}

func (d *virtualDir) Readdirnames(count int) ([]string, error) {
	infos, err := d.Readdir(count)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(infos))
	for i, info := range infos {
		names[i] = info.Name()
	}
	return names, nil
}

type virtualFileInfo struct {
	dir     string
	name    string
	size    int64
	mode    os.FileMode
	content []byte
}

func (f virtualFileInfo) Name() string      { return f.name }
func (f virtualFileInfo) Size() int64       { return f.size }
func (f virtualFileInfo) Mode() os.FileMode { return f.mode }
func (f virtualFileInfo) IsDir() bool       { return f.mode.IsDir() }
func (f virtualFileInfo) ModTime() time.Time { return time.Now().UTC() }
func (f virtualFileInfo) Sys() interface{}  { return nil }
