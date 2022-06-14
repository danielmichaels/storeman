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

	// Public
	r.Get("/", app.handleHomePage())
	r.Get("/login", app.handleLogin())
	r.Post("/login", app.handleLoginPost())

	// Private
	r.Group(func(r chi.Router) {
		r.Get("/containers/create", app.handleContainerCreate())
		r.Post("/containers/create", app.handleContainerCreate())

		r.Get("/containers/{id}/edit", app.handleContainerEdit())
		r.Post("/containers/{id}/edit", app.handleContainerEdit())

		r.Get("/containers/{id}", app.handleContainerViewAndAddItems())
		r.Post("/containers/{id}", app.handleContainerViewAndAddItems())

		r.Get("/containers/{id}/items/create", app.handleItemCreate())
		r.Post("/containers/{id}/items/create", app.handleItemCreate())

		r.Get("/containers/{id}/items/{item}", app.handleItemDetail())
		r.Get("/containers/{id}/items/{item}/edit", app.handleItemEdit())
		r.Post("/containers/{id}/items/{item}/edit", app.handleItemEdit())
	})
	return r
}
