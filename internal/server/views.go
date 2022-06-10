package server

import (
	"errors"
	"fmt"
	"github.com/danielmichaels/storeman/internal/templates"
	"github.com/danielmichaels/storeman/internal/validator"
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

type containerForm struct {
	Title               string `form:"title"`
	Notes               string `form:"notes"`
	Location            string `form:"location"`
	Image               []byte `form:"image"`
	validator.Validator `form:"-"`
}

func (app *Server) handleContainerCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form containerForm

		if r.Method == http.MethodPost {
			err := app.decodePostForm(r, &form)
			if err != nil {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
			// todo validator location and image in future

			if !form.Valid() {
				data := app.newTemplateData(r)
				data.Form = form

				app.render(w, http.StatusUnprocessableEntity, "container-create.tmpl", data)
				return
			}
			// insert
			id, err := app.Store.ContainerInsert(form.Title, form.Notes)
			if err != nil {
				app.serverError(w, err)
				return
			}

			// flash success

			// redirect to new container page
			http.Redirect(w, r, fmt.Sprintf("/containers/%d", id), http.StatusSeeOther)
		}

		data := app.newTemplateData(r)
		crumbs := []templates.BreadCrumb{
			{Name: "Containers", Href: "/"},
			{Name: "Create", Href: "/containers/create"},
		}
		data.BreadCrumbs = crumbs
		data.Form = form
		app.render(w, http.StatusOK, "container-create.tmpl", data)
		return
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
		if err != nil {
			app.serverError(w, errors.New("invalid container ID supplied"))
		}

		var form containerForm
		form.Title = container.Title
		form.Notes = container.Notes
		err = app.decodePostForm(r, &form)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		// todo validator location and image in future

		if !form.Valid() {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "container-edit.tmpl", data)
			return
		}

		if r.Method == http.MethodPost {
			row, err := app.Store.ContainerUpdate(form.Title, form.Notes, i)
			if err != nil {
				data := app.newTemplateData(r)
				data.Form = form
				app.render(w, http.StatusUnprocessableEntity, "container-edit.tmpl", data)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/containers/%d", row), http.StatusSeeOther)
		}

		data := app.newTemplateData(r)
		data.Form = form
		data.Container = container
		crumbs := []templates.BreadCrumb{
			{Name: "Containers", Href: "/"},
			{Name: "Edit", Href: fmt.Sprintf("/containers/edit/%s", id)},
		}
		data.BreadCrumbs = crumbs
		app.render(w, http.StatusOK, "container-edit.tmpl", data)

	}
}
func (app *Server) handleContainerViewAndAddItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		i, err := strconv.Atoi(id)
		if err != nil {
			app.serverError(w, errors.New("invalid container ID supplied"))
		}
		container, err := app.Store.ContainerGet(i)
		if err != nil {
			app.serverError(w, errors.New("invalid container ID supplied"))
		}

		items, err := app.Store.ItemGetAllByContainer(i)
		fmt.Println(items)
		if err != nil {
			fmt.Println(err.Error())
			app.serverError(w, errors.New("invalid container ID supplied"))
		}
		data := app.newTemplateData(r)
		data.Container = container
		data.Items = items
		crumbs := []templates.BreadCrumb{
			{Name: "Containers", Href: "/"},
			{Name: "View", Href: fmt.Sprintf("/containers/%s", id)},
		}
		data.BreadCrumbs = crumbs
		app.render(w, http.StatusOK, "container-view.tmpl", data)
	}
}

func (app *Server) handleItemCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
func (app *Server) handleItemEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := 2
		http.Redirect(w, r, fmt.Sprintf("/containers/edit/%d", id), http.StatusSeeOther)
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
