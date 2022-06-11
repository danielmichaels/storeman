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

// readParamInt accepts a key and attempts to convert it to an int.
func (app *Server) readParamInt(key string, r *http.Request) (int, error) {
	p := chi.URLParam(r, key)
	i, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return i, nil
}
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

type itemForm struct {
	Name                string `form:"name"`
	Description         string `form:"description"`
	Image               []byte `form:"image"`
	validator.Validator `form:"-"`
}

func (app *Server) handleContainerCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form containerForm
		data := app.newTemplateData(r)

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
			id, err := app.Store.ContainerInsert(form.Title, form.Notes)
			if err != nil {
				app.serverError(w, err)
				return
			}

			// todo flash success

			http.Redirect(w, r, fmt.Sprintf("/containers/%d", id), http.StatusSeeOther)
		}

		crumbs := []templates.BreadCrumb{
			{Name: "Containers", Href: "/"},
			{Name: "Create", Href: "/containers/create"},
		}
		data.BreadCrumbs = crumbs
		data.Form = form
		app.render(w, http.StatusOK, "container-create.tmpl", data)
	}
}

func (app *Server) handleContainerEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := app.readParamInt("id", r)
		if err != nil {
			fmt.Println(err)
			app.serverError(w, err)
			return
		}

		container, err := app.Store.ContainerGet(id)
		if err != nil {
			app.serverError(w, errors.New("invalid container ID supplied"))
			return
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
			row, err := app.Store.ContainerUpdate(form.Title, form.Notes, id)
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
			{Name: "Edit", Href: fmt.Sprintf("/containers/edit/%d", id)},
		}
		data.BreadCrumbs = crumbs
		app.render(w, http.StatusOK, "container-edit.tmpl", data)

	}
}
func (app *Server) handleContainerViewAndAddItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.newTemplateData(r)
		id, err := app.readParamInt("id", r)
		if err != nil {
			app.serverError(w, err)
			return
		}
		container, err := app.Store.ContainerGet(id)
		if err != nil {
			app.notFound(w, r)
			return
		}

		items, err := app.Store.ItemGetAllByContainer(id)
		if err != nil {
			app.serverError(w, errors.New("invalid container ID supplied"))
			return
		}
		data.Container = container
		data.Items = items
		crumbs := []templates.BreadCrumb{
			{Name: "Containers", Href: "/"},
			{Name: "Items", Href: fmt.Sprintf("/containers/%d", id)},
		}
		data.BreadCrumbs = crumbs
		app.render(w, http.StatusOK, "container-view.tmpl", data)
	}
}

func (app *Server) handleItemCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form itemForm
		id, err := app.readParamInt("id", r)
		if err != nil {
			app.serverError(w, errors.New("invalid container ID supplied"))
			return
		}
		data := app.newTemplateData(r)

		container, err := app.Store.ContainerGet(id)
		if err != nil {
			fmt.Println("invalid cont id")
			app.serverError(w, errors.New("invalid container ID supplied"))
			return
		}
		data.Container = container

		if r.Method == http.MethodPost {
			err = app.decodePostForm(r, &form)
			if err != nil {
				fmt.Println("form decode", err)
				app.clientError(w, http.StatusBadRequest)
				return
			}

			form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
			// todo validator location and image in future

			if !form.Valid() {
				//data := app.newTemplateData(r)
				fmt.Println("form validation", err)
				data.Form = form
				app.render(w, http.StatusUnprocessableEntity, "item-create.tmpl", data)
				return
			}

			_, err := app.Store.ItemInsert(id, form.Name, form.Description, []byte("image"))
			if err != nil {
				fmt.Println("insert", err)
				//data := app.newTemplateData(r)
				data.Form = form
				app.render(w, http.StatusUnprocessableEntity, "item-edit.tmpl", data)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/containers/%d", id), http.StatusSeeOther)
		}
		crumbs := []templates.BreadCrumb{
			{Name: "Containers", Href: "/"},
			{Name: "Items", Href: fmt.Sprintf("/containers/%d", id)},
			{Name: "Create", Href: fmt.Sprintf("/containers/%d/items/create", id)},
		}
		data.Container = container
		data.BreadCrumbs = crumbs
		data.Form = form
		app.render(w, http.StatusOK, "item-create.tmpl", data)
	}
}
func (app *Server) handleItemEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
func (app *Server) handleItemDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := app.readParamInt("id", r)
		if err != nil {
			app.serverError(w, errors.New("invalid container ID supplied"))
			return
		}
		itemId, err := app.readParamInt("item", r)
		if err != nil {
			app.serverError(w, errors.New("invalid item ID supplied"))
			return
		}
		data := app.newTemplateData(r)

		container, err := app.Store.ContainerGet(id)
		if err != nil {
			fmt.Println("invalid cont id")
			app.serverError(w, errors.New("invalid container ID supplied"))
			return
		}

		item, err := app.Store.ItemGet(itemId)
		if err != nil {
			app.serverError(w, errors.New("invalid container ID supplied"))
			return
		}
		data.Container = container
		data.Item = item
		crumbs := []templates.BreadCrumb{
			{Name: "Containers", Href: "/"},
			{Name: "Items", Href: fmt.Sprintf("/containers/%d", id)},
			{Name: "Detail", Href: fmt.Sprintf("/containers/%d/items/%d", id, itemId)},
		}
		data.BreadCrumbs = crumbs
		app.render(w, http.StatusOK, "item-view.tmpl", data)
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
