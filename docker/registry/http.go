package dockerregistry

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
)

var pattern = regexp.MustCompile(`\/v2\/(.+)\/(blobs|manifests)\/([a-z0-9:]+)`)

func Handler() http.HandlerFunc {

	static := map[string]Entry{}
	RegisterLayers(func(s string, e Entry) { static[s] = e })

	dynamicCache, _ := lru.New[string, Entry](200)

	return func(w http.ResponseWriter, r *http.Request) {
		tty := content.SingleText{
			Split:         "/",
			CommandPrefix: "docker run --rm jmattheis.de/",
			RemoteAddr:    util.GetRemoteAddr(r),
		}

		matches := pattern.FindStringSubmatch(r.URL.Path)

		if matches == nil {
			io.WriteString(w, "unknown url")
			w.WriteHeader(404)
			return
		}

		cmd, t, hash := matches[1], matches[2], matches[3]

		content, dynamic := tty.GetVerbose(cmd)

		var result Entry
		var ok bool

		if !strings.HasPrefix(hash, "sha256:") {
			ok = true
			result = DockerStdout(func(s string, e Entry) {
				if dynamic {
					dynamicCache.Add(s, e)
				} else {
					static[s] = e
				}
			}, content)
		} else {
			result, ok = static[hash]
			if !ok {
				result, ok = dynamicCache.Get(hash)
			}
		}

		if !ok {
			w.WriteHeader(404)
			io.WriteString(w, "unknown reference")
			return
		}

		w.Header().Add("content-type", result.MediaType)
		w.Header().Add("content-length", fmt.Sprint(len(result.Content)))
		w.Header().Add("etag", `"`+result.Digest.String()+`"`)
		if t == "manifests" {
			w.Header().Add("Docker-Content-Digest", result.Digest.String())
		}
		w.WriteHeader(200)
		if r.Method != "HEAD" {
			w.Write(result.Content)
		}
	}
}
