// A memory backend.
package backend

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmattheis/website/assets"
	"github.com/jmattheis/website/content"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type Backend struct {
	user *User
}

func (be *Backend) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {
	return be.user, nil
}

func New(port string) *Backend {
	user := &User{username: "guest", password: "abcde"}

	blog := `Choose a blog post:
`
	for i, entry := range assets.BlogList {
		blog += fmt.Sprintf("  curl -u \":\" \"imap://jmattheis.de/INBOX;UID=%d\"    %s\n", i+5, entry[2:])
	}

	user.mailbox = &Mailbox{
		name: "INBOX",
		user: user,
		Messages: []*Message{
			build(1, content.StartTXT(content.DnsSafeBanner, "imap", port)),
			build(2, content.ProjectsTXT),
			build(3, content.Cat),
			build(4, blog),
		},
	}

	for i, content := range assets.BlogContent {
		user.mailbox.Messages = append(user.mailbox.Messages, build(i+5, content))
	}
	return &Backend{
		user: user,
	}
}

func build(msgId int, content string) *Message {
	value := `From: hello@jmattheis.de
To: hello@jmattheis.de
Subject: jmattheis.de
Date: Wed, 11 May 2016 14:31:59 +0000
Content-Type: text/plain

`
	value += content
	value += `
Read more:

  curl -u ":" "imap://jmattheis.de/INBOX;UID=2"    read about projects
  curl -u ":" "imap://jmattheis.de/INBOX;UID=3"    show an image of my cat
  curl -u ":" "imap://jmattheis.de/INBOX;UID=4"    show my blog posts
`
	return &Message{
		Uid:   uint32(msgId),
		Date:  time.Now(),
		Size:  uint32(len(value)),
		Flags: []string{"\\Seen"},
		Body:  []byte(strings.ReplaceAll(value, "\n", "\r\n")),
	}
}
