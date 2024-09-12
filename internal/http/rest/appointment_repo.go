package rest

import (
	"context"
	"log"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/jmoiron/sqlx"
)

func (api *API) CreateLabTestAppointment(ctx context.Context, appointment model.Appointment, details model.LabTestAppointmentDetails) (int, error) {
	var appointmentID int

	err := api.Deps.DB.RunInTx(ctx, func(tx *sqlx.Tx) error {

		appointmentStmt := `
			INSERT INTO appointment (
				user_id,
				lab_test_id,
				appointment_type,
				appointment_date,
				appointment_time,
				status
			) VALUES (?, ?, ?, ?, ?, ?)`

		result, err := tx.ExecContext(ctx, appointmentStmt,
			appointment.UserID,
			appointment.LabTestID,
			"lab_test",
			appointment.AppointmentDate,
			appointment.AppointmentTime,
			appointment.Status,
		)
		if err != nil {
			return err
		}

		// Get the ID of the inserted appointment
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		appointmentID = int(lastInsertID)

		// Insert into the lab_test_appointment_details table
		detailsStmt := `
			INSERT INTO lab_test_appointment_details (
				appointment_id,
				pickup_type,
				home_location,
				test_type
			) VALUES (?, ?, ?, ?)`

		_, err = tx.ExecContext(ctx, detailsStmt,
			appointmentID,
			details.PickupType,
			details.HomeLocation,
			details.TestType,
		)
		return err
	})

	if err != nil {
		log.Println("error creating lab test appointment", err)
		return 0, err
	}
	return appointmentID, nil
}

func (api *API) CreateDoctorAppointment(ctx context.Context, appointment model.Appointment, details model.DoctorAppointmentDetails) (int, error) {
	var appointmentID int

	err := api.Deps.DB.RunInTx(ctx, func(tx *sqlx.Tx) error {
		// Insert into the appointment table
		appointmentStmt := `
			INSERT INTO appointment (
				user_id,
				doctor_id,
				appointment_type,
				appointment_date,
				appointment_time,
				status
			) VALUES (?, ?, ?, ?, ?, ?)`

		result, err := tx.ExecContext(ctx, appointmentStmt,
			appointment.UserID,
			appointment.DoctorID,
			"doctor", // Appointment type for doctor appointments
			appointment.AppointmentDate,
			appointment.AppointmentTime,
			appointment.Status,
		)
		if err != nil {
			return err
		}

		// Get the ID of the inserted appointment
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		appointmentID = int(lastInsertID)

		// Insert into the doctor_appointment_details table
		detailsStmt := `
			INSERT INTO doctor_appointment_details (
				appointment_id,
				reason_for_visit,
				symptoms,
				additional_notes
			) VALUES (?, ?, ?, ?)`

		_, err = tx.ExecContext(ctx, detailsStmt,
			appointmentID,
			details.ReasonForVisit,
			details.Symptoms,
			details.AdditionalNotes,
		)
		return err
	})

	if err != nil {
		log.Println("error creating doctor appointment", err)
		return 0, err
	}

	return appointmentID, nil
}
