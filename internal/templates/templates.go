package templates

import (
	"github.com/danielmichaels/storeman/internal/store/sqlite"
	"github.com/danielmichaels/storeman/ui"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

type BreadCrumb struct {
	// never include `home`, this is always available in HTML as a default.
	// Home > Containers > Edit is a logical structure.
	Name string
	Href string
}

// TemplateData holds any template data which is passed into the template.
// render and defaultData are added via this struct.
type TemplateData struct {
	Title           string
	Form            any
	IsAuthenticated bool
	CSRFToken       string
	Flash           string
	BreadCrumbs     []BreadCrumb
	Containers      []*sqlite.Container
	Container       *sqlite.Container
	//Items           []*sqlite.Item
	//Item            *sqlite.Item
}

// humanDate creates a human-readable datetime for use as a template filter.
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}
func titleCase(s string) string {
	return cases.Title(language.Und, cases.NoLower).String(s)
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"titleCase": titleCase,
}

// NewTemplateCache stores template data in memory. Creating a template cache
// prevents the disk being read for every invocation of a template in the server.
func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
