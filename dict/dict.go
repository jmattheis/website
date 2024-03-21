package dict

import (
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
)

func Listen() {
    port := util.PortOf(2628)
	server := NewServer()

	log.Info().Str("on", "init").Str("port", port.S).Msg("dict")
	go func() {
		if err := server.Start(port.S); err != nil {
			log.Fatal().Str("on", "init").Str("port", port.S).Err(err).Msg("dict")
		}
	}()
}
