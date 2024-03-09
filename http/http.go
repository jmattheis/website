package http

import (
	"crypto/tls"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/feeds"
	"github.com/jmattheis/website/assets"
	"github.com/jmattheis/website/content"
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
				Err(err).
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
				Err(err).
				Str("port", conf.SSLPort).
				Msg("https/*")
			return
		}

	}()
}

func handle(port string) http.HandlerFunc {
	handleWS := websocket.Handle(port)
	handleText := text.Handle(port)
	handleHTML := html.Handler()

	feed := feeds.Atom{Feed: content.BlogsRss()}
	atom, err := feed.ToAtom()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ssh" || r.URL.Path == "/key" || r.URL.Path == "/keys" {
			w.Header().Add("content-type", "text/plain")
			io.WriteString(w, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAxpgcVSnqwvdtBz8Vw0PAdP2sMelg5DsYpFbQdXqmxT ssh@jmattheis.de")
			return
		}

		if r.URL.Path == "/feed.xml" || r.URL.Path == "/blog/index.xml" { //
			w.Header().Add("content-type", "application/xml; charset=utf-8")
			io.WriteString(w, atom)
			return
		}

		if containsHeader(r, "connection", "upgrade") &&
			containsHeader(r, "upgrade", "websocket") {
			handleWS(w, r)
			return
		}

		ext := path.Ext(r.URL.Path)
		if ext != "" && ext != ".html" {
			m := mime.TypeByExtension(ext)

			content, err := assets.Assets.ReadFile(r.URL.Path)
			if err != nil {
				http.Error(w, "not found", 404)
				return
			}
			w.Header().Add("content-type", m)
			w.WriteHeader(200)

			_, _ = w.Write(content)
			return
		}

		if containsHeader(r, "user-agent", "httpie") ||
			containsHeader(r, "user-agent", "curl") ||
			containsHeader(r, "accept", "text/plain") {
			handleText(w, r)
			return
		}

		handleHTML(w, r)
	}
}

func containsHeader(r *http.Request, name, part string) bool {
	return strings.Contains(strings.ToLower(r.Header.Get(name)), part)
}
