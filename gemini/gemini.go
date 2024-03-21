package gemini

import (
	"context"
	"os"
	"path"
	"strings"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"
	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
)

func Listen() {
    port := util.PortOf(1965)
	log.Info().
		Str("on", "init").
		Str("port", port.S).
		Msg("gemini")

	tty := &content.SingleText{
		Split:          "/",
		CommandPrefix:  "=> gemini://jmattheis.de/",
		DisablePadding: true,
	}

	cwd, _ := os.Getwd()
	certificates := &certificate.Store{}
	p := path.Join(cwd, "gemini_cache")
	os.MkdirAll(p, 0755)
	certificates.SetPath(p)
	certificates.Load(p)
	certificates.Register("jmattheis.de")
	certificates.Register("localhost")

	mux := &gemini.Mux{}
	mux.Handle("/", gemini.HandlerFunc(func(c context.Context, rw gemini.ResponseWriter, r *gemini.Request) {
		rw.WriteHeader(gemini.StatusSuccess, "text/gemini")
		value := tty.Get(strings.TrimPrefix(r.URL.EscapedPath(), "/"))
		rw.Write([]byte(value))
	}))

	server := &gemini.Server{
		Addr:           ":" + port.S,
		Handler:        mux,
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		GetCertificate: certificates.Get,
	}
	go server.ListenAndServe(context.Background())

}
