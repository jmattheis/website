package html

import (
	"net/http"
	"strings"

	"github.com/jmattheis/website/content"
)

func Handler() http.HandlerFunc {
	displayBlog, blogLookup := content.ParseBlogs()

	return func(w http.ResponseWriter, r *http.Request) {
		resource := strings.TrimPrefix(r.URL.Path, "/")
		if strings.HasSuffix(r.URL.Path, ".jpg") {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Header().Set("Cache-Control", "max-age=84600, public")
			c, err := content.Assets.Find(r.URL.Path)
			if err != nil {
				resource = "404"
			} else {
				w.Write(c)
				return
			}
		}
		if strings.HasSuffix(r.URL.Path, ".ico") {
			w.Header().Set("Content-Type", "image/x-icon")
			w.Header().Set("Cache-Control", "max-age=84600, public")
			c, err := content.Assets.Find(r.URL.Path)
			if err != nil {
				resource = "404"
			} else {
				w.Write(c)
				return
			}
		}
		if strings.HasSuffix(r.URL.Path, ".svg") {
			w.Header().Set("Content-Type", "image/svg+xml")
			w.Header().Set("Cache-Control", "max-age=84600, public")
			c, err := content.Assets.Find(r.URL.Path)
			if err != nil {
				resource = "404"
			} else {
				w.Write(c)
				return
			}
		}
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
			w.Header().Set("Cache-Control", "max-age=84600, public")
			c, err := content.Assets.Find(r.URL.Path)
			if err != nil {
				resource = "404"
			} else {
				w.Write(c)
				return
			}
		}

		title := "Jannis Mattheis"
		description := "I'm a software engineer from Berlin, Germany. Since 2018, I'm creating and maintaining privacy focused open-source projects."
		blogContent := ""

		if resource == "" {
			resource = "index"
		}

		if resource == "404" {
			w.WriteHeader(404)
			return
		}

		if strings.HasPrefix(resource, "blog") {
			b, ok := blogLookup[resource]
			if !ok {
				resource = "404"
			} else {
				resource = "blog"
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

		err = content.HtmlTemplates.ExecuteTemplate(w, resource+".html", map[string]interface{}{
			"Title":       title,
			"URL":         "https://jmattheis.de" + r.URL.Path,
			"Description": description,
			"Blogs":       displayBlog,
			"Content":     blogContent,
		})
	}
}
