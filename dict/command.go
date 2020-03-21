package dict

import (
	"strings"
)

type Executable interface {
	Run(c *Client, args []string)
}

type QuitCommand struct{}

func (cmd QuitCommand) Run(c *Client, args []string) {
	c.isAlive = false
	c.printer.Write(221, "Closing Connection")
}

type ClientCommand struct{}

func (cmd ClientCommand) Run(c *Client, args []string) {

	if len(args) > 0 {
		s := args[0]
		if strings.Contains(strings.ToLower(s), "curl") {
			c.tty.CommandPrefix = "curl dict://jmattheis.de/d:"
			c.printer.Write(250, "ok curl it is")
			return
		}
	}
	c.printer.Write(250, "ok hello")
}

type StatusCommand struct{}

func (cmd StatusCommand) Run(c *Client, args []string) {
	c.printer.Write(210, "nah")
}

type HelpCommand struct{}

func (cmd HelpCommand) Run(c *Client, args []string) {
	c.printer.Write(113, "help text follows")
	c.printer.Plain("read the rfc: https://tools.ietf.org/html/rfc2229")
	c.printer.End()
}

type NotImplementedCommand struct{}

func (cmd NotImplementedCommand) Run(c *Client, args []string) {
	c.printer.NotImplemented()
}

type DefineCommand struct{}

func (cmd DefineCommand) Run(c *Client, args []string) {
	if len(args) != 2 {
		c.printer.IllegalParams()
		return
	}

	word := strings.ReplaceAll(args[1], `"`, "")

	find := c.tty.Get(word)
	if find == "" {
		c.printer.Write(552, "No match")
	} else {
		c.printer.Write(150, "1 definitions retrieved - definitions follow")
		c.printer.Write(151, "%s %s - text follows", word, "jmattheis.de")
		c.printer.Plain(find)
		c.printer.End()
		c.printer.Write(250, "ok")
	}
}

type MatchCommand struct{}

func (cmd MatchCommand) Run(c *Client, args []string) {
	if len(args) != 3 {
		c.printer.IllegalParams()
		return
	}

	commands := c.tty.Commands()

	c.printer.Write(152, "%d matches found - text follows", len(commands))
	for _, w := range commands {
		c.printer.Plain("%s %s", "jmattheis.de", w)
	}
	c.printer.End()
	c.printer.Write(250, "ok")
}

type ShowCommand struct{}

func (cmd ShowCommand) Run(c *Client, args []string) {
	if len(args) != 1 {
		c.printer.IllegalParams()
		return
	}

	word := strings.ToUpper(args[0])
	switch word {
	case "DATABASES":
		fallthrough
	case "DB":
		c.printer.Write(110, "1 databases present - text follows")
		c.printer.Plain("jmattheis.de \"this website\"")
		c.printer.End()
	case "STRAT":
		c.printer.Write(111, "2 strategies available - text follows")
		c.printer.Plain("exact \"does nothing\"")
		c.printer.Plain("prefix \"does nothing\"")
		c.printer.End()
	case "INFO":
		c.printer.Write(112, "database information follows")
		c.printer.Plain("nah nothing special")
		c.printer.End()
	case "SERVER":
		c.printer.Write(114, "server information follows")
		c.printer.Plain(c.tty.Get(""))
		c.printer.End()
	}
	c.printer.Write(250, "ok")
}
