package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (api *API) AdminRoutes() chi.Router {
	mux := chi.NewRouter()

	// Admin appointment routes
	mux.Route("/appointments", func(r chi.Router) {
		r.Use(api.RequireLogin)
		r.Use(api.RequireAdmin)
		r.Method(http.MethodGet, "/", Handler(api.AdminFetchAllAppointments))
		r.Method(http.MethodGet, "/{id}", Handler(api.AdminGetAppointmentDetails))
		r.Method(http.MethodGet, "/{id}/history", Handler(api.AdminGetAppointmentHistory))
		r.Method(http.MethodPost, "/{id}/confirm", Handler(api.AdminConfirmAppointment))
		r.Method(http.MethodPost, "/{id}/reject", Handler(api.AdminRejectAppointment))
		r.Method(http.MethodPost, "/{id}/reschedule", Handler(api.AdminRescheduleAppointment))
		r.Method(http.MethodPost, "/{id}/cancel", Handler(api.AdminCancelAppointment))
		r.Method(http.MethodPut, "/{id}/notes", Handler(api.AdminUpdateNotes))
	})

	// Admin lab test routes
	mux.Route("/tests", func(r chi.Router) {
		r.Use(api.RequireLogin)
		r.Use(api.RequireAdmin)
		r.Method(http.MethodGet, "/", Handler(api.GetAllLabTestsHandler))
		r.Method(http.MethodGet, "/available", Handler(api.GetAvailableTestsForSelectionHandler))
		r.Method(http.MethodPost, "/", Handler(api.CreateLabTestHandler))
		r.Method(http.MethodPut, "/{labTestID}", Handler(api.UpdateLabTestHandler))
		r.Method(http.MethodDelete, "/{labTestID}", Handler(api.DeleteLabTestHandler))
	})

	// Admin hospital routes
	mux.Route("/hospitals", func(r chi.Router) {
		r.Use(api.RequireLogin)
		r.Use(api.RequireAdmin)
		r.Method(http.MethodGet, "/", Handler(api.GetHospitals))
		r.Method(http.MethodPost, "/", Handler(api.CreateHospital))
		r.Method(http.MethodDelete, "/{hospitalID}", Handler(api.DeleteHospital))
		r.Method(http.MethodGet, "/{hospitalID}/lab-tests", Handler(api.GetHospitalLabTests))
		r.Method(http.MethodPost, "/{hospitalID}/lab-tests", Handler(api.CreateHospitalLabTest))
	})

	return mux
}