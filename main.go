package main

import (
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
	"github.com/jmattheis/website/whois"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme/autocert"
)

var (
	Mode = "prod"
)

type Config struct {
	AutoCert bool
	Domain   string
	PubIP    string
	Telnet   telnet.Config
	Whois    whois.Config
	DNS      dns.Config
	FTP      ftp.Config
	SSH      ssh.Config
	HTTP     http.Config
	POP      pop.Config
	IMAP     imap.Config
	DICT     dict.Config
	Gopher   gopher.Config
	Docker   docker.Config
	Gemini   gemini.Config
	Redis    redis.Config
}

var prodConf = Config{
	AutoCert: false,
	Domain:   "jmattheis.de",
	PubIP:    "78.47.104.216",
	Telnet: telnet.Config{
		Port: "23",
	},
	Whois: whois.Config{
		Port: "43",
	},
	DNS: dns.Config{
		Port: "53",
	},
	FTP: ftp.Config{
		Port: "21",
	},
	SSH: ssh.Config{
		Port:           "22",
		PrivateKeyPath: "./privkey",
	},
	HTTP: http.Config{
		Port:    "8080",
		SSLPort: "8083",
	},
	POP: pop.Config{
		Port: "110",
	},
	IMAP: imap.Config{
		Port: "143",
	},
	DICT: dict.Config{
		Port: "2628",
	},
	Gopher: gopher.Config{
		Port: "70",
	},
	Gemini: gemini.Config{
		Port: "1965",
	},
	Docker: docker.Config{
		Port: "2375",
	},
	Redis: redis.Config{
		Port: "6379",
	},
}

var devConf = Config{
	AutoCert: false,
	Domain:   "jmattheis.de",
	PubIP:    "127.0.0.1",
	Telnet: telnet.Config{
		Port: "10023",
	},
	Whois: whois.Config{
		Port: "10043",
	},
	DNS: dns.Config{
		Port: "10053",
	},
	FTP: ftp.Config{
		Port: "10021",
	},
	SSH: ssh.Config{
		Port:           "10022",
		PrivateKeyPath: "./privkey",
	},
	HTTP: http.Config{
		Port:    "8080",
		SSLPort: "10443",
	},
	POP: pop.Config{
		Port: "10110",
	},
	IMAP: imap.Config{
		Port: "10143",
	},
	DICT: dict.Config{
		Port: "2628",
	},
	Gopher: gopher.Config{
		Port: "10070",
	},
	Gemini: gemini.Config{
		Port: "1965",
	},
	Docker: docker.Config{
		Port: "2375",
	},
	Redis: redis.Config{
		Port: "6379",
	},
}

func main() {
	logger.Init(zerolog.DebugLevel)

	var config Config
	prod := Mode == "prod"

	if Mode == "prod" {
		config = prodConf
	} else {
		config = devConf
	}

	var certManager *autocert.Manager
	if config.AutoCert {
		certManager = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.Domain), //Your domain here
			Cache:      autocert.DirCache("certs"),            //Folder for storing certificates
		}
	}
	telnet.Listen(config.Telnet)
	whois.Listen(config.Whois)
	dns.Listen(config.DNS)
	ftp.Listen(config.FTP, certManager, config.PubIP)
	ssh.Listen(config.SSH)
	http.Listen(config.HTTP, certManager)
	pop.Listen(config.POP)
	imap.Listen(config.IMAP)
	dict.Listen(config.DICT)
	gopher.Listen(config.Gopher)
	docker.Listen(config.Docker)
	gemini.Listen(config.Gemini)
	redis.Listen(config.Redis)
	finger.Listen(prod)

	<-make(chan struct{})
}
