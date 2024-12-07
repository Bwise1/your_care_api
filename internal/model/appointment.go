package model

import "time"

type Appointment struct {
	ID                       int                        `json:"id" db:"id"`
	UserID                   int                        `json:"user_id" db:"user_id"`
	DoctorID                 *int                       `json:"doctor_id,omitempty" db:"doctor_id"`     // Pointer to handle nullability
	LabTestID                *int                       `json:"lab_test_id,omitempty" db:"lab_test_id"` // Pointer to handle nullability
	AppointmentType          string                     `json:"appointment_type" db:"appointment_type"`
	AppointmentDate          string                     `json:"appointment_date" db:"appointment_date"`
	AppointmentTime          string                     `json:"appointment_time" db:"appointment_time"`
	Status                   string                     `json:"status" db:"status"`
	LabTestDetails           *LabTestAppointmentDetails `json:"lab_test_details,omitempty"`
	DoctorAppointmentDetails *DoctorAppointmentDetails  `json:"doctor_appointment_details,omitempty"`
	CreatedAt                string                     `json:"created_at" db:"created_at"`
	UpdatedAt                string                     `json:"updated_at" db:"updated_at"`
}

type AppointmentDetails struct {
	ID                  int                 `db:"id" json:"id"`
	UserID              int                 `db:"user_id" json:"user_id"`
	AppointmentType     string              `db:"appointment_type" json:"appointment_type"`
	AppointmentDatetime *time.Time          `db:"appointment_datetime" json:"appointment_datetime"`
	Status              string              `db:"status" json:"status"`
	CreatedAt           *time.Time          `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt           *time.Time          `db:"updated_at,omitempty" json:"updated_at,omitempty"`
	DoctorDetails       *DoctorAppointment  `db:"doctor_details,omitempty" json:"doctor_details,omitempty"`
	LabTestDetails      *LabTestAppointment `db:"lab_test_details,omitempty" json:"lab_test_details,omitempty"`
}

type AppointmentRow struct {
	ID                  int        `db:"id"`
	UserID              int        `db:"user_id"`
	AppointmentType     string     `db:"appointment_type"`
	AppointmentDatetime *time.Time `db:"appointment_datetime"`
	Status              string     `db:"status"`
	CreatedAt           *time.Time `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`
	DoctorDetailsJSON   []byte     `db:"doctor_details"`
	LabTestDetailsJSON  []byte     `db:"lab_test_details"`
}

type DoctorAppointment struct {
	ID              int     `json:"id"`
	AppointmentID   int     `json:"appointment_id"`
	DoctorID        int     `json:"doctor_id,omitempty"`
	ReasonForVisit  *string `json:"reason_for_visit"`
	Symptoms        *string `json:"symptoms"`
	AdditionalNotes *string `json:"additional_notes"`
}

type LabTestAppointment struct {
	ID                     int     `json:"id,omitempty"`
	AppointmentID          int     `json:"appointment_id,omitempty" `
	TestTypeID             int     `json:"test_type_id"`
	PickupType             string  `json:"pickup_type"`
	HomeLocation           *string `json:"home_location,omitempty"`
	HospitalID             *int    `json:"hospital_id,omitempty"`
	AdditionalInstructions *string `json:"additional_instructions,omitempty"`
}

type LabTestAppointmentDetails struct {
	ID                     int     `json:"id"`
	AppointmentID          int     `json:"appointmentId"`
	PickupType             string  `json:"pickupType"`
	HomeLocation           *string `json:"homeLocation,omitempty"`
	TestTypeID             int     `json:"testTypeId"`
	HospitalID             *int    `json:"hospitalId,omitempty"`
	AdditionalInstructions *string `json:"additionalInstructions,omitempty"`
}

type DoctorAppointmentDetails struct {
	ID              int     `json:"id"`
	AppointmentID   int     `json:"appointment_id"`
	ReasonForVisit  *string `json:"reason_for_visit"`
	Symptoms        *string `json:"symptoms"`
	AdditionalNotes *string `json:"additional_notes"`
}

type LabAppointmentReq struct {
	UserID                 int     `json:"user"`
	DoctorID               *int    `json:"doctor,omitempty"`   // Pointer to handle nullability
	HospitalID             *int    `json:"hospital,omitempty"` // Pointer to handle nullability
	LabTestID              int     `json:"lab_test"`
	AppointmentDate        string  `json:"appointment_date"`
	AppointmentTime        string  `json:"appointment_time"`
	PickupType             string  `json:"pickup_type"`
	HomeLocation           *string `json:"home_location,omitempty"` // Pointer to handle nullability
	TestTypeID             int     `json:"test_type"`
	AdditionalInstructions *string `json:"additional_instructions,omitempty"` // Pointer to handle nullability
}

type DoctorAppointmentReq struct {
	UserID                 int     `json:"user"`
	DoctorID               *int    `json:"doctor,omitempty"`   // Pointer to handle nullability
	HospitalID             *int    `json:"hospital,omitempty"` // Pointer to handle nullability
	LabTestID              int     `json:"lab_test"`
	AppointmentDate        string  `json:"appointment_date"`
	AppointmentTime        string  `json:"appointment_time"`
	PickupType             string  `json:"pickup_type"`
	HomeLocation           *string `json:"home_location,omitempty"` // Pointer to handle nullability
	TestTypeID             int     `json:"test_type"`
	AdditionalInstructions *string `json:"additional_instructions,omitempty"` // Pointer to handle nullability
}

// Request structs for API
type CreateDoctorAppointmentRequest struct {
	UserID          int    `json:"user_id"`
	AppointmentDate string `json:"appointment_date"` // Format: "2024-10-17"
	AppointmentTime string `json:"appointment_time"` // Format: "14:30"
	DoctorID        int    `json:"doctor_id"`
	ReasonForVisit  string `json:"reason_for_visit"`
	Symptoms        string `json:"symptoms"`
	AdditionalNotes string `json:"additional_notes,omitempty"`
}

type CreateLabTestAppointmentRequest struct {
	UserID                 int     `json:"user"`
	AppointmentDate        string  `json:"appointment_date"`
	TestTypeID             int     `json:"test_type"`
	PickupType             string  `json:"pickup_type"`
	HomeLocation           *string `json:"home_location,omitempty"`
	AdditionalInstructions *string `json:"additional_instructions,omitempty"`
	HospitalID             *int    `json:"hospital,omitempty"`
}

type AppointmentFilter struct {
	UserID   int    `json:"user_id"`
	Date     string `json:"date"`
	Upcoming bool   `json:"upcoming"`
	History  bool   `json:"history"`
}
