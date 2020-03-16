package html

import (
	"github.com/jmattheis/website/content"
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".js") {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	}
	if strings.HasSuffix(r.URL.Path, ".html") {
		w.Header().Set("Content-Type", "text/html")
	}
	if strings.HasSuffix(r.URL.Path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	}
	content.HtmlServe.ServeHTTP(w, r)
}
