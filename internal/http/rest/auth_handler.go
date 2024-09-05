package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (api *API) AuthRoutes() chi.Router {
	mux := chi.NewRouter()
	mux.Method(http.MethodGet, "/login", Handler(api.Login))
	// mux.Method(http.MethodPost, "/register", Handler(api.Register))
	return mux
}

func (api *API) Login(_ http.ResponseWriter, r *http.Request) *ServerResponse {

	api.CreateUser()
	return &ServerResponse{
		Message:    "Login",
		Status:     "Success",
		StatusCode: http.StatusOK,
		Payload:    "Login",
	}
}
