package dockerregistry

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/jmattheis/website/content"
)

var pattern = regexp.MustCompile(`\/v2\/(.+)\/(blobs|manifests)\/([a-z0-9:]+)`)

func Handler() http.HandlerFunc {
	tty := content.SingleText{
		Split:         "/",
		CommandPrefix: "docker run --rm jmattheis.de/",
	}

	successCache, err := lru.New[string, map[string]Entry](50)
	check(err)
	errorCache, err := lru.New[string, map[string]Entry](20)
	check(err)

	return func(w http.ResponseWriter, r *http.Request) {
		matches := pattern.FindStringSubmatch(r.URL.Path)

		if matches == nil {
			io.WriteString(w, "unknown url")
			w.WriteHeader(404)
			return
		}

		cmd, t, hash := matches[1], matches[2], matches[3]

		content := tty.Get(cmd)

		cache := successCache
		if strings.HasPrefix(content, "error") {
			cache = errorCache
		}

		store, ok := cache.Get(cmd)
		if !ok {
			store = DockerStdout(content)
			cache.Add(cmd, store)
		}

		entry, ok := store[hash]
		if !ok {
			if strings.HasPrefix(hash, "sha256:") {
				w.WriteHeader(404)
				io.WriteString(w, "unknown hash")
				return
			}
			entry, _ = store["latest"]
		}

		w.Header().Add("content-type", entry.MediaType)
		w.Header().Add("content-length", fmt.Sprint(len(entry.Content)))
		w.Header().Add("etag", `"`+entry.Digest.String()+`"`)
		if t == "manifests" {
			w.Header().Add("Docker-Content-Digest", entry.Digest.String())
		}
		w.WriteHeader(200)
		if r.Method != "HEAD" {
			w.Write(entry.Content)
		}
	}
}
