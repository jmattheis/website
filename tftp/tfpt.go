package tftp

import (
	"fmt"
	"io"
	"strings"

	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
	"github.com/pin/tftp/v3"
)

func Listen() {
	port := util.PortOf(69)

	tty := &content.SingleText{
		Split:         "/",
		CommandPrefix: "curl tftp://jmattheis.de/",
		RemoteAddr:    "unknown",
	}

	svr := tftp.NewServer(func(filename string, rf io.ReaderFrom) error {
		content := tty.Get(filename)
		_, err := rf.ReadFrom(strings.NewReader(content))
		return err
	}, func(filename string, wt io.WriterTo) error {
		return fmt.Errorf("no thanks :P")
	})
	go svr.ListenAndServe(port.Addr)
}
