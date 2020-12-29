package content

import (
	"fmt"
	"strings"
)

type InteractiveText struct {
	Prompt   string
	Protocol string
	Port     string
}

func (i InteractiveText) Exec(command string) (string, bool) {
	words := strings.Fields(command)

	if len(words) == 0 {
		return i.Exec("start")
	}

	switch words[0] {
	case "start":
		return fmt.Sprintf(`%s
This is an interactive text interface,
write 'help' to show the available commands.
%s`, StartTXT(rdmBanner(), i.Protocol, i.Port), i.Prompt), false
	case "help":
		return fmt.Sprintf(`
blog      - show a list of my blog posts
blog <id> - show a blog post
cat       - show an ascii image of my cat
donate    - show donation info
exit      - terminate the connection
help      - show help
projects  - show my projects
%s`, i.Prompt), false
	case "blog":
		if len(words) > 1 {
			return txtBlog(words[1]) + "\n" + i.Prompt, false
		}
		fallthrough
	case "blogs":
		return fmt.Sprintf(`
%s
Show a blog post with:

  blog <number>

Example: blog 1
%s`, txtBlogs(), i.Prompt), false
	case "cat":
		return Cat + "\n" + i.Prompt, false
	case "donate":
		return Donate + "\n" + i.Prompt, false
	case "projects":
		return "\n" + ProjectsTXT + "\n" + i.Prompt, false
	case "exit":
		return fmt.Sprintf(`
Thanks for visiting :D.
`), true
	default:
		return fmt.Sprintf(`
error: %s doesn't exist
%s`, words[0], i.Prompt), false
	}
}
