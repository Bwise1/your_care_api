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
	mux.Method(http.MethodPost, "/register-admin", Handler(api.CreateAdminUser))

	mux.Method(http.MethodPost, "/token/refresh", Handler(api.TokenRefresh))

	// mux.Method(http.MethodPost, "/verify-email", Handler(api.VerifyEmail))
	mux.Method(http.MethodPost, "/resend-verification", Handler(api.ResendVerification))
	// mux.Method(http.MethodPost, "/forgot-password", Handler(api.ForgotPassword))
	// mux.Method(http.MethodPost, "/reset-password", Handler(api.ResetPassword))
	mux.Group(func(r chi.Router) {
		r.Use(api.RequireLogin)
		r.Method(http.MethodPut, "/change-password", Handler(api.ChangePassword))
		r.Method(http.MethodPost, "/logout", Handler(api.Logout))
	})

	// mux.Method(http.MethodDelete, "/delete-account", Handler(api.DeleteAccount))
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

func (api *API) TokenRefresh(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	var req model.RefreshTokenReq
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse token refresh request", values.BadRequestBody, &tc)
	}

	if req.RefreshToken == "" {
		return respondWithError(nil, "refresh token is required", values.BadRequestBody, &tc)
	}

	token, status, message, err := api.RefreshToken(req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       token,
	}
}

func (api *API) Logout(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	userID := r.Context().Value("user_id")
	// Invalidate the refresh token in the database
	err := api.invalidateRefreshToken(r.Context(), userID.(int))
	if err != nil {
		return respondWithError(err, "failed to logout", values.Error, &tc)
	}

	return &ServerResponse{
		Message:    "logged out successfully",
		Status:     values.Success,
		StatusCode: util.StatusCode(values.Success),
	}
}

// // VerifyEmail handles email verification using a token
// func (api *API) VerifyEmail(_ http.ResponseWriter, r *http.Request) *ServerResponse {
// 	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

// 	var req model.EmailVerificationReq
// 	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
// 		return respondWithError(decodeErr, "unable to parse verification request", values.BadRequestBody, &tc)
// 	}

// 	status, message, err := api.VerifyUserEmail(req)
// 	if err != nil {
// 		return respondWithError(err, message, status, &tc)
// 	}

// 	return &ServerResponse{
// 		Message:    message,
// 		Status:     status,
// 		StatusCode: util.StatusCode(status),
// 	}
// }

// ResendVerification handles resending verification email
func (api *API) ResendVerification(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	var req model.ResendVerificationReq
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse resend verification request", values.BadRequestBody, &tc)
	}

	status, message, err := api.ResendVerificationEmail(req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

// // ForgotPassword initiates the password reset process
// func (api *API) ForgotPassword(_ http.ResponseWriter, r *http.Request) *ServerResponse {
// 	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

// 	var req model.ForgotPasswordReq
// 	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
// 		return respondWithError(decodeErr, "unable to parse forgot password request", values.BadRequestBody, &tc)
// 	}

// 	status, message, err := api.InitiatePasswordReset(req)
// 	if err != nil {
// 		return respondWithError(err, message, status, &tc)
// 	}

// 	return &ServerResponse{
// 		Message:    message,
// 		Status:     status,
// 		StatusCode: util.StatusCode(status),
// 	}
// }

// // ResetPassword handles password reset using a token
// func (api *API) ResetPassword(_ http.ResponseWriter, r *http.Request) *ServerResponse {
// 	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

// 	var req model.ResetPasswordReq
// 	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
// 		return respondWithError(decodeErr, "unable to parse reset password request", values.BadRequestBody, &tc)
// 	}

// 	status, message, err := api.CompletePasswordReset(req)
// 	if err != nil {
// 		return respondWithError(err, message, status, &tc)
// 	}

// 	return &ServerResponse{
// 		Message:    message,
// 		Status:     status,
// 		StatusCode: util.StatusCode(status),
// 	}
// }

// ChangePassword allows authenticated users to change their password
func (api *API) ChangePassword(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	userID := r.Context().Value("user_id")

	var req model.ChangePasswordReq
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse change password request", values.BadRequestBody, &tc)
	}

	status, message, err := api.ChangeUserPassword(userID.(int), req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

// // DeleteAccount allows users to delete their account
// func (api *API) DeleteAccount(_ http.ResponseWriter, r *http.Request) *ServerResponse {
// 	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

// 	var req model.DeleteAccountReq
// 	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
// 		return respondWithError(decodeErr, "unable to parse delete account request", values.BadRequestBody, &tc)
// 	}

// 	status, message, err := api.RemoveAccount(req)
// 	if err != nil {
// 		return respondWithError(err, message, status, &tc)
// 	}

// 	return &ServerResponse{
// 		Message:    message,
// 		Status:     status,
// 		StatusCode: util.StatusCode(status),
// 	}
// }

func (api *API) CreateAdminUser(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	log.Println(tc, "creating admin user")
	var req model.UserRequest
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse registration request", values.BadRequestBody, &tc)
	}

	user, status, message, err := api.RegisterAdminUser(req)
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
