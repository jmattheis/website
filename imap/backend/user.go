package backend

import (
	"errors"

	"github.com/emersion/go-imap/backend"
)

type User struct {
	username string
	password string
	mailbox  *Mailbox
}

func (u *User) Username() string {
	return u.username
}

func (u *User) ListMailboxes(subscribed bool) (mailboxes []backend.Mailbox, err error) {
	return []backend.Mailbox{u.mailbox},nil
}

func (u *User) GetMailbox(name string) (mailbox backend.Mailbox, err error) {
	return u.mailbox, nil
}

func (u *User) CreateMailbox(name string) error {
	return errors.New("nah")
}

func (u *User) DeleteMailbox(name string) error {
	return errors.New("nah")
}

func (u *User) RenameMailbox(existingName, newName string) error {
	return errors.New("nah")
}

func (u *User) Logout() error {
	return nil
}
