package content

import (
	"fmt"
	"strings"
)

type SingleText struct {
	Protocol      string
	Port          string
	Prefix        string
	Split         string
	ForceBanner   string
	CommandPrefix string
}

func (i SingleText) Get(command string) string {
	words := strings.Split(command, i.Split)

	if len(words) == 0 {
		return i.Get("")
	}

	more := fmt.Sprintf(`
Read more:

  %sprojects   read about projects
  %scat        show an image of my cat
  %sblog       show my blog posts
`, i.CommandPrefix, i.CommandPrefix, i.CommandPrefix)
	switch words[0] {
	case "":
		if i.ForceBanner == "" {
			return StartTXT(rdmBanner(), i.Protocol, i.Port) + more
		}
		return StartTXT(i.ForceBanner, i.Protocol, i.Port) + more
	case "blog":
		if len(words) > 1 {
			return txtBlog(words[1]) + more
		}
		fallthrough
	case "blogs":
		return txtBlogs() + fmt.Sprintf(`
Read a post with
  %sblog.<nr>

Example:
  %sblog%s1
`, i.CommandPrefix, i.CommandPrefix, i.Split)
	case "projects":
		return ProjectsTXT + more
	case "cat":
		return Cat + more
	default:
		return fmt.Sprintf(`error: %s doesn't exist`, words[0]) + more
	}
}
