package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

//go:embed *
var Assets embed.FS
var Blogs, _ = fs.Sub(Assets, "blog")

var BlogList []string
var BlogContent []string

func init() {
	err := fs.WalkDir(Assets, "blog", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := Assets.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %w", err)
		}
		BlogList = append(BlogList, filepath.Base(path))
		BlogContent = append(BlogContent, string(content))
		return nil
	})
	if err != nil {
		panic(err)
	}
}
