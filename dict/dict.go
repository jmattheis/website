package dict

import (
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port string
}

func Listen(conf Config) {
	server := NewServer()

	log.Info().Str("on", "init").Str("port", conf.Port).Msg("dict")
	go func() {
		if err := server.Start(conf.Port); err != nil {
			log.Fatal().Str("on", "init").Str("port", conf.Port).Err(err).Msg("dict")
		}
	}()
}
