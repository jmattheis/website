package backend

import (
	"errors"
	"time"

	"github.com/emersion/go-imap"
)

var Delimiter = "/"

type Mailbox struct {
	Subscribed bool
	Messages   []*Message

	name string
	user *User
}

func (mbox *Mailbox) Name() string {
	return mbox.name
}

func (mbox *Mailbox) Info() (*imap.MailboxInfo, error) {
	info := &imap.MailboxInfo{
		Delimiter: Delimiter,
		Name:      mbox.name,
	}
	return info, nil
}

func (mbox *Mailbox) uidNext() uint32 {
	var uid uint32
	for _, msg := range mbox.Messages {
		if msg.Uid > uid {
			uid = msg.Uid
		}
	}
	uid++
	return uid
}

func (mbox *Mailbox) flags() []string {
	flagsMap := make(map[string]bool)
	for _, msg := range mbox.Messages {
		for _, f := range msg.Flags {
			if !flagsMap[f] {
				flagsMap[f] = true
			}
		}
	}

	var flags []string
	for f := range flagsMap {
		flags = append(flags, f)
	}
	return flags
}

func (mbox *Mailbox) unseenSeqNum() uint32 {
	for i, msg := range mbox.Messages {
		seqNum := uint32(i + 1)

		seen := false
		for _, flag := range msg.Flags {
			if flag == imap.SeenFlag {
				seen = true
				break
			}
		}

		if !seen {
			return seqNum
		}
	}
	return 0
}

func (mbox *Mailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	status := imap.NewMailboxStatus(mbox.name, items)
	status.Flags = mbox.flags()
	status.PermanentFlags = []string{"\\*"}
	status.UnseenSeqNum = mbox.unseenSeqNum()

	for _, name := range items {
		switch name {
		case imap.StatusMessages:
			status.Messages = uint32(len(mbox.Messages))
		case imap.StatusUidNext:
			status.UidNext = mbox.uidNext()
		case imap.StatusUidValidity:
			status.UidValidity = 1
		case imap.StatusRecent:
			status.Recent = 0 // TODO
		case imap.StatusUnseen:
			status.Unseen = 0 // TODO
		}
	}

	return status, nil
}

func (mbox *Mailbox) SetSubscribed(subscribed bool) error {
	return errors.New("nah")
}

func (mbox *Mailbox) Check() error {
	return nil
}

func (mbox *Mailbox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)

	for i, msg := range mbox.Messages {
		seqNum := uint32(i + 1)

		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = seqNum
		}
		if !seqSet.Contains(id) {
			continue
		}

		m, err := msg.Fetch(seqNum, items)
		if err != nil {
			continue
		}

		ch <- m
	}

	return nil
}

func (mbox *Mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	var ids []uint32
	for i, msg := range mbox.Messages {
		seqNum := uint32(i + 1)

		ok, err := msg.Match(seqNum, criteria)
		if err != nil || !ok {
			continue
		}

		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = seqNum
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (mbox *Mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	return errors.New("nah")
}

func (mbox *Mailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, op imap.FlagsOp, flags []string) error {
	return errors.New("nah")
}

func (mbox *Mailbox) CopyMessages(uid bool, seqset *imap.SeqSet, destName string) error {
	return errors.New("nah")
}

func (mbox *Mailbox) Expunge() error {
	return errors.New("nah")
}
