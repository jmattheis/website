package content

import (
	"io"
	"io/fs"
	"strings"
	"text/template"

	"github.com/jmattheis/website/assets"
)

var HtmlTemplates = MustLoadBoxedTemplate()

// https://github.com/gobuffalo/packr/issues/16#issuecomment-354905578
func MustLoadBoxedTemplate() *template.Template {
	t := template.New("")
	err := fs.WalkDir(assets.Assets, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == "" || d.IsDir() {
			return nil
		}

		// skip all files except .html
		if !strings.Contains(p, ".html") {
			return nil
		}

		// Normalize template name
		n := p
		if strings.HasPrefix(p, "\\") || strings.HasPrefix(p, "/") {
			n = n[1:] // don't want template name to start with / ie. /index.html
		}
		// replace windows path seperator \ to normalized /
		n = strings.Replace(n, "\\", "/", -1)

		x, err := assets.Assets.Open(n)
		if err != nil {
			return err
		}
		defer x.Close()
		b, err := io.ReadAll(x)
		if err != nil {
			return err
		}

		if _, err = t.New(n).Parse(string(b)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic("error loading template" + err.Error())
	}
	return t
}
