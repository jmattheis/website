package imap

import (
	"github.com/emersion/go-imap/server"
	"github.com/jmattheis/website/imap/backend"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port string
}

func Listen(conf Config) {
	be := backend.New(conf.Port)

	s := server.New(be)
	s.Addr = ":" + conf.Port
	s.AllowInsecureAuth = true

	log.Info().Str("on", "init").Str("port", conf.Port).Msg("imap")
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("imap")
		}
	}()
}
