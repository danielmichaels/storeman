package server

import (
	"errors"
	"github.com/danielmichaels/storeman/internal/validator"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"net/http"
)

type passwordForm struct {
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *Server) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.Form.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}
		return err
	}
	return nil
}

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
