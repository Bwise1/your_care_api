package rest

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bwise1/your_care_api/util/tracing"
	"github.com/bwise1/your_care_api/util/values"
	"github.com/lucsky/cuid"
)

// RequestTracing handles the request tracing context
func RequestTracing(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestSource := r.Header.Get(values.HeaderRequestSource)
		if requestSource == "" {
			errM := errors.New("X-Request-Source is empty")

			writeErrorResponse(w, errM, values.Error, errM.Error())
			return
		}

		requestID := r.Header.Get(values.HeaderRequestID)
		if requestID == "" {
			requestID = cuid.New()
		}

		tracingContext := tracing.Context{
			RequestID:     requestID,
			RequestSource: requestSource,
		}

		ctx = context.WithValue(ctx, values.ContextTracingKey, tracingContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func (api *API) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authorization) != 2 || authorization[0] != "Bearer" {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, "not-authorized")
			return
		}

		claims, err := api.verifyToken(authorization[1], false)
		if err != nil {
			log.Println("error verifyig token", err.Error())
			if err.Error() == "token expired" {
				// Handle the expired token case
				writeErrorResponse(w, err, values.TokenExpired, "token-expired")
				return
			}
			writeErrorResponse(w, err, values.NotAuthorised, "invalid-token")
			return
		}

		dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Println("claims", claims.UserID)
		// Get additional user info from database if needed
		user, err := api.GetUserByID(dbCtx, claims.UserID)
		if err != nil {
			writeErrorResponse(w, err, values.NotAuthorised, "user-not-found")
			return
		}

		// Add minimal information to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user", user) // Add full user object if needed

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
