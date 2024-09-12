package model

type Doctor struct {
	ID             int    `json:"id" db:"id"`
	HospitalID     int    `json:"hospital_id" db:"hospital_id"`
	Name           string `json:"name" db:"name"`
	Specialization string `json:"specialization" db:"specialization"`
	Email          string `json:"email" db:"email"`
	Phone          string `json:"phone" db:"phone"`
	AvailableFrom  string `json:"available_from" db:"available_from"`
	AvailableTo    string `json:"available_to" db:"available_to"`
}
