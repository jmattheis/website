package backends

import (
	"fmt"
	"github.com/jmattheis/website/content"
)

// NoAuth is a fake authorizator interface implementation used for test
type NoAuth struct {
}

// Authorize user for given username and password.
func (a NoAuth) Authorize(user, pass string) bool {
	return true
}

// ContentProvider is a fake backend interface implementation used for test
type ContentProvider struct {
	Port string
}

// Returns total message count and total mailbox size in bytes (octets).
// Deleted messages are ignored.
func (b ContentProvider) Stat(user string) (messages, octets int, err error) {
	return 5, 50, nil
}

// List of sizes of all messages in bytes (octets)
func (b ContentProvider) List(user string) (octets []int, err error) {
	return []int{10, 10, 10, 10, 10}, nil
}

// Returns whether message exists and if yes, then return size of the message in bytes (octets)
func (b ContentProvider) ListMessage(user string, msgId int) (exists bool, octets int, err error) {
	if msgId > 4 {
		return false, 0, nil
	}
	return true, 10, nil
}

// Retrieve whole message by ID - note that message ID is a message position returned
// by List() function, so be sure to keep that order unchanged while client is connected
// See Lock() function for more details
func (b ContentProvider) Retr(user string, msgId int) (message string, err error) {
	value := b.content(msgId)
	return value + `
Read more:

  curl pop3://jmattheis.de/2    read about projects
  curl pop3://jmattheis.de/3    show an image of my cat
  curl pop3://jmattheis.de/4    show my blog posts
`, nil
}

func (b ContentProvider) content(msgId int) string {
	switch (msgId) {
	case 1:
		return content.StartTXT(content.DnsSafeBanner, "pop3", b.Port)
	case 2:
		return content.ProjectsTXT
	case 3:
		return content.Cat
	case 4:
		result := `Choose a blog post:

`
		for i, entry := range content.BlogBox.List() {
			result += fmt.Sprintf("  curl pop3://jmattheis.de/%d    %s\n", i+5, entry[2:])
		}

		return result
	default:
		return content.TXTBlogByNR(msgId - 5)
	}

}

// Delete message by message ID - message should be just marked as deleted until
// Update() is called. Be aware that after Dele() is called, functions like List() etc.
// should ignore all these messages even if Update() hasn't been called yet
func (b ContentProvider) Dele(user string, msgId int) error {
	return nil
}

// Undelete all messages marked as deleted in single connection
func (b ContentProvider) Rset(user string) error {
	return nil
}

// List of unique IDs of all message, similar to List(), but instead of size there
// is a unique ID which persists the same across all connections. Uid (unique id) is
// used to allow client to be able to keep messages on the server.
func (b ContentProvider) Uidl(user string) (uids []string, err error) {
	return []string{"1", "2", "3", "4", "5"}, nil
}

// Similar to ListMessage, but returns unique ID by message ID instead of size.
func (b ContentProvider) UidlMessage(user string, msgId int) (exists bool, uid string, err error) {
	if msgId > 4 {
		return false, "", nil
	}
	return true, fmt.Sprintf("%d", msgId+1), nil
}

// Write all changes to persistent storage, i.e. delete all messages marked as deleted.
func (b ContentProvider) Update(user string) error {
	return nil
}

// Lock is called immediately after client is connected. The best way what to use Lock() for
// is to read all the messages into cache after client is connected. If another user
// tries to lock the storage, you should return an error to avoid data race.
func (b ContentProvider) Lock(user string) error {
	return nil
}

// Release lock on storage, Unlock() is called after client is disconnected.
func (b ContentProvider) Unlock(user string) error {
	return nil
}
