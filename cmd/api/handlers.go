package main

import (
	"app/internal/data"
	"app/internal/validator"
	"net/http"
	"strings"
	"time"
)

// enter requireAuth middleware
// pointer template data set to zero value
// then get user by token
// add user data to pointer template data

// func (app *application) HTML(templateName string, data ...any) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		ui.RenderTemplate(w, r, templateName, data)
// 	}
// }

func (app *application) Render(templateName string, templateData *data.TemplateData) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, templateName, templateData)
	})
}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user := data.User{
		Email:    r.Form["email"][0],
		Login:    r.Form["login"][0],
		Password: r.Form["password"][0],
		Name:     r.Form["name"][0],
	}
	v := validator.New()
	ValidateUser(v, &user)
	if !v.Valid() {
		errMsg := ""
		for k, v := range v.Errors {
			errMsg += k + " " + v + "\n"
		}
		app.render(w, r, "signup.page.html", &data.TemplateData{
			ErrorText: errMsg,
			Code:      409,
		})
		return
	}

	err := app.models.Users.Insert(user)
	if err != nil {
		if strings.Contains(err.Error(), "login") {
			app.render(w, r, "signup.page.html", &data.TemplateData{
				ErrorText: "user with such login already exists",
				Code:      409,
			})
			return
		}
		if strings.Contains(err.Error(), "email") {
			app.render(w, r, "signup.page.html", &data.TemplateData{
				ErrorText: "email already in use",
				Code:      409,
			})
			return
		}

	}

	cookie := &http.Cookie{
		HttpOnly: true,
		Secure:   true,
		Name:     "token",
		Path:     "/",
		Value:    "123",
		Expires:  time.Now().Add(15 * time.Second),
		SameSite: 2,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusCreated)
	// validate and then add the user
	// rw.WriteHeader(http.StatusOK)
	// rw.Write([]byte("User Created"))
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	login := r.Form["login"][0]
	password := r.Form["password"][0]

	if login == "" || password == "" {
		app.render(w, r, "login.page.html", &data.TemplateData{
			ErrorText: "login and/or password empty",
			Code:      409,
		})
		return
	}
	user, err := app.models.Users.GetByLogin(login)
	if err.Error() == "mongo: no documents in result" {
		app.render(w, r, "login.page.html", &data.TemplateData{
			ErrorText: "user not found",
			Code:      400,
		})
		return
	}

	if user.Password != password {
		app.render(w, r, "login.page.html", &data.TemplateData{
			ErrorText: "wrong password",
			Code:      409,
		})
		return
	}

	cookie := &http.Cookie{
		HttpOnly: true,
		Secure:   true,
		Name:     "token",
		Path:     "/",
		Value:    "122",
		Expires:  time.Now().Add(15 * time.Second),
		SameSite: 2,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusOK)
}

// func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
// 	app.render(w, r, "login.page.gohtml", &templateData{
// 		Form: forms.NewForm(nil),
// 	})
// }

// func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		app.clientError(w, http.StatusBadRequest)
// 		return
// 	}

// 	// Check whether the credentials are valid. If they're not, add a generic error
// 	// message to the form errors map and re-display the login page.
// 	form := forms.NewForm(r.PostForm)
// 	id, err := app.users.Authenticate(form.Get("email"),
// 		form.Get("password")) // Using embedded url.Values.Get method
// 	if err != nil {
// 		if errors.Is(err, models.ErrInvalidCredentials) {
// 			form.FormErrors.Add("generic", "Email or Password is incorrect")
// 			app.render(w, r, "login.page.gohtml", &templateData{Form: form})
// 		} else {
// 			app.serverError(w, err)
// 		}
// 		return
// 	}

// 	// Add the ID of the current user to the session, so that they are now "logged in".
// 	app.session.Put(r, "authenticatedUserID", id)

// 	// Redirect the user to the create snippet page.
// 	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
// }

// func (app application) logoutUser(w http.ResponseWriter, r *http.Request) {
// 	// Remove the authenticatedUserID from the session data so that the user is
// 	// 'logged out'.
// 	app.session.Remove(r, "authenticatedUserID")
// 	// Add a flash message to the session to confirm to the user that they've been
// 	// logged out.
// 	app.session.Put(r, "flash", "You've been logged out successfully!")
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
