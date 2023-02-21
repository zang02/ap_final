package main

import (
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

// func (app *application) requireAuth(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// If the user is not authenticated, redirect them to the login page and return
// 		// from the middleware chain so that no subsequent hanlders in the chain are executed
// 		if !app.isAuthenticated(r) {
// 			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
// 			return
// 		}

// 		// Otherwise set the "Cache-Control: no-store" header so that pages that
// 		// require authentication are not stored in the users browsers cache (or other
// 		// intermediary cache)
// 		w.Header().Add("Cache-Control", "no-store")
// 		next.ServeHTTP(w, r)
// 	})
// }

// // noSurf uses customized CSRF cookie with the Secure, Path and HttpOnly flags set.
// func noSurf(next http.Handler) http.Handler {
// 	csrfHandler := nosurf.New(next)
// 	csrfHandler.SetBaseCookie(http.Cookie{
// 		HttpOnly: true,
// 		Path:     "/",
// 		Secure:   true,
// 	})

// 	return csrfHandler
// }

// func (app *application) authenticate(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Check if authenticatedUserID value exists in the session. If this *is not present*
// 		// then call the next handler in the chain as normal
// 		exists := app.session.Exists(r, "authenticatedUserID")
// 		if !exists {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		// Fetch the details of the current user from the database based on session
// 		// authenticatedUserID. If no matching record is found,
// 		// or the current user has been deactivated, remove the (invalid) authenticatedUserID value
// 		// from their session and call the next handler in the chain as normal.
// 		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
// 		if errors.Is(err, models.ErrNoRecord) || !user.Active {
// 			app.session.Remove(r, "authenticatedUserID")
// 			next.ServeHTTP(w, r)
// 			return
// 		} else if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}

// 		// Otherwise, we know the that the request is coming from an activated, authenticated user.
// 		// We create a new copy of the request, with a true boolean value added to the request
// 		// context to indicate this, and call the next handler in the chain *using this new copy of
// 		// the request*.
// 		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
