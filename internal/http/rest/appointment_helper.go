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