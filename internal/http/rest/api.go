package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/bwise1/your_care_api/config"
	deps "github.com/bwise1/your_care_api/internal/debs"
	"github.com/bwise1/your_care_api/util/values"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

type Handler func(w http.ResponseWriter, r *http.Request) *ServerResponse

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := h(w, r)
	responseByte, err := json.Marshal(response)
	if err != nil {
		writeErrorResponse(w, err, values.Error, "unable to marshal server response")
		return
	}
	writeJSONResponse(w, responseByte, response.StatusCode)
}

type API struct {
	Server *http.Server
	Config *config.Config
	Deps   *deps.Dependencies
}

func (api *API) Serve() error {
	api.Server = &http.Server{
		Addr:         fmt.Sprintf(":%d", api.Config.Port),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Handler:      api.setUpServerHandler(),
	}
	return api.Server.ListenAndServe()
}

func (api *API) setUpServerHandler() http.Handler {
	mux := chi.NewRouter()

	// Add the RequestTracing middleware to all routes
	mux.Use(RequestTracing)

	mux.Get("/",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		},
	)

	mux.Mount("/health", HealthRoutes())
	mux.Mount("/auth", api.AuthRoutes())
	mux.Mount("/hospitals", api.HospitalRoutes())
	mux.Mount("/lab-tests", api.LabTestRoutes())
	mux.Mount("/appointments", api.AppointmentRoutes())
	mux.Mount("/admin", api.AdminRoutes())
	return mux
}

func (a *API) Shutdown() error {
	// err := a.Deps.DAL.DB.Close()
	// if err != nil {
	// 	return err
	// }

	err := a.Server.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}
