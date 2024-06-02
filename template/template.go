package template

import (
	"backend/logger"
	"log"
	"net/http"
)

func (t *TemplateEngine) Home(w http.ResponseWriter, r *http.Request) {
	err := t.Render(w, "home-page", map[string]any{
		"title":       "Home",
		"description": "This is the home page",
		"items": []map[string]any{
			{
				"selected": true,
				"url":      "/",
				"label":    "Home",
			},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(logger.INFO, err)
	}
}
