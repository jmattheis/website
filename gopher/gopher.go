package gopher

import (
	"strings"

	"git.mills.io/prologic/go-gopher"
	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
)

func Listen() {
	port := util.PortOf(70)

	mux := gopher.NewServeMux()
	mux.HandleFunc("/", func(w gopher.ResponseWriter, r *gopher.Request) {
		tty := &content.SingleText{
			Split:         "/",
			CommandPrefix: "curl gopher://jmattheis.de/0",
			RemoteAddr:    "unknown",
		}
		value := tty.Get(strings.TrimPrefix(r.Selector, "/"))
		_, _ = w.Write([]byte(value))
	})
	log.Info().
		Str("on", "init").
		Str("port", port.S).
		Msg("gopher")
	go func() {
		if err := gopher.ListenAndServe(port.Addr, mux); err != nil {
			log.Fatal().
				Str("on", "init").
				Str("port", port.S).
				Err(err).
				Msg("gopher")
		}
	}()
}
