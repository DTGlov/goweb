package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/DTGlov/goweb.git/pkg/forms"
	"github.com/DTGlov/goweb.git/pkg/models"
)

//Include a Post field in the templateData struct
type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	Post            *models.Post
	Posts           []*models.Post
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	//Return the empty string if time has zero value
	if t.IsZero() {
		return ""
	}
	//Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

//func to cache the templates
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	//Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	//Use the filepath.Glob func to get a slice of all filepaths  with the extension '.page.tmpl'.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	//Loop through the pages one-by-one.
	for _, page := range pages {
		//Extract the file name(eg. 'home.page.tmpl') from the full path and assign it the name variable
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		//Use the ParseGlob method to add any 'layout templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}
		//Use the ParseGlob method to add any 'partial' templates to the template set(eg . 'footer.partial.html')
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}
		//Add the template set to the cache,using the name of the page like ('home.page.html') as the key.
		cache[name] = ts
	}
	//Return the map.
	return cache, nil
}
