package text

import (
	"net/http"
	"strings"

	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
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
			RemoteAddr:    util.GetRemoteAddr(r),
		}

		value := tty.Get(strings.TrimPrefix(r.URL.Path, "/"))

		_, _ = w.Write([]byte(value))
	}

}
