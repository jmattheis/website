package content

import (
	"text/template"

	"github.com/jmattheis/website/assets"
)

var HtmlTemplates *template.Template

func init() {
	var err error
	HtmlTemplates, err = template.ParseFS(assets.Assets, "*.html")
	if err != nil {
		panic(err)
	}
}
