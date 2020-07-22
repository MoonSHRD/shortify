package router

import (
	"net/http"

	"github.com/MoonSHRD/shortify/controllers"
	"github.com/MoonSHRD/shortify/repositories"
	"github.com/MoonSHRD/shortify/services"

	"github.com/MoonSHRD/shortify/app"
	"github.com/MoonSHRD/shortify/middlewares"
	"github.com/gorilla/mux"
)

func NewRouter(a *app.App) (*mux.Router, error) {
	r := mux.NewRouter()

	// NOTE Create repositories here
	ur, err := repositories.NewLinksRepository(a)
	if err != nil {
		return nil, err
	}

	// NOTE Create services here
	gs := services.NewLinksService(a, ur)

	// NOTE Create controllers here
	gc := controllers.NewLinksController(a, gs)

	// NOTE Create middlewares here
	authMiddleware := middlewares.NewAuthMiddleware(a)

	r.HandleFunc("/", middlewares.Logger(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Shortify is up and running!"))
		},
	)).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()

	// NOTE Add routes here
	api.HandleFunc("/links", authMiddleware.ProcessRequest(middlewares.Logger(gc.CreateLink))).Methods(http.MethodPost)
	api.HandleFunc("/links", authMiddleware.ProcessRequest(middlewares.Logger(gc.GetLink))).Methods(http.MethodGet)

	return r, nil
}
