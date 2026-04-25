package ftp

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

var errReadOnly = errors.New("read-only filesystem")

// ClientDriver implements afero.Fs as a read-only virtual filesystem
type ClientDriver struct {
	Files []*virtualFileInfo
	Dirs  []string
}

func (d *ClientDriver) Name() string { return "virtual" }

func (d *ClientDriver) Create(name string) (afero.File, error)                      { return nil, errReadOnly }
func (d *ClientDriver) Mkdir(name string, perm os.FileMode) error                   { return errReadOnly }
func (d *ClientDriver) MkdirAll(path string, perm os.FileMode) error                { return errReadOnly }
func (d *ClientDriver) Remove(name string) error                                    { return errReadOnly }
func (d *ClientDriver) RemoveAll(path string) error                                 { return errReadOnly }
func (d *ClientDriver) Rename(oldname, newname string) error                        { return errReadOnly }
func (d *ClientDriver) Chmod(name string, mode os.FileMode) error                   { return errReadOnly }
func (d *ClientDriver) Chown(name string, uid, gid int) error                       { return errReadOnly }
func (d *ClientDriver) Chtimes(name string, atime time.Time, mtime time.Time) error { return errReadOnly }

func (d *ClientDriver) Stat(name string) (os.FileInfo, error) {
	// Check if it's a directory
	for _, dir := range d.Dirs {
		if name == dir {
			return &virtualFileInfo{
				dir:  filepath.Dir(name),
				name: filepath.Base(name),
				size: 4096,
				mode: 0666 | os.ModeDir,
			}, nil
		}
	}
	// Check files
	for _, file := range d.Files {
		if name == filepath.Join(file.dir, file.name) {
			return file, nil
		}
	}
	return nil, os.ErrNotExist
}

func (d *ClientDriver) Open(name string) (afero.File, error) {
	return d.OpenFile(name, os.O_RDONLY, 0)
}

func (d *ClientDriver) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if flag&(os.O_WRONLY|os.O_RDWR|os.O_CREATE|os.O_TRUNC) != 0 {
		return nil, errReadOnly
	}

	// Check if it's a directory
	for _, dir := range d.Dirs {
		if name == dir {
			return &virtualDir{
				info:  &virtualFileInfo{dir: filepath.Dir(name), name: filepath.Base(name), size: 4096, mode: 0666 | os.ModeDir},
				files: d.filesInDir(name),
			}, nil
		}
	}
	// Check files
	for _, file := range d.Files {
		if name == filepath.Join(file.dir, file.name) {
			return &virtualFile{info: file, content: file.content}, nil
		}
	}
	return nil, os.ErrNotExist
}

func (d *ClientDriver) filesInDir(dir string) []os.FileInfo {
	var result []os.FileInfo
	for _, file := range d.Files {
		if file.dir == dir {
			result = append(result, file)
		}
	}
	return result
}
