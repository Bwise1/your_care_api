package model

type Appointment struct {
	ID              int    `json:"id" db:"id"`
	UserID          int    `json:"user_id" db:"user_id"`
	DoctorID        *int   `json:"doctor_id,omitempty" db:"doctor_id"`     // Pointer to handle nullability
	LabTestID       *int   `json:"lab_test_id,omitempty" db:"lab_test_id"` // Pointer to handle nullability
	AppointmentType string `json:"appointment_type" db:"appointment_type"`
	AppointmentDate string `json:"appointment_date" db:"appointment_date"`
	AppointmentTime string `json:"appointment_time" db:"appointment_time"`
	Status          string `json:"status" db:"status"`
	CreatedAt       string `json:"created_at" db:"created_at"`
	UpdatedAt       string `json:"updated_at" db:"updated_at"`
}

type LabTestAppointmentDetails struct {
	ID            int     `json:"id"`
	AppointmentID int     `json:"appointmentId"`
	PickupType    string  `json:"pickupType"`
	HomeLocation  *string `json:"homeLocation,omitempty"` // Optional for hospital tests
	TestType      string  `json:"testType"`
}

type DoctorAppointmentDetails struct {
	ID              int    `json:"id" db:"id"`
	AppointmentID   int    `json:"appointment_id" db:"appointment_id"`
	ReasonForVisit  string `json:"reason_for_visit" db:"reason_for_visit"`
	Symptoms        string `json:"symptoms" db:"symptoms"`
	AdditionalNotes string `json:"additional_notes" db:"additional_notes"`
}
