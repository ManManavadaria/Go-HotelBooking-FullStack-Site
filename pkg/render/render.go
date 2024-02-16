package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/config"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

var PathToTemplate = "./templates/"

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	// if app.Session.Exists(r.Context(), "user_id") {
	// 	td.IsAuthenticated = 1
	// }
	return td
}

// RenderTemplate renders a template
func RenderTemplate(r *http.Request, w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
		// log.Fatalf("errors retriving tamplate catche ,err : %s ,tc:%s", err, tc)
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache", tc)
	}
	buf := new(bytes.Buffer)

	td2 := AddDefaultData(td, r)

	_ = t.Execute(buf, td2)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s*.page.templ", PathToTemplate))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s*.layout.templ", PathToTemplate))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s*.layout.templ", PathToTemplate))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
