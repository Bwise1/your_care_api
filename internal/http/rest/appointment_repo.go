package rest

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/jmoiron/sqlx"
)

func (api *API) CreateLabTestAppointment(ctx context.Context, appointment model.LabAppointmentReq) (int, error) {
	var appointmentID int

	err := api.Deps.DB.RunInTx(ctx, func(tx *sqlx.Tx) error {
		appointmentStmt := `
			INSERT INTO appointments (
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
			"pending",
		)
		if err != nil {
			return err
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		appointmentID = int(lastInsertID)

		detailsStmt := `
			INSERT INTO lab_test_appointment_details (
				appointment_id,
				pickup_type,
				home_location,
				test_type_id,
				hospital_id,
				additional_instructions
			) VALUES (?, ?, ?, ?, ?, ?)`

		var homeLocation sql.NullString
		var hospitalID sql.NullInt64

		if appointment.PickupType == "home" {
			if appointment.HomeLocation == nil {
				return errors.New("home location is required for home pickup type")
			}
			homeLocation = sql.NullString{String: *appointment.HomeLocation, Valid: true}
		} else if appointment.PickupType == "hospital" {
			if appointment.HospitalID == nil {
				return errors.New("hospital ID is required for hospital pickup type")
			}
			hospitalID = sql.NullInt64{Int64: int64(*appointment.HospitalID), Valid: true}
		} else {
			return errors.New("invalid pickup type")
		}

		_, err = tx.ExecContext(ctx, detailsStmt,
			appointmentID,
			appointment.PickupType,
			homeLocation,
			appointment.TestTypeID,
			hospitalID,
			sql.NullString{String: *appointment.AdditionalInstructions, Valid: appointment.AdditionalInstructions != nil},
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
