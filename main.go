package main

import (
	"os"

	"github.com/jmattheis/website/dict"
	"github.com/jmattheis/website/dns"
	"github.com/jmattheis/website/docker"
	"github.com/jmattheis/website/finger"
	"github.com/jmattheis/website/ftp"
	"github.com/jmattheis/website/gemini"
	"github.com/jmattheis/website/gopher"
	"github.com/jmattheis/website/http"
	"github.com/jmattheis/website/imap"
	"github.com/jmattheis/website/logger"
	"github.com/jmattheis/website/pop"
	"github.com/jmattheis/website/redis"
	"github.com/jmattheis/website/ssh"
	"github.com/jmattheis/website/telnet"
	"github.com/jmattheis/website/tftp"
	"github.com/jmattheis/website/util"
	"github.com/jmattheis/website/whois"
	"github.com/rs/zerolog"
)

func main() {
	logger.Init(zerolog.DebugLevel)

	ip := "127.0.0.1"

	if os.Args[1] == "prod" {
		ip = "78.47.104.216"
	} else {
		util.PortOffset = 10000
	}

	telnet.Listen()
	whois.Listen()
	dns.Listen()
	ftp.Listen(ip)
	ssh.Listen()
	http.Listen()
	pop.Listen()
	imap.Listen()
	dict.Listen()
	gopher.Listen()
	docker.Listen()
	gemini.Listen()
	redis.Listen()
	finger.Listen()
	tftp.Listen()

	<-make(chan struct{})
}
