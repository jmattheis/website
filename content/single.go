package content

import (
	"fmt"
	"strings"

	"github.com/jmattheis/website/assets"
)

type SingleText struct {
	Protocol       string
	Port           string
	Split          string
	ForceBanner    string
	CommandPrefix  string
	CommandSuffix  string
	DisablePadding bool
}

func (i SingleText) Get(command string) string {
	words := strings.Split(command, i.Split)

	if len(words) == 0 {
		return i.Get("")
	}

	prefix := i.CommandPrefix
	if !i.DisablePadding {
		prefix = "  " + prefix
	}

	more := fmt.Sprintf(`
Read more:

%sprojects%s   read about projects
%scat%s        show an image of my cat
%sblog%s       show my blog posts
`, prefix, i.CommandSuffix, prefix, i.CommandSuffix, prefix, i.CommandSuffix)
	switch words[0] {
	case "start":
		fallthrough
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
		result := "Choose a blog post:\n\n"

		for index, entry := range assets.BlogList {
			result += fmt.Sprintf("%sblog%s%d%s     %s\n", prefix, i.Split, index, i.CommandSuffix, entry[2:])
		}
		return result + more
	case "projects":
		return ProjectsTXT + more
	case "cat":
		return Cat + more
	default:
		return fmt.Sprintf(`error: %s doesn't exist
`, words[0]) + more
	}
}

func (i SingleText) Commands() []string {
	return []string{"blog", "projects", "start", "cat"}
}
