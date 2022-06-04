package html

import (
	"net/http"
	"strings"

	"github.com/jmattheis/website/content"
)

func Handler() http.HandlerFunc {
	displayBlog, blogLookup := content.ParseBlogs()

	return func(w http.ResponseWriter, r *http.Request) {
		title := "Jannis Mattheis"
		description := "I'm a software engineer from Berlin, Germany. Since 2018, I'm creating and maintaining privacy focused open-source projects."
		blogContent := ""
		url := "https://jmattheis.de"

		resource := strings.TrimPrefix(r.URL.Path, "/")
		if resource == "" {
			resource = "index"
		}

		if strings.HasPrefix(resource, "blog") {
			b, ok := blogLookup[resource]
			if !ok {
				resource = "404"
				url += "/404"
			} else {
				resource = "blog"
				blogContent = string(b.Content)
				description = b.Description
				title = b.Title
				url += "/" + b.URL
			}
		}

		_, err := content.Assets.Find(resource + ".html")
		if err != nil || resource == "resume" {
			resource = "404"
			w.WriteHeader(404)
		}

		err = content.HtmlTemplates.ExecuteTemplate(w, resource+".html", map[string]interface{}{
			"Title":       title,
			"URL":         url,
			"Description": description,
			"Blogs":       displayBlog,
			"Content":     blogContent,
		})
	}
}
