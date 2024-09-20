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
