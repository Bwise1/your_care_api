package rest

import (
	"net/http"
	"strconv"

	"github.com/bwise1/your_care_api/util"
	"github.com/bwise1/your_care_api/util/tracing"
	"github.com/bwise1/your_care_api/util/values"
	"github.com/go-chi/chi/v5"
)

func (api *API) HospitalRoutes() chi.Router {
	mux := chi.NewRouter()
	mux.Method(http.MethodGet, "/", Handler(api.GetHospitals))
	mux.Method(http.MethodGet, "/{hospitalID}/lab-tests", Handler(api.GetHospitalLabTests))
	return mux
}

func (api *API) GetHospitals(w http.ResponseWriter, r *http.Request) *ServerResponse {

	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	hospitals, status, message, err := api.GetHospitals_H()
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       hospitals,
	}
}
func (api *API) GetHospitalLabTests(_ http.ResponseWriter, r *http.Request) *ServerResponse {

	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	hospitalID := chi.URLParam(r, "hospitalID")
	id, err := strconv.Atoi(hospitalID)
	if err != nil {
		return respondWithError(err, "unable to parse id", values.BadRequestBody, &tc)

	}

	tests, status, message, err := api.GetLabTestsByHospital_H(id)
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
