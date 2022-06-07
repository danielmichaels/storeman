package server

import (
	"errors"
	"fmt"
	"github.com/danielmichaels/storeman/internal/templates"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *Server) handleHomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		containers, err := app.Store.ContainerGetAll()
		if err != nil {
			app.serverError(w, err)
		}
		data := app.newTemplateData(r)
		data.Containers = containers
		app.render(w, http.StatusOK, "home.tmpl", data)
	}
}
func (app *Server) handleContainerCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.newTemplateData(r)
		app.render(w, http.StatusOK, "create.tmpl", data)
	}
}

func (app *Server) handleContainerEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		i, err := strconv.Atoi(id)
		if err != nil {
			app.serverError(w, errors.New("invalid container ID supplied"))
		}
		container, err := app.Store.ContainerGet(i)
		data := app.newTemplateData(r)
		data.Container = container
		crumbs := []templates.BreadCrumb{
			{Name: "Containers", Href: "/"},
			{Name: "Edit", Href: fmt.Sprintf("/containers/edit/%s", id)},
		}
		data.BreadCrumbs = crumbs
		app.render(w, http.StatusOK, "edit.tmpl", data)
	}
}

func (app *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.newTemplateData(r)
		data.Form = passwordForm{}
		app.render(w, http.StatusOK, "login.tmpl", data)
	}
}

func (app *Server) handleLoginPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form passwordForm
		err := app.decodePostForm(r, &form)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// todo(ds) authentication
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
