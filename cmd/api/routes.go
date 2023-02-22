package main

import (
	"app/internal/data"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our app receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// add the authenticate() middleware to the chain
	// and use the noSurf middleware on all our dynamic routes.

	// dynamicMiddleware := alice.New(app.requireAuth)

	// dynamicMiddleware := alice.New()

	r := mux.NewRouter()
	templateData := &data.TemplateData{}

	// Register exact matches before wildcard route match (i.e. :id in Get method for
	// '/snippet/create').
	// Update these routes to use the dynamic middleware chain follow by the appropriate handler
	// function.
	// r.Handle("/", dynamicMiddleware.ThenFunc(app.home)).Methods("GET")
	r.Handle("/", app.authenticate(templateData, app.Render("home.page.html", templateData))).Methods("GET")

	// Require auth middleware for auth'd/logged-in actions
	// TODO: mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.createSnippetForm))
	// TODO: mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// Require auth middleware for auth'd/logged-in actions
	// TODO: mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.createSnippet))

	// Add the five new routes for user authentication.
	r.Handle("/user/signup", app.Render("signup.page.html", templateData)).Methods("GET")
	r.Handle("/user/login", app.Render("login.page.html", templateData)).Methods("GET")
	r.HandleFunc("/user/signup", app.signupHandler).Methods("POST")
	r.HandleFunc("/user/login", app.loginHandler).Methods("POST")

	// Require auth middleware for auth'd/logged-in actions
	// r.Handle("/user/logout", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.logoutUser))

	// gorilla mux file server
	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/").Handler(http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(r)
}
