package rest

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/bwise1/your_care_api/util"
	"github.com/bwise1/your_care_api/util/tracing"
	"github.com/bwise1/your_care_api/util/values"
	"github.com/go-chi/chi/v5"
)

func (api *API) AppointmentRoutes() chi.Router {
	mux := chi.NewRouter()

	// User routes
	mux.Route("/", func(r chi.Router) {
		r.Use(api.RequireLogin)
		r.Method(http.MethodGet, "/", Handler(api.FetchAllAppointmentsHandler))
		r.Method(http.MethodPost, "/lab-test-appointment", Handler(api.CreateLabTestAppointmentHandler))
		r.Method(http.MethodPost, "/lab-test", Handler(api.LabAppointment))
		r.Method(http.MethodGet, "/status-stages", Handler(api.GetAppointmentStatusStages))
		r.Method(http.MethodGet, "/{id}", Handler(api.GetAppointmentDetails))
		r.Method(http.MethodGet, "/{id}/history", Handler(api.GetAppointmentHistory))
		r.Method(http.MethodPut, "/{id}/reschedule/accept", Handler(api.AcceptRescheduleOffer))
		r.Method(http.MethodPut, "/{id}/reschedule/reject", Handler(api.RejectRescheduleOffer))
		r.Method(http.MethodDelete, "/{id}", Handler(api.CancelAppointment))
	})
	return mux
}


func (api *API) LabAppointment(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	userID := r.Context().Value("user_id")

	var appointment model.LabAppointmentReq
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &appointment); decodeErr != nil {
		log.Println(decodeErr)
		return respondWithError(decodeErr, "unable to parse lab appointment request", values.BadRequestBody, &tc)
	}

	appointment.UserID = userID.(int)
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

	userID := r.Context().Value("user_id")

	var req model.CreateLabTestAppointmentRequest
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		log.Println(decodeErr)
		return respondWithError(decodeErr, "unable to parse lab appointment request", values.BadRequestBody, &tc)
	}
	log.Println("req", req)
	appointmentDatetime, err := time.Parse("2006-01-02 15:04:05", req.AppointmentDate)

	log.Println("appointmentDate", appointmentDatetime)
	if err != nil {
		return respondWithError(err, "invalid date format", values.BadRequestBody, &tc)
	}

	//add user id to the req object
	req.UserID = userID.(int)
	// log.Println("userID", req.UserID)
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

	isAdmin, _ := r.Context().Value("is_admin").(bool)

	log.Println("admin? ", isAdmin)

	queryParams := r.URL.Query()
	filter := model.AppointmentFilter{
		Date:     queryParams.Get("date"),
		Upcoming: queryParams.Get("upcoming") == "true",
		History:  queryParams.Get("history") == "true",
	}

	// var userIDPtr *int
	if !isAdmin {
		userIDVal := r.Context().Value("user_id")
		userID, ok := userIDVal.(int)
		if !ok {
			return respondWithError(nil, "user_id not found in context", values.NotAuthorised, &tc)
		}
		filter.UserID = userID
		// userIDPtr = &userID
	}

	// log.Println("userID", *userIDPtr)

	// if !isAdmin && userIDPtr != nil {

	// 	filter.UserID = *userIDPtr
	// }

	log.Println("Filters", filter.UserID)

	appointments, status, message, err := api.FetchAllAppointments(filter)
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

// Admin Appointment Handlers

func (api *API) AdminFetchAllAppointments(_ http.ResponseWriter, r *http.Request) *ServerResponse {

	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	queryParams := r.URL.Query()
	filter := model.AdminAppointmentFilter{
		Page:  1,
		Limit: 50,
	}

	// Parse query parameters
	if page := queryParams.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filter.Page = p
		}
	}
	if limit := queryParams.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}
	if dateFrom := queryParams.Get("date_from"); dateFrom != "" {
		filter.DateFrom = &dateFrom
	}
	if dateTo := queryParams.Get("date_to"); dateTo != "" {
		filter.DateTo = &dateTo
	}
	if providerID := queryParams.Get("provider_id"); providerID != "" {
		if pid, err := strconv.Atoi(providerID); err == nil {
			filter.ProviderID = &pid
		}
	}

	// Parse status and appointment_type arrays
	if statuses := queryParams["status"]; len(statuses) > 0 {
		filter.Status = statuses
	}
	if appointmentTypes := queryParams["appointment_type"]; len(appointmentTypes) > 0 {
		filter.AppointmentType = appointmentTypes
	}

	appointments, status, message, err := api.AdminFetchAllAppointmentsHelper(filter)
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

func (api *API) AdminGetAppointmentDetails(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	appointment, status, message, err := api.AdminGetAppointmentDetailsHelper(appointmentID)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       appointment,
	}
}


func (api *API) AdminConfirmAppointment(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	var req model.AdminAppointmentAction
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse request", values.BadRequestBody, &tc)
	}

	adminID := r.Context().Value("user_id").(int)

	status, message, err := api.AdminConfirmAppointmentHelper(appointmentID, adminID, req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       nil,
	}
}

func (api *API) AdminRejectAppointment(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	var req model.AdminAppointmentAction
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse request", values.BadRequestBody, &tc)
	}

	adminID := r.Context().Value("user_id").(int)

	status, message, err := api.AdminRejectAppointmentHelper(appointmentID, adminID, req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       nil,
	}
}

func (api *API) AdminRescheduleAppointment(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	var req model.AdminAppointmentAction
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse request", values.BadRequestBody, &tc)
	}

	adminID := r.Context().Value("user_id").(int)

	status, message, err := api.AdminRescheduleAppointmentHelper(appointmentID, adminID, req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       nil,
	}
}

func (api *API) AdminCancelAppointment(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	var req model.AdminAppointmentAction
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse request", values.BadRequestBody, &tc)
	}

	adminID := r.Context().Value("user_id").(int)

	status, message, err := api.AdminCancelAppointmentHelper(appointmentID, adminID, req)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       nil,
	}
}

func (api *API) AdminUpdateNotes(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse request", values.BadRequestBody, &tc)
	}

	status, message, err := api.AdminUpdateNotesHelper(appointmentID, req.Notes)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       nil,
	}
}

func (api *API) AdminGetAppointmentHistory(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	history, status, message, err := api.AdminGetAppointmentHistoryHelper(appointmentID)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       history,
	}
}

// User Appointment Handlers

func (api *API) GetAppointmentStatusStages(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	statuses := []string{
		string(model.StatusPending),
		string(model.StatusConfirmed),
		string(model.StatusScheduled),
		string(model.StatusRescheduleOffered),
		string(model.StatusRescheduleAccepted),
		string(model.StatusInProgress),
		string(model.StatusCompleted),
		string(model.StatusCanceled),
		string(model.StatusRejected),
		string(model.StatusNoShow),
	}

	return &ServerResponse{
		Message:    "Status stages retrieved successfully",
		Status:     values.Success,
		StatusCode: util.StatusCode(values.Success),
		Data:       map[string][]string{"statuses": statuses},
	}
}

func (api *API) GetAppointmentDetails(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	userID := r.Context().Value("user_id").(int)

	appointment, status, message, err := api.GetAppointmentDetailsHelper(appointmentID, userID)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       appointment,
	}
}

func (api *API) AcceptRescheduleOffer(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	var req model.RescheduleResponse
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse request", values.BadRequestBody, &tc)
	}

	userID := r.Context().Value("user_id").(int)

	status, message, err := api.AcceptRescheduleOfferHelper(appointmentID, userID, req.OfferID)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       nil,
	}
}

func (api *API) RejectRescheduleOffer(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	var req model.RescheduleResponse
	if decodeErr := util.DecodeJSONBody(&tc, r.Body, &req); decodeErr != nil {
		return respondWithError(decodeErr, "unable to parse request", values.BadRequestBody, &tc)
	}

	userID := r.Context().Value("user_id").(int)

	status, message, err := api.RejectRescheduleOfferHelper(appointmentID, userID, req.OfferID, req.Reason)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       nil,
	}
}

func (api *API) CancelAppointment(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	userID := r.Context().Value("user_id").(int)

	status, message, err := api.CancelAppointmentHelper(appointmentID, userID)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       nil,
	}
}

func (api *API) GetAppointmentHistory(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	tc := r.Context().Value(values.ContextTracingKey).(tracing.Context)

	appointmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return respondWithError(err, "Invalid appointment ID", values.BadRequestBody, &tc)
	}

	userID := r.Context().Value("user_id").(int)

	history, status, message, err := api.GetAppointmentHistoryHelper(appointmentID, userID)
	if err != nil {
		return respondWithError(err, message, status, &tc)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Data:       history,
	}
}

// Helper function to determine next actions for user
func (api *API) getNextActionsForUser(status model.AppointmentStatus) []string {
	switch status {
	case model.StatusPending:
		return []string{"cancel"}
	case model.StatusConfirmed:
		return []string{"cancel"}
	case model.StatusScheduled:
		return []string{"cancel"}
	case model.StatusRescheduleOffered:
		return []string{"accept_reschedule", "reject_reschedule", "cancel"}
	case model.StatusRescheduleAccepted:
		return []string{"cancel"}
	case model.StatusInProgress:
		return []string{}
	case model.StatusCompleted:
		return []string{}
	case model.StatusCanceled:
		return []string{}
	case model.StatusRejected:
		return []string{}
	case model.StatusNoShow:
		return []string{}
	default:
		return []string{}
	}
}

// Helper function to determine next actions for admin
func (api *API) getNextActionsForAdmin(status model.AppointmentStatus) []string {
	switch status {
	case model.StatusPending:
		return []string{"confirm", "reject", "reschedule", "cancel"}
	case model.StatusConfirmed:
		return []string{"reschedule", "cancel", "mark_in_progress"}
	case model.StatusScheduled:
		return []string{"reschedule", "cancel", "mark_in_progress"}
	case model.StatusRescheduleOffered:
		return []string{"cancel", "offer_new_reschedule"}
	case model.StatusRescheduleAccepted:
		return []string{"confirm", "cancel", "mark_in_progress"}
	case model.StatusInProgress:
		return []string{"mark_completed", "mark_no_show"}
	case model.StatusCompleted:
		return []string{}
	case model.StatusCanceled:
		return []string{}
	case model.StatusRejected:
		return []string{}
	case model.StatusNoShow:
		return []string{}
	default:
		return []string{}
	}
}
