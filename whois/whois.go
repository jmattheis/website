package whois

import (
	"bufio"
	"github.com/jmattheis/website/content"
	"github.com/rs/zerolog/log"
	"net"
)

type Config struct {
	Port    string
}

func Listen(conf Config) {
	listener, err := net.Listen("tcp", ":"+conf.Port)

	if err != nil {
		log.Fatal().Str("on", "init").Str("port", conf.Port).Err(err).Msg("whois")
	}

	tty := &content.SingleText{
		Protocol: "whois",
		Port:     conf.Port,
		Split:    ".",
		CommandPrefix: "whois -h jmattheis.de ",
	}

	log.Info().Str("on", "init").Str("port", conf.Port).Msg("whois")
	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				continue
			}

			go accept(conn, tty)
		}
	}()
}

func accept(conn net.Conn, tty *content.SingleText) {
	defer conn.Close()
	if line, _, err := bufio.NewReader(conn).ReadLine(); err == nil {
		exec := tty.Get(string(line))
		_, _ = conn.Write([]byte(exec))
	}
}
