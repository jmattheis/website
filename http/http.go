package http

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/jmattheis/website/http/html"
	"github.com/jmattheis/website/http/text"
	"github.com/jmattheis/website/http/websocket"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"
)

type Config struct {
	Port    string
	SSLPort string
}

func Listen(conf Config, manager *autocert.Manager) {
	log.Info().
		Str("on", "init").
		Str("port", conf.Port).
		Msg("http/*")

	go func() {

		var handler http.Handler = gziphandler.GzipHandler(handle(conf.Port))
		if manager != nil {
			handler = manager.HTTPHandler(handler)
		}

		err := http.ListenAndServe(":"+conf.Port, handler)
		if err != nil {
			log.Fatal().
				Str("on", "init").
				Str("port", conf.Port).
				Msg("http/*")
		}
	}()

	log.Info().
		Str("on", "init").
		Bool("autocert", manager != nil).
		Str("port", conf.SSLPort).
		Msg("https/*")
	go func() {
		server := &http.Server{
			Addr:      ":" + conf.SSLPort,
			Handler:   gziphandler.GzipHandler(handle(conf.SSLPort)),
			TLSConfig: &tls.Config{},
		}
		if manager == nil {
			server.TLSConfig.Certificates = []tls.Certificate{*util.NewUntrustedCert()}
		} else {
			server.TLSConfig.GetCertificate = manager.GetCertificate
		}

		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatal().
				Str("on", "init").
				Str("port", conf.SSLPort).
				Msg("https/*")
			return
		}

	}()
}

func handle(port string) http.HandlerFunc {
	ws := websocket.Handle(port)
	t := text.Handle(port)

	return func(w http.ResponseWriter, r *http.Request) {
		if containsHeader(r, "connection", "upgrade") &&
			containsHeader(r, "upgrade", "websocket") {
			ws(w, r)
			return
		}

		if containsHeader(r, "user-agent", "httpie") ||
			containsHeader(r, "user-agent", "curl") ||
			containsHeader(r, "accept", "text/plain") {
			t(w, r)
			return
		}

		html.Handle(w, r)
	}
}

func containsHeader(r *http.Request, name, part string) bool {
	return strings.Contains(strings.ToLower(r.Header.Get(name)), part)
}
