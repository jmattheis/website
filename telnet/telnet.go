package telnet

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
		log.Fatal().Str("on", "init").Str("port", conf.Port).Err(err).Msg("tcp")
	}

	tty := &content.InteractiveText{
		Prompt:   "\nguest@jmattheis.de > ",
		Protocol: "telnet",
		Port:     conf.Port,
	}

	log.Info().Str("on", "init").Str("port", conf.Port).Msg("tcp")
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

func accept(conn net.Conn, tty *content.InteractiveText) {
	defer conn.Close()
	data, _ := tty.Exec("start")
	if _, err := conn.Write([]byte(data)); err != nil {
		return
	}

	reader := bufio.NewReader(conn)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		exec, exit := tty.Exec(string(line))
		if _, err := conn.Write([]byte(exec)); err != nil {
			return
		}
		if exit {
			return
		}
	}
}
