package model

import "time"

// Appointment status constants
type AppointmentStatus string

const (
	StatusPending          AppointmentStatus = "pending"
	StatusConfirmed        AppointmentStatus = "confirmed"
	StatusScheduled        AppointmentStatus = "scheduled"
	StatusRescheduleOffered AppointmentStatus = "reschedule_offered"
	StatusRescheduleAccepted AppointmentStatus = "reschedule_accepted"
	StatusInProgress       AppointmentStatus = "in_progress"
	StatusCompleted        AppointmentStatus = "completed"
	StatusCanceled         AppointmentStatus = "canceled"
	StatusRejected         AppointmentStatus = "rejected"
	StatusNoShow           AppointmentStatus = "no_show"
)

// Appointment type constants
type AppointmentType string

const (
	TypeDoctor AppointmentType = "doctor"
	TypeLab    AppointmentType = "lab_test"
	TypeIVF    AppointmentType = "ivf"
)

type Appointment struct {
	ID                       int                        `json:"id" db:"id"`
	UserID                   int                        `json:"user_id" db:"user_id"`
	DoctorID                 *int                       `json:"doctor_id,omitempty" db:"doctor_id"`
	LabTestID                *int                       `json:"lab_test_id,omitempty" db:"lab_test_id"`
	ProviderID               *int                       `json:"provider_id,omitempty" db:"provider_id"`
	AppointmentType          AppointmentType            `json:"appointment_type" db:"appointment_type"`
	AppointmentDate          string                     `json:"appointment_date" db:"appointment_date"`
	AppointmentTime          string                     `json:"appointment_time" db:"appointment_time"`
	Status                   AppointmentStatus          `json:"status" db:"status"`
	AdminNotes               *string                    `json:"admin_notes,omitempty" db:"admin_notes"`
	UserNotes                *string                    `json:"user_notes,omitempty" db:"user_notes"`
	RejectionReason          *string                    `json:"rejection_reason,omitempty" db:"rejection_reason"`
	LabTestDetails           *LabTestAppointmentDetails `json:"lab_test_details,omitempty"`
	DoctorAppointmentDetails *DoctorAppointmentDetails  `json:"doctor_appointment_details,omitempty"`
	IVFDetails               *IVFAppointmentDetails     `json:"ivf_details,omitempty"`
	RescheduleOffers         []RescheduleOffer          `json:"reschedule_offers,omitempty"`
	StatusHistory            []AppointmentStatusLog     `json:"status_history,omitempty"`
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

// New structs for enhanced appointment system
type RescheduleOffer struct {
	ID              int                `json:"id" db:"id"`
	AppointmentID   int                `json:"appointment_id" db:"appointment_id"`
	ProposedDate    string             `json:"proposed_date" db:"proposed_date"`
	ProposedTime    string             `json:"proposed_time" db:"proposed_time"`
	AdminNotes      *string            `json:"admin_notes,omitempty" db:"admin_notes"`
	Status          string             `json:"status" db:"status"`
	CreatedAt       *time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time         `json:"updated_at" db:"updated_at"`
}

type AppointmentStatusLog struct {
	ID              int       `json:"id" db:"id"`
	AppointmentID   int       `json:"appointment_id" db:"appointment_id"`
	Status          string    `json:"status" db:"status"`
	Notes           *string   `json:"notes,omitempty" db:"notes"`
	ChangedByUserID *int      `json:"changed_by_user_id,omitempty" db:"changed_by_user_id"`
	ChangedAt       *time.Time `json:"changed_at" db:"changed_at"`
}

type IVFAppointmentDetails struct {
	ID                  int     `json:"id" db:"id"`
	AppointmentID       int     `json:"appointment_id" db:"appointment_id"`
	TreatmentType       *string `json:"treatment_type,omitempty" db:"treatment_type"`
	CycleDay            *int    `json:"cycle_day,omitempty" db:"cycle_day"`
	SpecialInstructions *string `json:"special_instructions,omitempty" db:"special_instructions"`
	PreparationNotes    *string `json:"preparation_notes,omitempty" db:"preparation_notes"`
}

// Admin request structs
type AdminAppointmentAction struct {
	Notes           *string    `json:"notes,omitempty"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`
	ProposedDate    *string    `json:"proposed_date,omitempty"`
	ProposedTime    *string    `json:"proposed_time,omitempty"`
}

type AdminStatusUpdateRequest struct {
	AppointmentID   int     `json:"appointmentId"`
	Status          string  `json:"status"`
	AdminNotes      *string `json:"adminNotes,omitempty"`
	RejectionReason *string `json:"rejectionReason,omitempty"`
	NewDateTime     *string `json:"newDateTime,omitempty"`
}

type AdminAppointmentFilter struct {
	Status          []string `json:"status,omitempty"`
	AppointmentType []string `json:"appointment_type,omitempty"`
	DateFrom        *string  `json:"date_from,omitempty"`
	DateTo          *string  `json:"date_to,omitempty"`
	ProviderID      *int     `json:"provider_id,omitempty"`
	Page            int      `json:"page"`
	Limit           int      `json:"limit"`
}

// Detailed appointment structures for get by ID
type DetailedAppointment struct {
	ID                  int                        `json:"id"`
	UserID              int                        `json:"user_id"`
	AppointmentType     string                     `json:"appointment_type"`
	AppointmentDatetime *time.Time                 `json:"appointment_datetime"`
	Status              string                     `json:"status"`
	CreatedAt           *time.Time                 `json:"created_at"`
	UpdatedAt           *time.Time                 `json:"updated_at"`
	
	// User details
	User                *UserInfo                  `json:"user"`
	
	// Appointment specific details
	DoctorDetails       *DoctorAppointmentDetails  `json:"doctor_details,omitempty"`
	LabTestDetails      *LabTestAppointmentDetails `json:"lab_test_details,omitempty"`
	
	// Related information
	Hospital            *HospitalInfo              `json:"hospital,omitempty"`
	TestType            *TestTypeInfo              `json:"test_type,omitempty"`
	Doctor              *DoctorInfo                `json:"doctor,omitempty"`
	
	// Status and actions
	StatusHistory       []AppointmentStatusLog     `json:"status_history"`
	RescheduleOffers    []RescheduleOffer          `json:"reschedule_offers"`
	NextActions         []string                   `json:"next_actions"`
}

type UserInfo struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`
	Sex         string `json:"sex"`
	Phone       string `json:"phone,omitempty"`
}

type HospitalInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type TestTypeInfo struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type DoctorInfo struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
}

// User response structs  
type UserAppointmentResponse struct {
	Appointment   AppointmentDetails      `json:"appointment"`
	StatusHistory []AppointmentStatusLog  `json:"status_history"`
	NextActions   []string                `json:"next_actions"`
	EstimatedTime *time.Time              `json:"estimated_time,omitempty"`
}

// Enhanced detailed response for get by ID
type DetailedAppointmentResponse struct {
	Appointment DetailedAppointment `json:"appointment"`
}

type RescheduleResponse struct {
	OfferID int     `json:"offer_id"`
	Reason  *string `json:"reason,omitempty"`
}
