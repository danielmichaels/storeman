package server

import (
	"bytes"
	"github.com/danielmichaels/storeman/internal/templates"
	"net/http"
)

// addDefaultData is a helper which will pre-fill the templateData struct with
// default information that is used across several templates.
func (app *Server) addDefaultData(td *templates.TemplateData, r *http.Request) *templates.TemplateData {
	if td == nil {
		td = &templates.TemplateData{}
	}
	return td
}

// render is a template rendering helper. It uses a template cache to speed up delivery of templates
func (app *Server) render(w http.ResponseWriter, r *http.Request, name string, td *templates.TemplateData) {
	ts, ok := app.Template[name]
	if !ok {
		http.Error(w, "Template does not exist", 500)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf.WriteTo(w)
}
