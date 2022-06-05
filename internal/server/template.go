package server

import (
	"bytes"
	"fmt"
	"github.com/danielmichaels/storeman/internal/templates"
	"github.com/justinas/nosurf"
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
func (app *Server) render(w http.ResponseWriter, status int, name string, td *templates.TemplateData) {
	ts, ok := app.Template[name]
	if !ok {
		err := fmt.Errorf("the template %q does not exist", name)
		app.serverError(w, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", td)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *Server) newTemplateData(r *http.Request) *templates.TemplateData {
	return &templates.TemplateData{
		Containers: []string{},
		//Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: false,
		CSRFToken:       nosurf.Token(r),
	}
}
