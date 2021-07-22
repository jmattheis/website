package redis

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/jmattheis/website/content"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port string
}

func Listen(conf Config) {
	listener, err := net.Listen("tcp", ":"+conf.Port)

	if err != nil {
		log.Fatal().Str("on", "init").Str("port", conf.Port).Err(err).Msg("tcp")
	}

	tty := &content.SingleText{
		Protocol:      "redis",
		Port:          conf.Port,
		Split:         ".",
		CommandPrefix: "redis-cli -h jmattheis.de lrange ",
		CommandSuffix: " 0 0",
	}

	log.Info().Str("on", "init").Str("port", conf.Port).Msg("redis")
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
