package whois

import (
	"bufio"
	"net"

	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
)

func Listen() {
	port := util.PortOf(43)
	listener, err := net.Listen("tcp", port.Addr)

	if err != nil {
		log.Fatal().Str("on", "init").Str("port", port.S).Err(err).Msg("whois")
	}

	log.Info().Str("on", "init").Str("port", port.S).Msg("whois")
	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				continue
			}

			tty := &content.SingleText{
				Split:         ".",
				CommandPrefix: "whois -h jmattheis.de ",
				RemoteAddr:    conn.RemoteAddr().String(),
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
