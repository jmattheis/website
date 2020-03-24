package docker

import (
	"encoding/json"
	"github.com/jmattheis/website/content"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type Config struct {
	Port    string
}

func Listen(conf Config) {
	log.Info().
		Str("on", "init").
		Str("port", conf.Port).
		Msg("docker")

	tty := &content.SingleText{
		Protocol:      "http@docker",
		Port:          conf.Port,
		Split:         "/",
		CommandPrefix: "docker -H jmattheis.de inspect -f '{{.Value}}' ",
	}

	go func() {

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if !strings.Contains(path, "containers") {
				w.WriteHeader(404);
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

		err := http.ListenAndServe(":"+conf.Port, mux)
		if err != nil {
			log.Fatal().
				Str("on", "init").
				Str("port", conf.Port).
				Msg("docker")
		}
	}()
}

