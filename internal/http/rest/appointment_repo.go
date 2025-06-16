package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

func (api *API) CreateLabTestAppRepo(ctx context.Context, appointment model.AppointmentDetails, labApt model.LabTestAppointment) (int, error) {
	var appointmentID int
	err := api.Deps.DB.RunInTx(ctx, func(tx *sqlx.Tx) error {
		appointmentStmt := `
			INSERT INTO appointments (
				user_id,
				appointment_type,
				appointment_datetime
			) VALUES (?, ?, ?)`

		result, err := tx.ExecContext(ctx, appointmentStmt,
			appointment.UserID,
			"lab_test",
			appointment.AppointmentDatetime,
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
			INSERT INTO lab_test_appointments (
				appointment_id,
				test_type_id,
				pickup_type,
				home_location,
				hospital_id,
				additional_instructions
			) VALUES (?, ?, ?, ?, ?, ?)`

		var homeLocation sql.NullString
		var hospitalID sql.NullInt64

		if labApt.PickupType == "home" {
			if labApt.HomeLocation == nil {
				return errors.New("home location is required for home pickup type")
			}
			homeLocation = sql.NullString{String: *labApt.HomeLocation, Valid: true}
		} else if labApt.PickupType == "hospital" {
			if labApt.HospitalID == nil {
				return errors.New("hospital ID is required for hospital pickup type")
			}
			hospitalID = sql.NullInt64{Int64: int64(*labApt.HospitalID), Valid: true}
		} else {
			return errors.New("invalid pickup type")
		}

		_, err = tx.ExecContext(ctx, detailsStmt,
			appointmentID,
			labApt.TestTypeID,
			labApt.PickupType,
			homeLocation,
			hospitalID,
			sql.NullString{String: *labApt.AdditionalInstructions, Valid: labApt.AdditionalInstructions != nil},
		)
		return err
	})

	if err != nil {
		log.Println(appointment)
		log.Println("error creating lab test appointment", err)
		return 0, err
	}
	log.Println("appointmentID from repo", appointmentID)
	return appointmentID, nil

}

func (api *API) FetchAllAppointmentsRepo(ctx context.Context, userID *int) ([]model.AppointmentDetails, error) {

	query := `
		SELECT
			a.id,
			a.user_id,
			a.appointment_type,
			a.appointment_datetime,
			a.status,
			a.created_at,
			a.updated_at,
			CASE
				WHEN a.appointment_type = 'doctor' THEN
					JSON_OBJECT(
						'id', da.id,
						'doctor_id', da.doctor_id,
						'reason_for_visit', da.reason_for_visit,
						'symptoms', da.symptoms,
						'additional_notes', da.additional_notes
					)
			END as doctor_details,
			CASE
				WHEN a.appointment_type = 'lab_test' THEN
					JSON_OBJECT(
						'id', la.id,
						'test_type_id', la.test_type_id,
						'pickup_type', la.pickup_type,
						'home_location', la.home_location,
						'hospital_id', la.hospital_id,
						'additional_instructions', la.additional_instructions
					)
			END as lab_test_details
		FROM
			appointments a
			LEFT JOIN doctor_appointments da ON a.id = da.appointment_id AND a.appointment_type = 'doctor'
			LEFT JOIN lab_test_appointments la ON a.id = la.appointment_id AND a.appointment_type = 'lab_test'
		WHERE
			(? IS NULL OR a.user_id = ?)
		ORDER BY
			a.appointment_datetime DESC`

	var rows []model.AppointmentRow
	err := api.Deps.DB.SelectContext(ctx, &rows, query, userID, userID)
	if err != nil {
		log.Println("error fetching appointments", err)
		return nil, fmt.Errorf("failed to fetch appointments: %w", err)
	}

	appointments := make([]model.AppointmentDetails, len(rows))
	for i, row := range rows {
		appointments[i] = model.AppointmentDetails{
			ID:                  row.ID,
			UserID:              row.UserID,
			AppointmentType:     row.AppointmentType,
			AppointmentDatetime: row.AppointmentDatetime,
			Status:              row.Status,
			CreatedAt:           row.CreatedAt,
			UpdatedAt:           row.UpdatedAt,
		}

		if row.AppointmentType == "doctor" && len(row.DoctorDetailsJSON) > 0 {
			var doctorDetails model.DoctorAppointment
			if err := json.Unmarshal(row.DoctorDetailsJSON, &doctorDetails); err != nil {
				return nil, fmt.Errorf("failed to unmarshal doctor details: %w", err)
			}
			appointments[i].DoctorDetails = &doctorDetails
		}

		if row.AppointmentType == "lab_test" && len(row.LabTestDetailsJSON) > 0 {
			var labTestDetails model.LabTestAppointment
			if err := json.Unmarshal(row.LabTestDetailsJSON, &labTestDetails); err != nil {
				return nil, fmt.Errorf("failed to unmarshal lab test details: %w", err)
			}
			appointments[i].LabTestDetails = &labTestDetails
		}
	}
	return appointments, nil
}

func (api *API) FetchFilteredAppointmentsRepo(ctx context.Context, filter model.AppointmentFilter) ([]model.AppointmentDetails, error) {

	query := `
		SELECT
			a.id,
			a.user_id,
			a.appointment_type,
			a.appointment_datetime,
			a.status,
			a.created_at,
			a.updated_at,
			CASE
				WHEN a.appointment_type = 'doctor' THEN
					JSON_OBJECT(
						'id', da.id,
						'doctor_id', da.doctor_id,
						'reason_for_visit', da.reason_for_visit,
						'symptoms', da.symptoms,
						'additional_notes', da.additional_notes
					)
			END as doctor_details,
			CASE
				WHEN a.appointment_type = 'lab_test' THEN
					JSON_OBJECT(
						'id', la.id,
						'test_type_id', la.test_type_id,
						'pickup_type', la.pickup_type,
						'home_location', la.home_location,
						'hospital_id', la.hospital_id,
						'additional_instructions', la.additional_instructions
					)
			END as lab_test_details
		FROM
			appointments a
			LEFT JOIN doctor_appointments da ON a.id = da.appointment_id AND a.appointment_type = 'doctor'
			LEFT JOIN lab_test_appointments la ON a.id = la.appointment_id AND a.appointment_type = 'lab_test'
		 WHERE a.user_id = ?`

	var args []interface{}
	args = append(args, filter.UserID)

	if filter.Date != "" {
		query += " AND DATE(a.appointment_datetime) = ?"
		args = append(args, filter.Date)
	}

	if filter.Upcoming {
		query += " AND a.appointment_datetime > NOW()"
		query += " ORDER BY a.appointment_datetime ASC"
	} else if filter.History {
		query += " AND a.appointment_datetime < NOW()"
		query += " ORDER BY a.appointment_datetime DESC"
	} else {
		query += " ORDER BY a.appointment_datetime DESC"
	}

	log.Println("query", query)

	log.Println("Arguments:")
	for i, arg := range args {
		log.Printf("  Arg[%d]: %v (type: %T)", i, arg, arg)
	}
	var rows []model.AppointmentRow
	err := api.Deps.DB.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		log.Println("error fetching appointments", err)
		return nil, fmt.Errorf("failed to fetch appointments: %w", err)
	}

	appointments := make([]model.AppointmentDetails, len(rows))
	for i, row := range rows {
		appointments[i] = model.AppointmentDetails{
			ID:                  row.ID,
			UserID:              row.UserID,
			AppointmentType:     row.AppointmentType,
			AppointmentDatetime: row.AppointmentDatetime,
			Status:              row.Status,
			CreatedAt:           row.CreatedAt,
			UpdatedAt:           row.UpdatedAt,
		}

		if row.AppointmentType == "doctor" && len(row.DoctorDetailsJSON) > 0 {
			var doctorDetails model.DoctorAppointment
			if err := json.Unmarshal(row.DoctorDetailsJSON, &doctorDetails); err != nil {
				return nil, fmt.Errorf("failed to unmarshal doctor details: %w", err)
			}
			appointments[i].DoctorDetails = &doctorDetails
		}

		if row.AppointmentType == "lab_test" && len(row.LabTestDetailsJSON) > 0 {
			var labTestDetails model.LabTestAppointment
			if err := json.Unmarshal(row.LabTestDetailsJSON, &labTestDetails); err != nil {
				return nil, fmt.Errorf("failed to unmarshal lab test details: %w", err)
			}
			appointments[i].LabTestDetails = &labTestDetails
		}
	}
	return appointments, nil
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

func (api *API) GetLabTestAppointments(ctx context.Context, userID int) ([]model.Appointment, error) {
	var appointments []model.Appointment

	query := `
		SELECT
			a.id,
			a.appointment_date,
			a.appointment_time,
			a.status,
			lta.pickup_type,
			lta.home_location,
			lta.test_type_id,
			lta.hospital_id,
			lta.additional_instructions
		FROM appointments a
		JOIN lab_test_appointment_details lta ON a.id = lta.appointment_id
		WHERE a.user_id = ?`

	err := api.Deps.DB.SelectContext(ctx, &appointments, query, userID)
	if err != nil {
		log.Println("error getting lab test appointments", err)
		return nil, err
	}

	return appointments, nil
}
