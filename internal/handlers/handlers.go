package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/helpers"
	"github.com/triplq/snippetbox/internal/models"
	"github.com/triplq/snippetbox/internal/templates"
	"github.com/triplq/snippetbox/internal/validator"
)

type SnippetForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type UserSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func Home(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			helpers.NotFound(w)
			return
		}

		snippets, err := app.Snippets.Latest()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data := templates.NewTemplateData(app, r)
		data.Snippets = snippets

		err = app.Render(w, http.StatusOK, "home.html", data)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}
}

func ShowSnippet(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			helpers.NotFound(w)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				helpers.NotFound(w)
			} else {
				helpers.ServerError(w, err)
			}
			return
		}

		data := templates.NewTemplateData(app, r)
		data.Snippet = snippet

		err = app.Render(w, http.StatusOK, "view.html", data)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}
}

func CreateSnippet(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewTemplateData(app, r)
		data.Form = SnippetForm{
			Expires: 365,
		}
		err := app.Render(w, http.StatusOK, "create.html", data)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}
}

func CreateSnippetPost(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Form SnippetForm

		err := app.DecodePostForm(r, &Form)
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		Form.CheckError(validator.NotBlank(Form.Title), "title", "Field can not be empty")
		Form.CheckError(validator.MaxChars(Form.Title, 100), "title", "Field can not be large then 100")
		Form.CheckError(validator.NotBlank(Form.Content), "content", "Field can not be empty")
		Form.CheckError(validator.MaxChars(Form.Content, 100), "title", "Field can not be large then 100")
		Form.CheckError(validator.PermittedValue(Form.Expires, 1, 7, 365), "expires", "This field must be 1, 7 or 365")

		if !Form.Valid() {
			data := templates.NewTemplateData(app, r)
			data.Form = Form
			err := app.Render(w, http.StatusUnprocessableEntity, "create.html", data)
			if err != nil {
				helpers.ClientError(w, http.StatusBadRequest)
			}
			return
		}

		id, err := app.Snippets.Insert(Form.Title, Form.Content, Form.Expires)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		app.SessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
	}
}

func SignUp(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewTemplateData(app, r)
		data.Form = new(UserSignUpForm)
		err := app.Render(w, http.StatusOK, "signup.html", data)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}
}

func SignUpPost(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Form UserSignUpForm

		err := app.DecodePostForm(r, &Form)
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		Form.CheckError(validator.NotBlank(Form.Name), "name", "Field can not be empty")
		Form.CheckError(validator.MaxChars(Form.Name, 100), "name", "Field can not be lagre then 100")
		Form.CheckError(validator.NotBlank(Form.Email), "email", "Field can not be empty")
		Form.CheckError(validator.Matches(Form.Email, validator.EmailRX), "email", "Field must be a email")
		Form.CheckError(validator.NotBlank(Form.Password), "password", "Field can not be empty")
		Form.CheckError(validator.MinChars(Form.Password, 8), "password", "Field can not be shorter then 8")

		if !Form.Valid() {
			data := templates.NewTemplateData(app, r)
			data.Form = Form
			err := app.Render(w, http.StatusUnprocessableEntity, "signup.html", data)
			if err != nil {
				helpers.ClientError(w, http.StatusBadRequest)
			}
			return
		}

		err = app.Users.Insert(Form.Name, Form.Email, Form.Password)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				Form.AddError("email", "This email is already used")
				data := templates.NewTemplateData(app, r)
				data.Form = Form
				err := app.Render(w, http.StatusUnprocessableEntity, "signup.html", data)
				if err != nil {
					helpers.ClientError(w, http.StatusBadRequest)
				}
				return
			}
			helpers.ServerError(w, err)
			return
		}

		app.SessionManager.Put(r.Context(), "flash", "User successfully signup")

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func LogIn(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewTemplateData(app, r)
		data.Form = new(UserLoginForm)
		err := app.Render(w, http.StatusOK, "login.html", data)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}
}

func LogInPost(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Form UserLoginForm

		err := app.DecodePostForm(r, &Form)
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		Form.CheckError(validator.NotBlank(Form.Email), "email", "Field can not be empty")
		Form.CheckError(validator.Matches(Form.Email, validator.EmailRX), "email", "Field must be a email")
		Form.CheckError(validator.NotBlank(Form.Password), "password", "Field can not be empty")
		Form.CheckError(validator.MinChars(Form.Password, 8), "password", "Field can not be shorter then 8")

		if !Form.Valid() {
			data := templates.NewTemplateData(app, r)
			data.Form = Form
			err = app.Render(w, http.StatusUnprocessableEntity, "login.html", data)
			if err != nil {
				helpers.ClientError(w, http.StatusBadRequest)
			}
			return
		}

		id, err := app.Users.Authenticate(Form.Email, Form.Password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				Form.AddNonFieldErrors("Email or password is incorrect")
				data := templates.NewTemplateData(app, r)
				data.Form = Form
				err = app.Render(w, http.StatusUnprocessableEntity, "login.html", data)
				if err != nil {
					helpers.ClientError(w, http.StatusBadRequest)
				}
				return
			} else {
				helpers.ServerError(w, err)
			}
			return
		}

		err = app.SessionManager.RenewToken(r.Context())
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		app.SessionManager.Put(r.Context(), "authenticatedUserID", id)
		http.Redirect(w, r, "/snippets/create", http.StatusSeeOther)
	}
}

func LogOutPost(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.SessionManager.RenewToken(r.Context())
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		app.SessionManager.Remove(r.Context(), "authenticatedUserID")
		app.SessionManager.Put(r.Context(), "flash", "You are successfully logged out")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
