package content

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

func HtmlStart() string {
    htmlFlags := html.CommonFlags | html.HrefTargetBlank
    opts := html.RendererOptions{Flags: htmlFlags}
    renderer := html.NewRenderer(opts)

    return string(markdown.ToHTML([]byte(MdStart()), nil, renderer))
}

func MdStart() string {
    return `\
    # Jannis Mattheis

Hey there! I'm a software engineer from Berlin, Germany. Since 2018, I've
created and maintained multiple privacy focused open-source projects. I enjoy
writing simple, maintainable and testable code. It's a kind of art, some
do paint, I write code (:.

You can contact me via [hello@jmattheis.de](mailto@hello.jmattheis.de) and find
me on [GitHub](https://github.com/jmattheis),
[StackOverflow](https://stackoverflow.com/users/4244993/jmattheis).

This website is available via various protocols, which may or may not be
intended to be used that way. Currently there is support for: dict, dns(tcp),
ftp, gopher, http/https, imap, pop3, ssh, telnet/tcp, websocket and whois.

Try one of the following commands in your terminal:
`+"```"+`
curl   dict://jmattheis.de/show:server
curl    ftp://jmattheis.de
curl gopher://jmattheis.de
curl   http://jmattheis.de
curl  https://jmattheis.de
curl  "imap://jmattheis.de/INBOX;UID=1" -u ":"
curl   pop3://jmattheis.de/1
dict -h       jmattheis.de -I
dig          @jmattheis.de +tcp +short
docker -H     jmattheis.de inspect start -f '{{.Value}}'
netcat        jmattheis.de 23
ssh           jmattheis.de
telnet        jmattheis.de 23
whois -h      jmattheis.de .
wscat -c      jmattheis.de
`+"```"+`

## Projects

[Gotify](https://gotify.net):
- [gotify/server](https://github.com/gotify/server): A simple server for
  sending and receiving messages in real-time per WebSocket.
- [gotify/android](https://github.com/gotify/android): An app that creates push
  notifications for new messages posted to gotify/server.
- [gotify/cli](https://github.com/gotify/cli): A cli for pushing messages to
  gotify/server.

[screego/server](https://github.com/screego/server)
([Demo](https://app.screego.net)): screen sharing for developers.

[traggo/server](https://github.com/traggo/server): A self-hosted tag-based time
tracking server.

???/??? (coming soon): A self-hosted podcast player and manager.

I create and manage these projects in my free time besides my full time job. If
you want to support me, you can donate to me. See
[jmattheis.de/donate](https://jmattheis.de/donate).

## Blog

TODO
    `
}
