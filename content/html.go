package content

import (
	"github.com/gobuffalo/packr/v2"
	"strings"
	"text/template"
)

var Assets = packr.New("html", "../assets")

var HtmlTemplates = MustLoadBoxedTemplate(Assets);

// https://github.com/gobuffalo/packr/issues/16#issuecomment-354905578
func MustLoadBoxedTemplate(b *packr.Box) *template.Template {
	t := template.New("")
	err := b.Walk(func(p string, f packr.File) error {
		if p == "" {
			return nil
		}
		var err error
		var csz int64
		if finfo, err := f.FileInfo(); err != nil {
			return err
		}else{
			// skip directory path
			if finfo.IsDir() {
				return nil
			}
			csz= finfo.Size()
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

		var h = make([]byte, 0, csz)

		if h, err = b.Find(p); err != nil {
			return err
		}

		if _, err = t.New(n).Parse(string(h)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic("error loading template")
	}
	return t
}
