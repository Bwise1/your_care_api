package rest

import (
	"log"
	"net/http"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/bwise1/your_care_api/util"
	"github.com/bwise1/your_care_api/util/tracing"
	"github.com/bwise1/your_care_api/util/values"
	"github.com/go-chi/chi/v5"
)

func (api *API) AuthRoutes() chi.Router {
	mux := chi.NewRouter()
	mux.Method(http.MethodPost, "/login", Handler(api.Login))
	mux.Method(http.MethodPost, "/register", Handler(api.CreateUser))
	return mux
}

func (api *API) Login(_ http.ResponseWriter, r *http.Request) *ServerResponse {

	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	var req model.UserLoginReq
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse login request", values.BadRequestBody, &tc)
	}

	user, status, message, err := api.LoginUser(req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	// api.CreateUser()
	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       user,
	}
}

func (api *API) CreateUser(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// Create user
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	log.Println(tc, "creating user")
	var req model.UserRequest
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse registration request", values.BadRequestBody, &tc)
	}

	user, status, message, err := api.RegisterUser(req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       user,
	}
}
