package util

import (
	"net/http"
	"strings"
)

func GetRemoteAddr(r *http.Request) string {
	forwarded := r.Header.Get("x-real-ip")
	if forwarded == "" {
		forwarded = r.RemoteAddr
	}
	forwarded = strings.SplitN(forwarded, ",", 2)[0]
	return forwarded
}
