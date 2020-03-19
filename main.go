package main

import (
	"github.com/jmattheis/website/dns"
	"github.com/jmattheis/website/ftp"
	"github.com/jmattheis/website/http"
	"github.com/jmattheis/website/logger"
	"github.com/jmattheis/website/ssh"
	"github.com/jmattheis/website/telnet"
	"github.com/jmattheis/website/whois"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme/autocert"
)

var (
	Mode = "dev"
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
}

var prodConf = Config{
	AutoCert: true,
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
		Port:    "80",
		SSLPort: "443",
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
		Port:    "10080",
		SSLPort: "10443",
	},
}

func main() {
	logger.Init(zerolog.DebugLevel)

	var config Config

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
	ssh.Listen(config.SSH, )
	http.Listen(config.HTTP, certManager)

	<-make(chan struct{})
}
