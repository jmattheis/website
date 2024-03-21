package docker

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
)

func Listen() {
	port := util.PortOf(2375)
	log.Info().
		Str("on", "init").
		Str("port", port.S).
		Msg("docker")

	tty := &content.SingleText{
		Split:         ".",
		CommandPrefix: "docker -H jmattheis.de inspect -f '{{.Value}}' ",
	}

	go func() {

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if !strings.Contains(path, "containers") {
				w.WriteHeader(404)
				return
			}

			split := strings.Split(path, "/")

			if len(split) < 3 {
				w.WriteHeader(404)
				return
			}

			help := struct {
				Value string `json:"Value"`
			}{
				Value: tty.Get(split[len(split)-2]),
			}
			json.NewEncoder(w).Encode(help)
		})

		err := http.ListenAndServe(port.Addr, mux)
		if err != nil {
			log.Fatal().
				Str("on", "init").
				Str("port", port.S).
				Msg("docker")
		}
	}()
}
