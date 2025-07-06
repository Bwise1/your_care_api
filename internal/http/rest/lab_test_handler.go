package rest

import (
	"log"
	"net/http"
	"strconv"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/bwise1/your_care_api/util"
	"github.com/bwise1/your_care_api/util/tracing"
	"github.com/bwise1/your_care_api/util/values"
	"github.com/go-chi/chi/v5"
)

func (api *API) LabTestRoutes() chi.Router {
	mux := chi.NewRouter()
	mux.Method(http.MethodGet, "/", Handler(api.GetAllLabTestsHandler))
	mux.Method(http.MethodGet, "/available", Handler(api.GetAvailableTestsForSelectionHandler))

	// Admin endpoints
	mux.Group(func(r chi.Router) {
		r.Use(api.RequireLogin)
		r.Use(api.RequireAdmin)
		r.Method(http.MethodPost, "/", Handler(api.CreateLabTestHandler))
		r.Method(http.MethodPut, "/{labTestID}", Handler(api.UpdateLabTestHandler))
		r.Method(http.MethodDelete, "/{labTestID}", Handler(api.DeleteLabTestHandler))
	})

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

func (api *API) GetAvailableTestsForSelectionHandler(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)
	tests, status, message, err := api.GetAvailableTestsForSelectionHelper()
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

// Admin: Create Lab Test
func (api *API) CreateLabTestHandler(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)
	var req model.LabTest
	if err := util.DecodeJSONBody(&tc, r.Body, &req); err != nil {
		return respondWithError(err, "Invalid request", values.BadRequestBody, &tc)
	}
	test, status, message, err := api.CreateLabTestHelper(req)
	if err != nil {
		log.Println(err)
		return respondWithError(err, message, status, &tc)
	}
	return &ServerResponse{Message: message, Status: status, StatusCode: util.StatusCode(status), Data: test}
}

// Admin: Update Lab Test
func (api *API) UpdateLabTestHandler(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)
	labTestID := chi.URLParam(r, "labTestID")
	var req model.LabTest
	if err := util.DecodeJSONBody(&tc, r.Body, &req); err != nil {
		return respondWithError(err, "Invalid request", values.BadRequestBody, &tc)
	}
	// Set the ID from the URL
	id, err := strconv.Atoi(labTestID)
	if err != nil {
		return respondWithError(err, "Invalid lab test ID", values.BadRequestBody, &tc)
	}
	req.ID = id
	status, message, err := api.UpdateLabTestHelper(req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}
	return &ServerResponse{Message: message, Status: status, StatusCode: util.StatusCode(status)}
}

// Admin: Delete Lab Test
func (api *API) DeleteLabTestHandler(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)
	labTestID := chi.URLParam(r, "labTestID")
	id, err := strconv.Atoi(labTestID)
	if err != nil {
		return respondWithError(err, "Invalid lab test ID", values.BadRequestBody, &tc)
	}
	status, message, err := api.DeleteLabTestHelper(id)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}
	return &ServerResponse{Message: message, Status: status, StatusCode: util.StatusCode(status)}
}
