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

func (api *API) AppointmentRoutes() chi.Router {
	mux := chi.NewRouter()
	// mux.Method(http.MethodPost, "/login", Handler(api.Login))
	mux.Method(http.MethodPost, "/lab-test", Handler(api.LabAppointment))
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
