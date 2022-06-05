package server

import (
	"github.com/danielmichaels/storeman/ui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"net/http"
)

func (app *Server) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(5))
	r.Use(httplog.RequestLogger(app.Logger))
	r.Use(noSurf)

	r.NotFound(app.notFound)
	r.MethodNotAllowed(app.methodNotAllowed)
	fileServer := http.FileServer(neuteredFileSystem{http.FS(ui.Files)})
	r.Handle("/static/*", fileServer)

	r.Get("/", app.handleHomePage())

	return r
}

func (app *Server) handleHomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "home.page.tmpl", &td.TemplateData{
			Title: "Home",
		})
	}
}
