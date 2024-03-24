package ssh

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
	"github.com/rs/zerolog/log"
	xssh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func Listen() {
	port := util.PortOf(22)
	privateKey, err := readKey("./privkey")
	if err != nil {
		log.Fatal().
			Str("on", "init").
			Str("port", port.S).
			Err(fmt.Errorf("reading private key %s", err)).
			Msg("ssh")
	}

	server := ssh.Server{
		IdleTimeout: time.Minute,
		MaxTimeout:  time.Minute * 10,
		Addr:        port.Addr,
		HostSigners: []ssh.Signer{privateKey},
		Handler: ssh.Handler(func(s ssh.Session) {
			tty := &content.InteractiveText{
				RemoteAddr: s.RemoteAddr().String(),
			}
			defer s.Close()
			term := term.NewTerminal(s, "\nguest@jmattheis.de > ")
			term.AutoCompleteCallback = autocomplete(s)
			exec, _ := tty.Exec(s.RawCommand())

			if _, err := io.WriteString(term, exec); err != nil {
				return
			}

			if s.RawCommand() != "" {
				return
			}

			for {
				line, err := term.ReadLine()
				if err != nil {
					return
				}
				response, exit := tty.Exec(line)
				if _, err := io.WriteString(term, response); err != nil {
					return
				}
				if exit {
					return
				}
			}
		}),
	}
	log.Info().
		Str("on", "init").
		Str("port", port.S).
		Msg("ssh")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("ssh")
		}
	}()
}

func readKey(path string) (ssh.Signer, error) {
	privateBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return xssh.ParsePrivateKey(privateBytes)
}

var auto = []string{"help", "exit", "blog", "projects", "ip", "time"}

func autocomplete(s ssh.Session) func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	return func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		if key == 3 {
			_ = s.Close()
			return
		}

		if key == '\t' {
			for _, cmd := range auto {
				if strings.HasPrefix(cmd, strings.ToLower(line)) {
					newLine = cmd
					ok = true
					newPos = len(cmd)
					return
				}
			}
		}

		return
	}
}
