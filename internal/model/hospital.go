package model

type LabTest struct {
	ID          int      `json:"id" db:"id"`
	Name        string   `json:"name,omitempty" db:"name"`
	Description string   `json:"description,omitempty" db:"description"`
	Price       *float64 `json:"price,omitempty" db:"price"`
}

type Hospital struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Address string `json:"address" db:"address"`
	Phone   string `json:"phone" db:"phone"`
	Email   string `json:"email" db:"email"`
}

type HospitalLabTest struct {
	ID         int     `json:"id" db:"id"`
	HospitalID int     `json:"hospital_id" db:"hospital_id"`
	LabTestID  int     `json:"lab_test_id" db:"lab_test_id"`
	Name       string  `json:"name" db:"name"`
	Price      float64 `json:"price" db:"price"`
	Details    string  `json:"details" db:"details"`
}

type TestForSelection struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
}
