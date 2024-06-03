package template

import (
	"backend/logger"
	"log"
	"net/http"
)

func (t *TemplateEngine) Home(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := map[string]any{
		"title":       "Home",
		"description": "This is the home page",
		"items": []map[string]any{
			{
				"selected": true,
				"url":      "/",
				"label":    "Home",
			},
			{
				"selected": false,
				"url":      "/something",
				"label":    "Something",
			},
		},
	}
	if GetBoosted(r) {
		err = t.Render(w, "home-comp", ctx)
	} else {
		err = t.Render(w, "home-page", ctx)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(logger.INFO, err)
	}
}

func (t *TemplateEngine) Something(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := map[string]any{
		"title":       "Something",
		"description": "This is the something page",
		"items": []map[string]any{
			{
				"selected": false,
				"url":      "/",
				"label":    "Home",
			},
			{
				"selected": true,
				"url":      "/something",
				"label":    "Something",
			},
		},
	}
	if GetBoosted(r) {
		err = t.Render(w, "something-comp", ctx)
	} else {
		err = t.Render(w, "something-page", ctx)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(logger.INFO, err)
	}
}
