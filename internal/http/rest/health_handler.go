package rest

import (
	"github.com/bwise1/your_care_api/util/values"
	"github.com/go-chi/chi/v5"
	"github.com/lucsky/cuid"

	"net/http"
)

func HealthRoutes() chi.Router {
	mux := chi.NewRouter()
	mux.Method(http.MethodGet, "/", Handler(func(w http.ResponseWriter, r *http.Request) *ServerResponse {
		return &ServerResponse{
			Message:    values.Success,
			Status:     values.Success,
			StatusCode: http.StatusOK,
			Data:       cuid.New(),
		}
	}))
	return mux
}
