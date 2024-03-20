package finger

import (
	"bufio"
	"net"
	"strconv"

	"github.com/jmattheis/website/content"
	"github.com/rs/zerolog/log"
)

func Listen(prod bool) {
	port := 79
	if !prod {
		port = 10079
	}
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatal().Str("on", "init").Int("port", port).Err(err).Msg("finger")
	}

	tty := &content.SingleText{
		Protocol:      "finger",
		Port:          strconv.Itoa(port),
		Split:         ".",
		CommandPrefix: "finger ",
		CommandSuffix: "@jmattheis.de",
	}

	log.Info().Str("on", "init").Int("port", port).Err(err).Msg("finger")
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
