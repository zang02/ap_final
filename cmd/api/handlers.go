package main

import (
	"app/internal/data"
	"net/http"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// set lastest as home page
	// s, err := app.snippets.Latest()
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	app.render(w, r, "home.page.html", &data.TemplateData{

		Envelope: data.Envelope{
			"hello": "World",
		},
	})
}
func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	err := app.models.Users.Insert(
		r.Form["email"][0],
		r.Form["login"][0],
		r.Form["password"][0],
		r.Form["name"][0],
	)
	if err != nil {
		app.serverError(w, err)
		return
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

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.html", &data.TemplateData{})
}

// func getSignedToken() (string, error) {
// 	// we make a JWT Token here with signing method of ES256 and claims.
// 	// claims are attributes.
// 	// aud - audience
// 	// iss - issuer
// 	// exp - expiration of the Token
// 	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 	// 	"aud": "frontend.knowsearch.ml",
// 	// 	"iss": "knowsearch.ml",
// 	// 	"exp": string(time.Now().Add(time.Minute * 1).Unix()),
// 	// })
// 	claimsMap := map[string]string{
// 		"aud": "frontend.knowsearch.ml",
// 		"iss": "knowsearch.ml",
// 		"exp": fmt.Sprint(time.Now().Add(time.Minute * 1).Unix()),
// 	}
// 	// here we provide the shared secret. It should be very complex.\
// 	// Aslo, it should be passed as a System Environment variable

// 	secret := "S0m3_R4n90m_sss"
// 	header := "HS256"
// 	tokenString, err := jwt.GenerateToken(header, claimsMap, secret)
// 	if err != nil {
// 		return tokenString, err
// 	}
// 	return tokenString, nil
// }

// searches the user in the database.
// func validateUser(email string, passwordHash string) (bool, error) {

// 	usr, exists := data.GetUserObject(email)
// 	if !exists {
// 		return false, errors.New("user does not exist")
// 	}
// 	passwordCheck := usr.ValidatePasswordHash(passwordHash)

// 	if !passwordCheck {
// 		return false, nil
// 	}
// 	return true, nil
// }

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.html", &data.TemplateData{})
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {

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
