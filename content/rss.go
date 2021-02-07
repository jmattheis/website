package content

import (
	"github.com/gernest/front"
	strip "github.com/grokify/html-strip-tags-go"
	"time"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gorilla/feeds"
)

type Blog struct {
	Title       string
	Description string
	Date        string
	Content     []byte
	URL         string
}

func ParseBlogs() ([]Blog, map[string]Blog) {
	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	blogLookup := map[string]Blog{}
	displayBlog := []Blog{}

	Assets.WalkPrefix("blog/", func(name string, file packr.File) error {
		tags, data, _ := m.Parse(file)
		md := markdown.ToHTML([]byte(data), nil, renderer)
		b := Blog{
			Title:       tags["title"].(string),
			Description: tags["description"].(string),
			Date:        tags["date"].(string),
			Content:     md,
			URL:         tags["url"].([]interface{})[0].(string),
		}
		displayBlog = append(displayBlog, b)
		for _, url := range tags["url"].([]interface{}) {
			blogLookup[url.(string)] = b
		}
		return nil
	})

	return displayBlog, blogLookup
}

func BlogsRss() *feeds.Feed {
    f := feeds.Feed {
        Title: "Jannis Mattheis Blog",
        Author: &feeds.Author{
            Name: "Jannis Mattheis",
            Email: "hello@jmattheis.de",
        },
        Link: &feeds.Link{
            Href: "https://jmattheis.de/",
            Type: "website",
        },
    }
    blogs, _ :=ParseBlogs()
    for _, b := range blogs {
        t, _ := time.Parse("2006-01-02", b.Date)
        f.Items = append(f.Items, &feeds.Item{
            Title: b.Title,
            Author: &feeds.Author{
                Name: "Jannis Mattheis",
                Email: "hello@jmattheis.de",
            },
            Link: &feeds.Link{Href: fmt.Sprintf("https://jmattheis.de/%s", b.URL)},
            Id: fmt.Sprintf("https://jmattheis.de/%s", b.URL),
            Created: t,
            Description: strip.StripTags(string(b.Content))[:500],
            Content: string(b.Content),
        })
    }
    return &f
}
