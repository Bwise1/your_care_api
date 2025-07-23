package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

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
		if err != nil {
			return err
		}

		// Create initial status history entry
		historyStmt := `
			INSERT INTO appointment_status_history (appointment_id, status, notes, changed_by_user_id)
			VALUES (?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, historyStmt, appointmentID, "pending", "Appointment created", appointment.UserID)
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
				appointment_datetime,
				status
			) VALUES (?, ?, ?, ?)`

		result, err := tx.ExecContext(ctx, appointmentStmt,
			appointment.UserID,
			"lab_test",
			appointment.AppointmentDatetime,
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

		// Create initial status history entry
		historyStmt := `
			INSERT INTO appointment_status_history (appointment_id, status, notes, changed_by_user_id)
			VALUES (?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, historyStmt, appointmentID, "pending", "Appointment created", appointment.UserID)
		if err != nil {
			return err
		}

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
            a.appointment_date,
            a.appointment_time,
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
        WHERE 1=1`

	var args []interface{}

	if filter.UserID != 0 {
		query += " AND a.user_id = ?"
		args = append(args, filter.UserID)
	}

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

	// log.Println("query", query)

	// log.Println("Arguments:")
	// for i, arg := range args {
	// 	log.Printf("  Arg[%d]: %v (type: %T)", i, arg, arg)
	// }
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

		// Create initial status history entry
		historyStmt := `
			INSERT INTO appointment_status_history (appointment_id, status, notes, changed_by_user_id)
			VALUES (?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, historyStmt, appointmentID, appointment.Status, "Appointment created", appointment.UserID)
		if err != nil {
			return err
		}

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
			a.appointment_datetime,
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

// New Repository Functions for Enhanced Appointment System

func (api *API) AdminFetchAllAppointmentsRepo(ctx context.Context, filter model.AdminAppointmentFilter) ([]model.AppointmentDetails, error) {
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
		WHERE 1=1`

	var args []interface{}

	// Add filters
	if len(filter.Status) > 0 {
		placeholders := "?" + strings.Repeat(",?", len(filter.Status)-1)
		query += " AND a.status IN (" + placeholders + ")"
		for _, status := range filter.Status {
			args = append(args, status)
		}
	}

	if len(filter.AppointmentType) > 0 {
		placeholders := "?" + strings.Repeat(",?", len(filter.AppointmentType)-1)
		query += " AND a.appointment_type IN (" + placeholders + ")"
		for _, appointmentType := range filter.AppointmentType {
			args = append(args, appointmentType)
		}
	}

	if filter.DateFrom != nil {
		query += " AND DATE(a.appointment_datetime) >= ?"
		args = append(args, *filter.DateFrom)
	}

	if filter.DateTo != nil {
		query += " AND DATE(a.appointment_datetime) <= ?"
		args = append(args, *filter.DateTo)
	}

	if filter.ProviderID != nil {
		query += " AND a.provider_id = ?"
		args = append(args, *filter.ProviderID)
	}

	query += " ORDER BY a.created_at DESC"

	// Add pagination
	query += " LIMIT ? OFFSET ?"
	args = append(args, filter.Limit, (filter.Page-1)*filter.Limit)

	var rows []model.AppointmentRow
	err := api.Deps.DB.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		log.Println("error fetching admin appointments", err)
		return nil, fmt.Errorf("failed to fetch admin appointments: %w", err)
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

func (api *API) AdminGetAppointmentDetailsRepo(ctx context.Context, appointmentID int) (model.AppointmentDetails, error) {
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
		WHERE a.id = ?`

	var row model.AppointmentRow
	err := api.Deps.DB.GetContext(ctx, &row, query, appointmentID)
	if err != nil {
		log.Println("error fetching appointment details", err)
		return model.AppointmentDetails{}, fmt.Errorf("failed to fetch appointment details: %w", err)
	}

	
	appointment := model.AppointmentDetails{
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
			return model.AppointmentDetails{}, fmt.Errorf("failed to unmarshal doctor details: %w", err)
		}
		appointment.DoctorDetails = &doctorDetails
	}

	if row.AppointmentType == "lab_test" && len(row.LabTestDetailsJSON) > 0 {
		var labTestDetails model.LabTestAppointment
		if err := json.Unmarshal(row.LabTestDetailsJSON, &labTestDetails); err != nil {
			return model.AppointmentDetails{}, fmt.Errorf("failed to unmarshal lab test details: %w", err)
		}
		appointment.LabTestDetails = &labTestDetails
	}

	return appointment, nil
}

func (api *API) GetAppointmentDetailsRepo(ctx context.Context, appointmentID, userID int) (model.AppointmentDetails, error) {
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
		WHERE a.id = ? AND a.user_id = ?`

	var row model.AppointmentRow
	err := api.Deps.DB.GetContext(ctx, &row, query, appointmentID, userID)
	if err != nil {
		log.Println("error fetching user appointment details", err)
		return model.AppointmentDetails{}, fmt.Errorf("failed to fetch appointment details: %w", err)
	}

	
	appointment := model.AppointmentDetails{
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
			return model.AppointmentDetails{}, fmt.Errorf("failed to unmarshal doctor details: %w", err)
		}
		appointment.DoctorDetails = &doctorDetails
	}

	if row.AppointmentType == "lab_test" && len(row.LabTestDetailsJSON) > 0 {
		var labTestDetails model.LabTestAppointment
		if err := json.Unmarshal(row.LabTestDetailsJSON, &labTestDetails); err != nil {
			return model.AppointmentDetails{}, fmt.Errorf("failed to unmarshal lab test details: %w", err)
		}
		appointment.LabTestDetails = &labTestDetails
	}

	return appointment, nil
}

// Enhanced detailed appointment repository function
func (api *API) GetDetailedAppointmentRepo(ctx context.Context, appointmentID int, userID *int) (model.DetailedAppointment, error) {
	// Main appointment query with all joins
	query := `
		SELECT
			a.id,
			a.user_id,
			a.appointment_type,
			a.appointment_datetime,
			a.status,
			a.created_at,
			a.updated_at,

			-- User details
			u.firstName,
			u.lastName,
			u.email,
			u.dateOfBirth,
			u.sex,

			-- Hospital details (if applicable)
			h.id as hospital_id,
			h.name as hospital_name,
			h.address as hospital_address,
			h.phone as hospital_phone,
			h.email as hospital_email,

			-- Lab test details (if lab appointment)
			la.id as lab_appointment_id,
			la.test_type_id,
			la.pickup_type,
			la.home_location,
			la.additional_instructions,

			-- Test type details (if lab appointment)
			lt.name as test_name,
			lt.description as test_description,

			-- Doctor appointment details (if doctor appointment)
			da.id as doctor_appointment_id,
			da.doctor_id,
			da.reason_for_visit,
			da.symptoms,
			da.additional_notes as doctor_notes,

			-- Doctor details (if doctor appointment)
			d.name as doctor_name,
			d.specialization,
			d.email as doctor_email,
			d.phone as doctor_phone

		FROM appointments a

		-- Join user details
		LEFT JOIN users u ON a.user_id = u.id

		-- Join lab test appointment details
		LEFT JOIN lab_test_appointments la ON a.id = la.appointment_id AND a.appointment_type = 'lab_test'
		LEFT JOIN lab_tests lt ON la.test_type_id = lt.id
		LEFT JOIN hospitals h ON la.hospital_id = h.id

		-- Join doctor appointment details
		LEFT JOIN doctor_appointments da ON a.id = da.appointment_id AND a.appointment_type = 'doctor'
		LEFT JOIN doctors d ON da.doctor_id = d.id
		LEFT JOIN hospitals h2 ON d.hospital_id = h2.id

		WHERE a.id = ?`

	// Add user restriction if not admin
	args := []interface{}{appointmentID}
	if userID != nil {
		query += " AND a.user_id = ?"
		args = append(args, *userID)
	}

	// Struct to scan the complex joined result
	var result struct {
		// Appointment fields
		ID                  int        `db:"id"`
		UserID              int        `db:"user_id"`
		AppointmentType     string     `db:"appointment_type"`
		AppointmentDatetime *time.Time `db:"appointment_datetime"`
		Status              string     `db:"status"`
		CreatedAt           *time.Time `db:"created_at"`
		UpdatedAt           *time.Time `db:"updated_at"`

		// User fields
		FirstName   string `db:"firstName"`
		LastName    string `db:"lastName"`
		Email       string `db:"email"`
		DateOfBirth string `db:"dateOfBirth"`
		Sex         string `db:"sex"`

		// Hospital fields
		HospitalID      sql.NullInt64  `db:"hospital_id"`
		HospitalName    sql.NullString `db:"hospital_name"`
		HospitalAddress sql.NullString `db:"hospital_address"`
		HospitalPhone   sql.NullString `db:"hospital_phone"`
		HospitalEmail   sql.NullString `db:"hospital_email"`

		// Lab test fields
		LabAppointmentID       sql.NullInt64  `db:"lab_appointment_id"`
		TestTypeID             sql.NullInt64  `db:"test_type_id"`
		PickupType             sql.NullString `db:"pickup_type"`
		HomeLocation           sql.NullString `db:"home_location"`
		AdditionalInstructions sql.NullString `db:"additional_instructions"`

		// Test type fields
		TestName        sql.NullString  `db:"test_name"`
		TestDescription sql.NullString  `db:"test_description"`
		TestPrice       sql.NullFloat64 `db:"test_price"`

		// Doctor appointment fields
		DoctorAppointmentID sql.NullInt64  `db:"doctor_appointment_id"`
		DoctorID            sql.NullInt64  `db:"doctor_id"`
		ReasonForVisit      sql.NullString `db:"reason_for_visit"`
		Symptoms            sql.NullString `db:"symptoms"`
		DoctorNotes         sql.NullString `db:"doctor_notes"`

		// Doctor fields
		DoctorName     sql.NullString `db:"doctor_name"`
		Specialization sql.NullString `db:"specialization"`
		DoctorEmail    sql.NullString `db:"doctor_email"`
		DoctorPhone    sql.NullString `db:"doctor_phone"`
	}

	err := api.Deps.DB.GetContext(ctx, &result, query, args...)
	if err != nil {
		log.Println("error fetching detailed appointment", err)
		return model.DetailedAppointment{}, fmt.Errorf("failed to fetch detailed appointment: %w", err)
	}

	// Build the detailed appointment object
	detailed := model.DetailedAppointment{
		ID:                  result.ID,
		UserID:              result.UserID,
		AppointmentType:     result.AppointmentType,
		AppointmentDatetime: result.AppointmentDatetime,
		Status:              result.Status,
		CreatedAt:           result.CreatedAt,
		UpdatedAt:           result.UpdatedAt,

		// User info
		User: &model.UserInfo{
			ID:          result.UserID,
			FirstName:   result.FirstName,
			LastName:    result.LastName,
			Email:       result.Email,
			DateOfBirth: result.DateOfBirth,
			Sex:         result.Sex,
		},
	}

	// Add hospital info if available
	if result.HospitalID.Valid {
		detailed.Hospital = &model.HospitalInfo{
			ID:      int(result.HospitalID.Int64),
			Name:    result.HospitalName.String,
			Address: result.HospitalAddress.String,
			Phone:   result.HospitalPhone.String,
			Email:   result.HospitalEmail.String,
		}
	}

	// Add lab test details if lab appointment
	if result.AppointmentType == "lab_test" && result.LabAppointmentID.Valid {
		detailed.LabTestDetails = &model.LabTestAppointmentDetails{
			ID:                     int(result.LabAppointmentID.Int64),
			AppointmentID:          result.ID,
			TestTypeID:             int(result.TestTypeID.Int64),
			PickupType:             result.PickupType.String,
			AdditionalInstructions: &result.AdditionalInstructions.String,
		}

		if result.HospitalID.Valid {
			hospitalID := int(result.HospitalID.Int64)
			detailed.LabTestDetails.HospitalID = &hospitalID
		}

		if result.HomeLocation.Valid {
			detailed.LabTestDetails.HomeLocation = &result.HomeLocation.String
		}

		// Add test type info
		if result.TestTypeID.Valid {
			detailed.TestType = &model.TestTypeInfo{
				ID:          int(result.TestTypeID.Int64),
				Name:        result.TestName.String,
				Description: result.TestDescription.String,
				Price:       result.TestPrice.Float64,
			}
		}
	}

	// Add doctor details if doctor appointment
	if result.AppointmentType == "doctor" && result.DoctorAppointmentID.Valid {
		detailed.DoctorDetails = &model.DoctorAppointmentDetails{
			ID:            int(result.DoctorAppointmentID.Int64),
			AppointmentID: result.ID,
		}

		if result.ReasonForVisit.Valid {
			detailed.DoctorDetails.ReasonForVisit = &result.ReasonForVisit.String
		}
		if result.Symptoms.Valid {
			detailed.DoctorDetails.Symptoms = &result.Symptoms.String
		}
		if result.DoctorNotes.Valid {
			detailed.DoctorDetails.AdditionalNotes = &result.DoctorNotes.String
		}

		// Add doctor info
		if result.DoctorID.Valid {
			detailed.Doctor = &model.DoctorInfo{
				ID:             int(result.DoctorID.Int64),
				Name:           result.DoctorName.String,
				Specialization: result.Specialization.String,
				Email:          result.DoctorEmail.String,
				Phone:          result.DoctorPhone.String,
			}
		}
	}

	// Get status history
	statusHistory, _ := api.GetAppointmentStatusHistoryRepo(ctx, appointmentID)
	detailed.StatusHistory = statusHistory

	// Get reschedule offers
	rescheduleOffers, _ := api.GetRescheduleOffersRepo(ctx, appointmentID)
	detailed.RescheduleOffers = rescheduleOffers

	return detailed, nil
}

// Helper function to get reschedule offers
func (api *API) GetRescheduleOffersRepo(ctx context.Context, appointmentID int) ([]model.RescheduleOffer, error) {
	query := `
		SELECT id, appointment_id, proposed_date, proposed_time, admin_notes, status, created_at, updated_at
		FROM reschedule_offers
		WHERE appointment_id = ?
		ORDER BY created_at DESC`

	var offers []model.RescheduleOffer
	err := api.Deps.DB.SelectContext(ctx, &offers, query, appointmentID)
	if err != nil {
		log.Println("error fetching reschedule offers", err)
		return nil, err
	}

	return offers, nil
}

func (api *API) UpdateAppointmentStatus(ctx context.Context, appointmentID int, status string, notes *string, changedByUserID *int) error {
	return api.Deps.DB.RunInTx(ctx, func(tx *sqlx.Tx) error {
		// Update appointment status
		updateQuery := `UPDATE appointments SET status = ?, updated_at = NOW() WHERE id = ?`
		_, err := tx.ExecContext(ctx, updateQuery, status, appointmentID)
		if err != nil {
			return err
		}

		// Log status change in history
		historyQuery := `
			INSERT INTO appointment_status_history (appointment_id, status, notes, changed_by_user_id)
			VALUES (?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, historyQuery, appointmentID, status, notes, changedByUserID)
		return err
	})
}

func (api *API) RejectAppointment(ctx context.Context, appointmentID int, rejectionReason, notes *string, changedByUserID *int) error {
	return api.Deps.DB.RunInTx(ctx, func(tx *sqlx.Tx) error {
		// Update appointment with rejection
		updateQuery := `
			UPDATE appointments
			SET status = ?, rejection_reason = ?, admin_notes = ?, updated_at = NOW()
			WHERE id = ?`
		_, err := tx.ExecContext(ctx, updateQuery, string(model.StatusRejected), rejectionReason, notes, appointmentID)
		if err != nil {
			return err
		}

		// Log status change in history
		historyQuery := `
			INSERT INTO appointment_status_history (appointment_id, status, notes, changed_by_user_id)
			VALUES (?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, historyQuery, appointmentID, string(model.StatusRejected), notes, changedByUserID)
		return err
	})
}

func (api *API) CreateRescheduleOffer(ctx context.Context, appointmentID int, proposedDate, proposedTime string, notes *string, changedByUserID *int) error {
	return api.Deps.DB.RunInTx(ctx, func(tx *sqlx.Tx) error {
		// Create reschedule offer
		offerQuery := `
			INSERT INTO reschedule_offers (appointment_id, proposed_date, proposed_time, admin_notes)
			VALUES (?, ?, ?, ?)`
		_, err := tx.ExecContext(ctx, offerQuery, appointmentID, proposedDate, proposedTime, notes)
		if err != nil {
			return err
		}

		// Update appointment status
		updateQuery := `UPDATE appointments SET status = ?, updated_at = NOW() WHERE id = ?`
		_, err = tx.ExecContext(ctx, updateQuery, string(model.StatusRescheduleOffered), appointmentID)
		if err != nil {
			return err
		}

		// Log status change in history
		historyQuery := `
			INSERT INTO appointment_status_history (appointment_id, status, notes, changed_by_user_id)
			VALUES (?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, historyQuery, appointmentID, string(model.StatusRescheduleOffered), notes, changedByUserID)
		return err
	})
}

func (api *API) UpdateAppointmentAdminNotes(ctx context.Context, appointmentID int, notes string) error {
	query := `UPDATE appointments SET admin_notes = ?, updated_at = NOW() WHERE id = ?`
	_, err := api.Deps.DB.ExecContext(ctx, query, notes, appointmentID)
	return err
}

func (api *API) AcceptRescheduleOfferRepo(ctx context.Context, appointmentID, userID, offerID int) error {
	return api.Deps.DB.RunInTx(ctx, func(tx *sqlx.Tx) error {
		// Get the reschedule offer details
		var offer model.RescheduleOffer
		offerQuery := `
			SELECT proposed_date, proposed_time
			FROM reschedule_offers
			WHERE id = ? AND appointment_id = ? AND status = 'pending'`
		err := tx.GetContext(ctx, &offer, offerQuery, offerID, appointmentID)
		if err != nil {
			return err
		}

		// Update reschedule offer status
		updateOfferQuery := `UPDATE reschedule_offers SET status = 'accepted', updated_at = NOW() WHERE id = ?`
		_, err = tx.ExecContext(ctx, updateOfferQuery, offerID)
		if err != nil {
			return err
		}

		// Update appointment with new date/time
		updateAppointmentQuery := `
			UPDATE appointments
			SET appointment_date = ?, appointment_time = ?, status = ?, updated_at = NOW()
			WHERE id = ?`
		_, err = tx.ExecContext(ctx, updateAppointmentQuery, offer.ProposedDate, offer.ProposedTime, string(model.StatusRescheduleAccepted), appointmentID)
		if err != nil {
			return err
		}

		// Log status change
		historyQuery := `
			INSERT INTO appointment_status_history (appointment_id, status, notes, changed_by_user_id)
			VALUES (?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, historyQuery, appointmentID, string(model.StatusRescheduleAccepted), "User accepted reschedule offer", userID)
		return err
	})
}

func (api *API) RejectRescheduleOfferRepo(ctx context.Context, appointmentID, userID, offerID int, reason *string) error {
	return api.Deps.DB.RunInTx(ctx, func(tx *sqlx.Tx) error {
		// Update reschedule offer status
		updateOfferQuery := `UPDATE reschedule_offers SET status = 'rejected', updated_at = NOW() WHERE id = ?`
		_, err := tx.ExecContext(ctx, updateOfferQuery, offerID)
		if err != nil {
			return err
		}

		// Update appointment status back to pending for admin to review
		updateAppointmentQuery := `UPDATE appointments SET status = ?, updated_at = NOW() WHERE id = ?`
		_, err = tx.ExecContext(ctx, updateAppointmentQuery, string(model.StatusPending), appointmentID)
		if err != nil {
			return err
		}

		// Log status change
		notes := "User rejected reschedule offer"
		if reason != nil {
			notes += ": " + *reason
		}
		historyQuery := `
			INSERT INTO appointment_status_history (appointment_id, status, notes, changed_by_user_id)
			VALUES (?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, historyQuery, appointmentID, string(model.StatusPending), notes, userID)
		return err
	})
}

func (api *API) GetAppointmentStatusHistoryRepo(ctx context.Context, appointmentID int) ([]model.AppointmentStatusLog, error) {
	query := `
		SELECT id, appointment_id, status, notes, changed_by_user_id, changed_at
		FROM appointment_status_history
		WHERE appointment_id = ?
		ORDER BY changed_at ASC`

	var history []model.AppointmentStatusLog
	err := api.Deps.DB.SelectContext(ctx, &history, query, appointmentID)
	if err != nil {
		log.Println("error fetching appointment status history", err)
		return nil, fmt.Errorf("failed to fetch status history: %w", err)
	}

	return history, nil
}
