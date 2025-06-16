package rest

import (
	"net/http"

	"github.com/bwise1/your_care_api/util"
	"github.com/bwise1/your_care_api/util/tracing"
	"github.com/bwise1/your_care_api/util/values"
	"github.com/go-chi/chi/v5"
)

func (api *API) LabTestRoutes() chi.Router {
	mux := chi.NewRouter()
	mux.Method(http.MethodGet, "/", Handler(api.GetAllLabTestsHandler))

	return mux
}

func (api *API) GetAllLabTestsHandler(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)
	tests, status, message, err := api.GetAllLabTestsHelper()
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}
	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       tests,
	}
}
