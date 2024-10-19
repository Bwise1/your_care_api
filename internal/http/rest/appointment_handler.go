package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/bwise1/your_care_api/util"
	"github.com/bwise1/your_care_api/util/tracing"
	"github.com/bwise1/your_care_api/util/values"
	"github.com/go-chi/chi/v5"
)

func (api *API) AppointmentRoutes() chi.Router {
	mux := chi.NewRouter()

	// mux.Method(http.MethodPost, "/login", Handler(api.Login))
	mux.Method(http.MethodPost, "/lab-test", Handler(api.LabAppointment))
	mux.Method(http.MethodPost, "/lab-test-appointment", Handler(api.CreateLabTestAppointmentHandler))
	mux.Route("/", func(r chi.Router) {
		r.Use(api.RequireLogin)
		r.Method(http.MethodGet, "/", Handler(api.FetchAllAppointmentsHandler))
	})
	return mux
}

func (api *API) LabAppointment(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	var appointment model.LabAppointmentReq
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &appointment); decodeErr != nil {
		log.Println(decodeErr)
		return respondWithError(decodeErr, "unable to parse lab appointment request", values.BadRequestBody, &tc)
	}
	log.Println("appointment", appointment)
	// var details model.LabTestAppointmentDetails
	// if decodeErr := util.DecodeJSONBody(&tc, r.Body, &details); decodeErr != nil {
	// 	log.Println(decodeErr)
	// 	return respondWithError(decodeErr, "unable to parse lab appointment details request", values.BadRequestBody, &tc)
	// }

	newAppointment, status, message, err := api.CreateLabTestAppointmentH(appointment)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       newAppointment,
	}
}

func (api *API) CreateLabTestAppointmentHandler(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	var req model.CreateLabTestAppointmentRequest
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		log.Println(decodeErr)
		return respondWithError(decodeErr, "unable to parse lab appointment request", values.BadRequestBody, &tc)
	}
	appointmentDatetime, err := time.Parse("2006-01-02 15:04:05", req.AppointmentDate)

	log.Println("appointmentDate", appointmentDatetime)
	if err != nil {
		return respondWithError(err, "invalid date format", values.BadRequestBody, &tc)
	}

	appointment := &model.AppointmentDetails{
		UserID:              req.UserID,
		AppointmentType:     "lab_test",
		AppointmentDatetime: &appointmentDatetime,
		Status:              "pending",
	}

	labAppt := &model.LabTestAppointment{
		TestTypeID:             req.TestTypeID,
		PickupType:             req.PickupType,
		HomeLocation:           req.HomeLocation,
		AdditionalInstructions: req.AdditionalInstructions,
		HospitalID:             req.HospitalID,
	}

	newAppointment, status, message, err := api.CreateLabTestAppointmentHelper(*appointment, *labAppt)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       newAppointment,
	}
}

func (api *API) FetchAllAppointmentsHandler(_ http.ResponseWriter, r *http.Request) *ServerResponse {

	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	userID := r.Context().Value("user_id")
	log.Println("userID", userID)
	appointments, status, message, err := api.FetchAllAppointments(userID.(int))
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       appointments,
	}

}
