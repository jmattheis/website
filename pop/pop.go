package pop

import (
	"github.com/jmattheis/website/pop/popgun"
	"github.com/jmattheis/website/pop/popgun/backends"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port string
}

func Listen(conf Config) {
	server := popgun.NewServer(popgun.Config{ListenInterface: ":" + conf.Port}, backends.NoAuth{}, backends.ContentProvider{Port: conf.Port})
	log.Info().
		Str("on", "init").
		Str("port", conf.Port).
		Msg("pop3")
	go func() {
		if err := server.Start(); err != nil {
			log.Error().Err(err).Msg("pop3")
		}
	}()
}
