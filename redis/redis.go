package redis

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
)

func Listen() {
	port := util.PortOf(6379)
	listener, err := net.Listen("tcp", port.Addr)

	if err != nil {
		log.Fatal().Str("on", "init").Str("port", port.S).Err(err).Msg("tcp")
	}

	log.Info().Str("on", "init").Str("port", port.S).Msg("redis")
	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				continue
			}
			tty := &content.SingleText{
				Split:         ".",
				CommandPrefix: "redis-cli -h jmattheis.de lrange ",
				CommandSuffix: " 0 0",
				RemoteAddr:    conn.RemoteAddr().String(),
			}

			go accept(conn, tty)
		}
	}()
}

func accept(conn net.Conn, tty *content.SingleText) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		if string(line) != "lrange" {
			continue
		}
		reader.ReadLine()
		line, _, _ = reader.ReadLine()
		resp := tty.Get(string(line))
		respLines := strings.Split(resp, "\n")
		conn.Write([]byte(fmt.Sprintf("*%d\r\n", len(respLines))))
		for _, p := range respLines {
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(p), p)))
		}
	}
}
