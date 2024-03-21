package text

import (
	"github.com/jmattheis/website/content"
	"net/http"
	"strings"
)

func Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain")

		cmd := "curl"
		if strings.Contains(strings.ToLower(r.Header.Get("user-agent")), "httpie") {
			cmd = "http"
		}
		tty := &content.SingleText{
			Split:         "/",
			CommandPrefix: cmd + " jmattheis.de/",
		}

		value := tty.Get(strings.TrimPrefix(r.URL.Path, "/"))

		_, _ = w.Write([]byte(value))
	}

}
