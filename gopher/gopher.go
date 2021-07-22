package gopher

import (
	"github.com/jmattheis/website/content"
	"git.mills.io/prologic/go-gopher"
	"github.com/rs/zerolog/log"
	"strings"
)

type Config struct {
	Port string
}

func Listen(conf Config) {

	tty := &content.SingleText{
		Protocol:      "gopher",
		Port:          conf.Port,
		Split:         "/",
		CommandPrefix: "curl gopher://jmattheis.de/0",
	}
	mux := gopher.NewServeMux()
	mux.HandleFunc("/", func(w gopher.ResponseWriter, r *gopher.Request) {
		value := tty.Get(strings.TrimPrefix(r.Selector, "/"))
		_, _ = w.Write([]byte(value))
	})
	log.Info().
		Str("on", "init").
		Str("port", conf.Port).
		Msg("gopher")
	go func() {
		if err := gopher.ListenAndServe(":"+conf.Port, mux); err != nil {
			log.Fatal().
				Str("on", "init").
				Str("port", conf.Port).
				Err(err).
				Msg("gopher")
		}
	}()
}
