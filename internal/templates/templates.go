package templates

import (
	"github.com/danielmichaels/storeman/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use fs.Glob() to get a slice of all the filepaths in the ui.Files embedded filesystem
	// which match the pattern 'html/*page.tmpl'.
	pages, err := fs.Glob(ui.Files, "html/*.page.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		// Use ParseFS() to parse a specific page template from ui.Files
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, page)
		if err != nil {
			return nil, err
		}

		// Collect any 'partials'
		//ts, err = ts.ParseFS(templates.Files, "html/*.partial.tmpl")
		//if err != nil {
		//	return nil, err
		//}
		// Collect any layouts
		ts, err = ts.ParseFS(ui.Files, "html/*.layout.tmpl")
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

type TemplateData struct {
	Title string
}
