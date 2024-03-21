package imap

import (
	"github.com/emersion/go-imap/server"
	"github.com/jmattheis/website/imap/backend"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
)

func Listen() {
	port := util.PortOf(143)
	be := backend.New()

	s := server.New(be)
	s.Addr = port.Addr
	s.AllowInsecureAuth = true

	log.Info().Str("on", "init").Str("port", port.S).Msg("imap")
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("imap")
		}
	}()
}
