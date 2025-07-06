package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/bwise1/your_care_api/util/values"
)

func (api *API) CreateDoctorAppointment_H(w http.ResponseWriter, r *http.Request) {
	var appointment model.Appointment
	var details model.DoctorAppointmentDetails

	// Decode the request body into the appointment and details structs
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, fmt.Sprintf("%s [CrDoAp]", values.BadRequestBody), http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		http.Error(w, fmt.Sprintf("%s [CrDoAp]", values.BadRequestBody), http.StatusBadRequest)
		return
	}

	// Set context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create the doctor's appointment
	appointmentID, err := api.CreateDoctorAppointment(ctx, appointment, details)
	if err != nil {
		log.Println("error creating doctor appointment", err)
		http.Error(w, fmt.Sprintf("%s [CrDoAp]", values.SystemErr), http.StatusInternalServerError)
		return
	}

	// Respond with the created appointment ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  values.Success,
		"message": "Doctor appointment created successfully",
		"data":    appointmentID,
	})
}

func (api *API) CreateLabTestAppointmentH(appointment model.LabAppointmentReq) (model.Appointment, string, string, error) {

	// Set context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create the lab test appointment
	appointmentID, err := api.CreateLabTestAppointment(ctx, appointment)
	if err != nil {
		return model.Appointment{}, values.Error, fmt.Sprintf("%s [CrLaAp]", values.SystemErr), err
	}

	newAppointment := model.Appointment{
		ID:     appointmentID,
		UserID: appointment.UserID,
		// DoctorID:  appointment.DoctorID,
		// LabTestID: appointment.LabTestID,
	}

	return newAppointment, values.Created, "Lab test appointment created successfully", nil
}

func (api *API) CreateLabTestAppointmentHelper(appointment model.AppointmentDetails, labAppt model.LabTestAppointment) (model.AppointmentDetails, string, string, error) {

	// Set context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Create the lab test appointment
	appointmentID, err := api.CreateLabTestAppRepo(ctx, appointment, labAppt)
	log.Println("appointmentID", appointmentID)
	if err != nil {
		return model.AppointmentDetails{}, values.Error, fmt.Sprintf("%s [CrLaAp]", values.SystemErr), err
	}
	newAppointment := model.AppointmentDetails{
		ID:                  appointmentID,
		UserID:              appointment.UserID,
		AppointmentType:     appointment.AppointmentType,
		AppointmentDatetime: appointment.AppointmentDatetime,
		Status:              "pending",
		LabTestDetails: &model.LabTestAppointment{
			ID:                     labAppt.ID,
			PickupType:             labAppt.PickupType,
			HomeLocation:           labAppt.HomeLocation,
			TestTypeID:             labAppt.TestTypeID,
			HospitalID:             labAppt.HospitalID,
			AdditionalInstructions: labAppt.AdditionalInstructions,
		},
	}
	return newAppointment, values.Created, "Lab test appointment created successfully", nil
}

func (api *API) FetchAllAppointments(filter model.AppointmentFilter) ([]model.AppointmentDetails, string, string, error) {
	// Set context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("userID", filter.UserID)
	// Fetch all appointments
	// appointments, err := api.FetchAllAppointmentsRepo(ctx, &filter.UserID)
	appointments, err := api.FetchFilteredAppointmentsRepo(ctx, filter)

	if err != nil {
		return nil, values.Error, fmt.Sprintf("%s [FtAlAp]", values.SystemErr), err
	}

	return appointments, values.Success, "Appointments fetched successfully", nil

}

// Admin Helper Functions

func (api *API) AdminFetchAllAppointmentsHelper(filter model.AdminAppointmentFilter) ([]model.AppointmentDetails, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	appointments, err := api.AdminFetchAllAppointmentsRepo(ctx, filter)
	if err != nil {
		return nil, values.Error, fmt.Sprintf("%s [AdFtAlAp]", values.SystemErr), err
	}

	return appointments, values.Success, "Appointments fetched successfully", nil
}

func (api *API) AdminGetAppointmentDetailsHelper(appointmentID int) (model.AppointmentDetails, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the same simple query as the working list function
	appointment, err := api.AdminGetAppointmentDetailsRepo(ctx, appointmentID)
	if err != nil {
		return model.AppointmentDetails{}, values.Error, fmt.Sprintf("%s [AdGtApDt]", values.SystemErr), err
	}

	return appointment, values.Success, "Appointment details retrieved successfully", nil
}

func (api *API) AdminConfirmAppointmentHelper(appointmentID, adminID int, req model.AdminAppointmentAction) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := api.UpdateAppointmentStatus(ctx, appointmentID, string(model.StatusConfirmed), req.Notes, &adminID)
	if err != nil {
		return values.Error, fmt.Sprintf("%s [AdCfAp]", values.SystemErr), err
	}

	return values.Success, "Appointment confirmed successfully", nil
}

func (api *API) AdminRejectAppointmentHelper(appointmentID, adminID int, req model.AdminAppointmentAction) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := api.RejectAppointment(ctx, appointmentID, req.RejectionReason, req.Notes, &adminID)
	if err != nil {
		return values.Error, fmt.Sprintf("%s [AdRjAp]", values.SystemErr), err
	}

	return values.Success, "Appointment rejected", nil
}

func (api *API) AdminRescheduleAppointmentHelper(appointmentID, adminID int, req model.AdminAppointmentAction) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if req.ProposedDate == nil || req.ProposedTime == nil {
		return values.BadRequestBody, "Proposed date and time are required", fmt.Errorf("missing proposed date/time")
	}

	err := api.CreateRescheduleOffer(ctx, appointmentID, *req.ProposedDate, *req.ProposedTime, req.Notes, &adminID)
	if err != nil {
		return values.Error, fmt.Sprintf("%s [AdRsAp]", values.SystemErr), err
	}

	return values.Success, "Reschedule offer created", nil
}

func (api *API) AdminCancelAppointmentHelper(appointmentID, adminID int, req model.AdminAppointmentAction) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := api.UpdateAppointmentStatus(ctx, appointmentID, string(model.StatusCanceled), req.Notes, &adminID)
	if err != nil {
		return values.Error, fmt.Sprintf("%s [AdCnAp]", values.SystemErr), err
	}

	return values.Success, "Appointment canceled", nil
}

func (api *API) AdminUpdateNotesHelper(appointmentID int, notes string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := api.UpdateAppointmentAdminNotes(ctx, appointmentID, notes)
	if err != nil {
		return values.Error, fmt.Sprintf("%s [AdUpNt]", values.SystemErr), err
	}

	return values.Success, "Notes updated successfully", nil
}

// User Helper Functions

func (api *API) GetAppointmentDetailsHelper(appointmentID, userID int) (model.UserAppointmentResponse, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the same simple query as the working list function
	appointment, err := api.GetDetailedAppointmentRepo(ctx, appointmentID, &userID)
	if err != nil {
		return model.UserAppointmentResponse{}, values.Error, fmt.Sprintf("%s [GtApDt]", values.SystemErr), err
	}

	statusHistory, _ := api.GetAppointmentStatusHistory(appointmentID)
	nextActions := api.getNextActionsForUser(model.AppointmentStatus(appointment.Status))

	response := model.UserAppointmentResponse{
		//Appointment:   appointment,
		StatusHistory: statusHistory,
		NextActions:   nextActions,
	}

	return response, values.Success, "Appointment details retrieved successfully", nil
}

func (api *API) AcceptRescheduleOfferHelper(appointmentID, userID, offerID int) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := api.AcceptRescheduleOfferRepo(ctx, appointmentID, userID, offerID)
	if err != nil {
		return values.Error, fmt.Sprintf("%s [AcRsOf]", values.SystemErr), err
	}

	return values.Success, "Reschedule offer accepted", nil
}

func (api *API) RejectRescheduleOfferHelper(appointmentID, userID, offerID int, reason *string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := api.RejectRescheduleOfferRepo(ctx, appointmentID, userID, offerID, reason)
	if err != nil {
		return values.Error, fmt.Sprintf("%s [RjRsOf]", values.SystemErr), err
	}

	return values.Success, "Reschedule offer rejected", nil
}

func (api *API) CancelAppointmentHelper(appointmentID, userID int) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := api.UpdateAppointmentStatus(ctx, appointmentID, string(model.StatusCanceled), nil, &userID)
	if err != nil {
		return values.Error, fmt.Sprintf("%s [CnAp]", values.SystemErr), err
	}

	return values.Success, "Appointment canceled successfully", nil
}

func (api *API) GetAppointmentStatusHistory(appointmentID int) ([]model.AppointmentStatusLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return api.GetAppointmentStatusHistoryRepo(ctx, appointmentID)
}

func (api *API) AdminGetAppointmentHistoryHelper(appointmentID int) ([]model.AppointmentStatusLog, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	history, err := api.GetAppointmentStatusHistoryRepo(ctx, appointmentID)
	if err != nil {
		return nil, values.Error, fmt.Sprintf("%s [AdGtHist]", values.SystemErr), err
	}

	return history, values.Success, "Appointment history retrieved successfully", nil
}

func (api *API) GetAppointmentHistoryHelper(appointmentID, userID int) ([]model.AppointmentStatusLog, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First verify that the appointment belongs to the user
	_, err := api.GetAppointmentDetailsRepo(ctx, appointmentID, userID)
	if err != nil {
		return nil, values.NotFound, "Appointment not found or access denied", err
	}

	history, err := api.GetAppointmentStatusHistoryRepo(ctx, appointmentID)
	if err != nil {
		return nil, values.Error, fmt.Sprintf("%s [GtHist]", values.SystemErr), err
	}

	return history, values.Success, "Appointment history retrieved successfully", nil
}
