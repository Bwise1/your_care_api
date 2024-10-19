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
	appointmentID, err := api.CreateLabTestRepo(ctx, appointment, labAppt)
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

func (api *API) FetchAllAppointments(userID int) ([]model.AppointmentDetails, string, string, error) {
	// Set context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("userID", userID)
	// Fetch all appointments
	appointments, err := api.FetchAllAppointmentsRepo(ctx, &userID)
	if err != nil {
		return nil, values.Error, fmt.Sprintf("%s [FtAlAp]", values.SystemErr), err
	}

	return appointments, values.Success, "Appointments fetched successfully", nil

}
