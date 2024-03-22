package content

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jmattheis/website/assets"
)

type SingleText struct {
	Split          string
	ForceBanner    string
	CommandPrefix  string
	CommandSuffix  string
	DisablePadding bool
	RemoteAddr     string
}

func (i SingleText) Get(command string) string {
	content, _ := i.GetVerbose(command)
	return content
}

func (i SingleText) GetVerbose(command string) (string, bool) {
	words := strings.Split(command, i.Split)

	if len(words) == 0 {
		return i.GetVerbose("")
	}

	prefix := i.CommandPrefix
	if !i.DisablePadding {
		prefix = "  " + prefix
	}

	more := fmt.Sprintf(`
Read more:

%sblog%s       show my blog posts
%scat%s        show an image of my cat
%sip%s         show your ip
%sprojects%s   read about projects
%stime%s       show current time
`, prefix, i.CommandSuffix, prefix, i.CommandSuffix, prefix, i.CommandSuffix, prefix, i.CommandSuffix, prefix, i.CommandSuffix)
	switch words[0] {
	case "start":
		fallthrough
	case "":
		if i.ForceBanner == "" {
			return StartTXT(rdmBanner()) + more, false
		}
		return StartTXT(i.ForceBanner) + more, false
	case "blog":
		if len(words) > 1 {
			return txtBlog(words[1]) + more, false
		}
		fallthrough
	case "blogs":
		result := "Choose a blog post:\n\n"

		for index, entry := range assets.BlogList {
			result += fmt.Sprintf("%sblog%s%d%s     %s\n", prefix, i.Split, index, i.CommandSuffix, entry[2:])
		}
		return result + more, false
	case "projects":
		return ProjectsTXT + more, false
	case "cat":
		return Cat + more, false
	case "ip":
		if strings.Contains(i.RemoteAddr, ":") {
			host, _, _ := net.SplitHostPort(i.RemoteAddr)
			return host, true
		}
		return i.RemoteAddr, true
	case "time":
		return time.Now().UTC().Format(time.RFC3339), true
	default:
		return fmt.Sprintf(`error: %s doesn't exist
`, words[0]) + more, false
	}
}

func (i SingleText) Commands() []string {
	return []string{"blog", "projects", "start", "cat"}
}
