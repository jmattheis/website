package dict

import (
	"bufio"
	"fmt"
	"github.com/jmattheis/website/content"
	"net"
	"strings"
	"time"
)

type Client struct {
	commands map[string]Executable
	printer  *Printer
	isAlive  bool
	tty      *content.SingleText
}

func newClient() *Client {
	commands := make(map[string]Executable)

	commands["QUIT"] = QuitCommand{}
	commands["SHOW"] = ShowCommand{}
	commands["HELP"] = HelpCommand{}
	commands["DEFINE"] = DefineCommand{}
	commands["MATCH"] = MatchCommand{}
	commands["STATUS"] = StatusCommand{}
	commands["CLIENT"] = ClientCommand{}
	commands["OPTION"] = NotImplementedCommand{}
	commands["AUTH"] = NotImplementedCommand{}
	commands["SASLAUTH"] = NotImplementedCommand{}

	text := content.SingleText{
		Split:         ".",
		CommandPrefix: "dict -h jmattheis.de ",
	}

	return &Client{
		commands: commands,
		tty:      &text,
	}
}

func (c Client) handle(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
	c.printer = NewPrinter(conn)

	c.isAlive = true
	reader := bufio.NewReader(conn)

	c.printer.Write(220, "jmattheis.de <guest@jmattheis.de>")

	for c.isAlive {
		input, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		cmd, args := c.parseInput(string(input))
		exec, ok := c.commands[cmd]
		if !ok {
			c.printer.Write(500, "Syntax error, command not recognized")
			continue
		}
		exec.Run(&c, args)
	}
}

func (c Client) parseInput(input string) (string, []string) {
	cmd := strings.Split(input, " ")
	return strings.ToUpper(cmd[0]), cmd[1:]
}

type Server struct {
	listener net.Listener
}

func NewServer() *Server {
	return &Server{
	}
}

func (s Server) Start(port string) error {

	var err error
	s.listener, err = net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}

			c := newClient()
			go c.handle(conn)
		}
	}()

	return nil
}

//---------------PRINTER

type Printer struct {
	conn net.Conn
}

func NewPrinter(conn net.Conn) *Printer {
	return &Printer{conn}
}

func (p Printer) Write(status int, message string, a ...interface{}) {
	fmt.Fprintf(p.conn, "%d %s\r\n", status, fmt.Sprintf(message, a...))
}

func (p Printer) Plain(message string, a ...interface{}) {
	fmt.Fprintf(p.conn, "%s\r\n", fmt.Sprintf(message, a...))
}

func (p Printer) End() {
	fmt.Fprint(p.conn, ".\r\n")
}

func (p Printer) IllegalParams() {
	p.Write(501, "Syntax error, illegal parameters")
}

func (p Printer) NotImplemented() {
	p.Write(502, "Command not implemented")
}

func (p Printer) Ok(msg string, a ...interface{}) {
	fmt.Fprintf(p.conn, "+OK %s\r\n", fmt.Sprintf(msg, a...))
}

func (p Printer) Err(msg string, a ...interface{}) {
	fmt.Fprintf(p.conn, "-ERR %s\r\n", fmt.Sprintf(msg, a...))
}
