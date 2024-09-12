package rest

import (
	"github.com/go-chi/chi/v5"
)

func (api *API) AppointmentRoutes() chi.Router {
	mux := chi.NewRouter()
	// mux.Method(http.MethodPost, "/login", Handler(api.Login))
	// mux.Method(http.MethodPost, "/register", Handler(api.CreateUser))
	return mux
}
