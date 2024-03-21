package pop

import (
	"github.com/jmattheis/website/pop/popgun"
	"github.com/jmattheis/website/pop/popgun/backends"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
)

func Listen() {
	port := util.PortOf(110)
	server := popgun.NewServer(popgun.Config{ListenInterface: port.Addr}, backends.NoAuth{}, backends.ContentProvider{})
	log.Info().
		Str("on", "init").
		Str("port", port.S).
		Msg("pop3")
	go func() {
		if err := server.Start(); err != nil {
			log.Error().Err(err).Msg("pop3")
		}
	}()
}
