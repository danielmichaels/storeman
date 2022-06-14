package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Server) notFound(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusNotFound, "404.tmpl", data)
}
func (app *Server) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
	w.Write([]byte("Method Not Allowed"))
}
func (app *Server) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.Logger.Error().Err(err).Msgf(trace)

	//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (app *Server) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
