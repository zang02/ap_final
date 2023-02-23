package main

import (
	"app/internal/data"
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.PrintInfo(fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI()), "")
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// trigger to make Go's http server automatically close
				w.Header().Set("Connection", "close")

				// Call the app.serverError helper method to return a 500 status code.
				// Also, panic returns an interface{}, so we normalize error into an
				// error object by using the fmt.Errorf() function (which app.serverError expects).
				// using fmt.Errorf() with err will create a new error object containing the default
				// textual representation of the interface{} value panic returns.
				// TODO: app.serverError(w, fmt.Errorf("%s", err))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%s", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and return
		// from the middleware chain so that no subsequent hanlders in the chain are executed
		// if !app.isAuthenticated(r) {
		// 	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		// 	return
		// }
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		fmt.Println(tokenCookie)
		// Otherwise set the "Cache-Control: no-store" header so that pages that
		// require authentication are not stored in the users browsers cache (or other
		// intermediary cache)
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// // // noSurf uses customized CSRF cookie with the Secure, Path and HttpOnly flags set.
// // func noSurf(next http.Handler) http.Handler {
// // 	csrfHandler := nosurf.New(next)
// // 	csrfHandler.SetBaseCookie(http.Cookie{
// // 		HttpOnly: true,
// // 		Path:     "/",
// // 		Secure:   true,
// // 	})

// // 	return csrfHandler
// // }

func (app *application) authenticate(td *data.TemplateData, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		td.User.Email = ""
		td.User.Login = ""
		td.User.Name = ""
		td.User.CreateDate = ""
		td.IsAuthenticated = false
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		TokenDocument, err := app.models.Tokens.GetTokenDocumentByToken(tokenCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := app.models.Users.GetByLogin(TokenDocument.UserLogin)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		} else {
			td.User.Email = user.Email
			td.User.Login = user.Login
			td.User.Name = user.Name
			td.User.CreateDate = user.CreateDate
			td.IsAuthenticated = true
		}
		next.ServeHTTP(w, r)
	})
}
