package html

import (
	"net/http"
	"strings"

	"github.com/gernest/front"
	"github.com/gobuffalo/packr/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/jmattheis/website/content"
)

type Blog struct {
    Title string
    Description string
    Date string
    Content []byte
    URL string
}

func Handler() http.HandlerFunc {
    m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)

    htmlFlags := html.CommonFlags | html.HrefTargetBlank
    opts := html.RendererOptions{Flags: htmlFlags}
    renderer := html.NewRenderer(opts)

    blogLookup := map[string]Blog{}
    displayBlog := []Blog{}

    content.Assets.WalkPrefix("blog/", func(name string, file packr.File) error {
        tags, data, _ := m.Parse(file)
        md := markdown.ToHTML([]byte(data), nil, renderer)
        b := Blog {
            Title: tags["title"].(string),
            Description: tags["description"].(string),
            Date: tags["date"].(string),
            Content: md,
            URL: tags["url"].([]interface{})[0].(string),
        }
        displayBlog = append(displayBlog, b)
        for _, url := range tags["url"].([]interface{}) {
            blogLookup[url.(string)] = b
        }
        return nil
    })

	return func(w http.ResponseWriter, r *http.Request) {
        resource := strings.TrimPrefix(r.URL.Path, "/")
		if strings.HasSuffix(r.URL.Path, ".jpg") {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Header().Set("Cache-Control", "max-age=84600, public")
			c, err := content.Assets.Find(r.URL.Path)
			if err != nil {
                resource = "404";
			} else {
                w.Write(c)
            }
		}
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
			w.Header().Set("Cache-Control", "max-age=84600, public")
			c, err := content.Assets.Find(r.URL.Path)
			if err != nil {
                resource = "404";
			} else {
                w.Write(c)
            }
		}

        title := "Jannis Mattheis"
        description := "I'm a software engineer from Berlin, Germany. Since 2018, I'm creating and maintaining privacy focused open-source projects."
        blogContent := ""

        if resource == "" {
            resource = "index"
        }

        if strings.HasPrefix(resource, "blog") {
            b, ok := blogLookup[resource]
            if !ok {
                resource = "404";
            } else {
                resource = "blog";
                blogContent = string(b.Content)
                description = b.Description
                title = b.Title
            }
        }

		_, err := content.Assets.Find(resource + ".html")
		if err != nil {
            w.Write([]byte("no thanks"))
			return
		}

		err = content.HtmlTemplates.ExecuteTemplate(w, resource + ".html", map[string]interface{}{
			"Title":       title,
			"URL":         "https://jmattheis.de" + r.URL.Path,
			"Description": description,
            "Blogs": displayBlog,
            "Content": blogContent,
		})
	}
}
