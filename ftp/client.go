package ftp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fclairamb/ftpserver/server"
)

// ClientDriver defines a very basic client driver
type ClientDriver struct {
	BaseDir string // Base directory from which to server file
	Files   []*virtualFileInfo
	Dirs    []string
}

// ChangeDirectory changes the current working directory
func (driver *ClientDriver) ChangeDirectory(cc server.ClientContext, directory string) error {
	for _, dir := range driver.Dirs {
		if directory == dir {
			return nil
		}
	}
	return errors.New("directory does not exist " + directory)
}

// MakeDirectory creates a directory
func (driver *ClientDriver) MakeDirectory(server.ClientContext, string) error {
	return errors.New("mhhmm, don't do it.")
}

// ListFiles lists the files of a directory
func (driver *ClientDriver) ListFiles(cc server.ClientContext, directory string) ([]os.FileInfo, error) {
	visible := []os.FileInfo{}

	for _, file := range driver.Files {
		if file.dir == directory {
			visible = append(visible, file)
		}
	}

	return visible, nil
}

// OpenFile opens a file in 3 possible modes: read, write, appending write (use appropriate flags)
func (driver *ClientDriver) OpenFile(cc server.ClientContext, path string, flag int) (server.FileStream, error) {
	for _, file := range driver.Files {
		if path == filepath.Join(file.dir, file.name) {
			return &virtualFile{content: file.content}, nil
		}
	}
	return nil, fmt.Errorf("file does not exist %s", path)
}

// GetFileInfo gets some info around a file or a directory
func (driver *ClientDriver) GetFileInfo(cc server.ClientContext, path string) (os.FileInfo, error) {
	for _, file := range driver.Files {
		if path == filepath.Join(file.dir, file.name) {
			return file, nil
		}
	}
	return nil, fmt.Errorf("file does not exist %s", path)
}

// SetFileMtime changes file mtime
func (driver *ClientDriver) SetFileMtime(cc server.ClientContext, path string, mtime time.Time) error {
	return errors.New("not supported")
}

// CanAllocate gives the approval to allocate some data
func (driver *ClientDriver) CanAllocate(cc server.ClientContext, size int) (bool, error) {
	return true, nil
}

// ChmodFile changes the attributes of the file
func (driver *ClientDriver) ChmodFile(cc server.ClientContext, path string, mode os.FileMode) error {
	return errors.New("not supported")
}

// DeleteFile deletes a file or a directory
func (driver *ClientDriver) DeleteFile(cc server.ClientContext, path string) error {
	return errors.New("not supported")
}

// RenameFile renames a file or a directory
func (driver *ClientDriver) RenameFile(cc server.ClientContext, from, to string) error {
	return errors.New("not supported")
}
