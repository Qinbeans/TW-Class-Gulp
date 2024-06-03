package template

import (
	"backend/logger"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type TemplateEngine struct {
	templates map[string]*template.Template
}

type CompRecipe struct {
	Name      string `json:"name"`
	Component string `json:"component"`
}

type Recipe struct {
	Root                string       `json:"root"`
	PageDefinition      []string     `json:"page_definition"`
	ComponentDefinition []string     `json:"component_definition"`
	Recipes             []CompRecipe `json:"recipes"`
}

func NewEngine(mode string) *TemplateEngine {
	var file *os.File
	var err error

	templates := make(map[string]*template.Template)
	// open "public/recipe.json"
	if mode == "dev" {
		file, err = os.Open("public/recipe.json")
	} else {
		file, err = os.Open("recipe.json")
	}
	if err != nil {
		log.Fatalf("%srecipe.json not found: %s\n", logger.ERROR, err)
	}
	defer file.Close()
	// decode json
	decoder := json.NewDecoder(file)
	recipe := Recipe{}
	err = decoder.Decode(&recipe)
	if err != nil {
		log.Fatal(err)
	}
	if recipe.Root == "" {
		log.Fatal("root is required")
	}
	if mode == "dev" {
		recipe.Root = "tmp" + recipe.Root
	} else {
		recipe.Root = "." + recipe.Root
	}
	if recipe.PageDefinition == nil || len(recipe.PageDefinition) == 0 {
		log.Fatal("page_definition is required")
	}
	if recipe.ComponentDefinition == nil || len(recipe.ComponentDefinition) == 0 {
		log.Fatal("component_definition is required")
	}

	for i, comp := range recipe.ComponentDefinition {
		recipe.ComponentDefinition[i] = fmt.Sprintf("%s/%s.go.html", recipe.Root, comp)
	}

	for i, comp := range recipe.PageDefinition {
		recipe.PageDefinition[i] = fmt.Sprintf("%s/%s.go.html", recipe.Root, comp)
	}

	for _, comp := range recipe.Recipes {
		compName := fmt.Sprintf("%s/%s.go.html", recipe.Root, comp.Component)
		// parse component
		componentList := []string{
			compName,
		}
		componentList = append(componentList, recipe.ComponentDefinition...)
		component := template.Must(template.ParseFiles(componentList...))
		templates[fmt.Sprintf("%s-comp", comp.Name)] = component
		// parse page
		pageList := []string{
			compName,
		}
		pageList = append(pageList, recipe.PageDefinition...)
		page, err := template.ParseFiles(pageList...)
		if err != nil {
			log.Printf("%sparsing page %s: %s", logger.ERROR, comp.Name, err)
			continue
		}
		templates[fmt.Sprintf("%s-page", comp.Name)] = page
	}
	return &TemplateEngine{
		templates: templates,
	}
}

func GetBoosted(r *http.Request) bool {
	_, ok := r.Header["Hx-Boosted"]
	return ok
}

func (t *TemplateEngine) Render(w http.ResponseWriter, name string, data interface{}) error {
	var ctx map[string]any = map[string]any{}
	var ok bool
	if data != nil {
		ctx, ok = data.(map[string]any)
		if !ok {
			log.Printf("%sno template context found: %v\n", logger.ERROR, data)
			return errors.New("no template context found")
		}
	}
	// check if the templates exists
	var tpl *template.Template
	if tpl, ok = t.templates[name]; !ok {
		log.Printf("%stemplate not found: %s\n", logger.ERROR, name)
		return errors.New("template not found")
	}
	err := tpl.ExecuteTemplate(w, "html", ctx)
	if err != nil {
		log.Printf("%serror executing template: %s\n", logger.ERROR, err)
		return err
	}
	return nil
}
