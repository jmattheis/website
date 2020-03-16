package content

import (
	"github.com/gobuffalo/packr/v2"
	"net/http"
)

var HtmlServe = http.FileServer(packr.New("html", "../website/docs"))
