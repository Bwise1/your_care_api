package rest

import (
	"net/http"
	"strconv"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/bwise1/your_care_api/util"
	"github.com/bwise1/your_care_api/util/tracing"
	"github.com/bwise1/your_care_api/util/values"
	"github.com/go-chi/chi/v5"
)

func (api *API) HospitalRoutes() chi.Router {
	mux := chi.NewRouter()
	mux.Method(http.MethodGet, "/", Handler(api.GetHospitals))
	mux.Method(http.MethodGet, "/{hospitalID}/lab-tests", Handler(api.GetHospitalLabTests))

	mux.Group(func(r chi.Router) {
		r.Use(api.RequireLogin)
		r.Use(api.RequireAdmin)
		r.Method(http.MethodPost, "/", Handler(api.CreateHospital))
		r.Method(http.MethodDelete, "/{hospitalID}", Handler(api.DeleteHospital))
		//under review
		r.Method(http.MethodPut, "/lab-tests/{labTestID}", Handler(api.UpdateHospitalLabTest))
		r.Method(http.MethodDelete, "/lab-tests/{labTestID}", Handler(api.DeleteHospitalLabTest))
	})

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

func (api *API) CreateHospital(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	var req model.Hospital
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse hospital creation request", values.BadRequestBody, &tc)
	}

	hospital, status, message, err := api.CreateHospital_H(req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       hospital,
	}
}

func (api *API) DeleteHospital(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	hospitalID := chi.URLParam(r, "hospitalID")
	id, err := strconv.Atoi(hospitalID)
	if err != nil {
		return respondWithError(err, "unable to parse id", values.BadRequestBody, &tc)
	}

	status, message, err := api.DeleteHospital_H(id)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

func (api *API) CreateHospitalLabTest(w http.ResponseWriter, r *http.Request) *ServerResponse {
	var req model.HospitalLabTest
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)
	if err := util.DecodeJSONBody(&tc, r.Body, &req); err != nil {
		return respondWithError(err, "Invalid request", values.BadRequestBody, &tc)
	}
	test, status, message, err := api.CreateHospitalLabTest_H(req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}
	return &ServerResponse{Message: message, Status: status, StatusCode: util.StatusCode(status), Data: test}
}

func (api *API) GetAHospitalLabTests(w http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)
	hospitalIDStr := chi.URLParam(r, "hospitalID")
	hospitalID, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		return respondWithError(err, "Invalid hospital ID", values.BadRequestBody, &tc)
	}
	tests, status, message, err := api.GetHospitalLabTests_H(hospitalID)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}
	return &ServerResponse{Message: message, Status: status, StatusCode: util.StatusCode(status), Data: tests}
}

func (api *API) UpdateHospitalLabTest(w http.ResponseWriter, r *http.Request) *ServerResponse {
	var req model.HospitalLabTest
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)
	if err := util.DecodeJSONBody(&tc, r.Body, &req); err != nil {
		return respondWithError(err, "Invalid request", values.BadRequestBody, &tc)
	}
	status, message, err := api.UpdateHospitalLabTest_H(req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}
	return &ServerResponse{Message: message, Status: status, StatusCode: util.StatusCode(status)}
}

func (api *API) DeleteHospitalLabTest(w http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)
	idStr := chi.URLParam(r, "labTestID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return respondWithError(err, "Invalid lab test ID", values.BadRequestBody, &tc)
	}
	status, message, err := api.DeleteHospitalLabTest_H(id)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}
	return &ServerResponse{Message: message, Status: status, StatusCode: util.StatusCode(status)}
}
