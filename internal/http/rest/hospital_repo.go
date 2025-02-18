package rest

import (
	"context"

	"github.com/bwise1/your_care_api/internal/model"
)

func (api *API) GetAllHospitals(ctx context.Context) ([]model.Hospital, error) {
	stmt := `SELECT
		id,
		name,
		address,
		phone,
		email
	FROM hospitals`

	rows, err := api.Deps.DB.QueryContext(ctx, stmt)
	if err != nil {
		return []model.Hospital{}, err
	}
	defer rows.Close()

	var hospitals []model.Hospital
	for rows.Next() {
		var h model.Hospital
		err := rows.Scan(&h.ID, &h.Name, &h.Address, &h.Phone, &h.Email)
		if err != nil {
			return nil, err
		}
		hospitals = append(hospitals, h)
	}

	return hospitals, nil
}

func (api *API) GetLabTestsByHospital(ctx context.Context, hospitalID int) ([]model.LabTest, error) {

	stmt := `SELECT
		id,
		hospital_id,
		name,
		description,
		price
	FROM lab_tests
	WHERE hospital_id = ?`

	rows, err := api.Deps.DB.QueryContext(ctx, stmt, hospitalID)
	if err != nil {
		return []model.LabTest{}, err
	}
	defer rows.Close()

	var tests []model.LabTest
	for rows.Next() {
		var t model.LabTest
		err := rows.Scan(&t.ID, &t.HospitalID, &t.Name, &t.Description, &t.Price)
		if err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}

	return tests, nil
}

func (api *API) CreateHospitalRepo(ctx context.Context, req model.Hospital) (int, error) {
	stmt := `INSERT INTO hospitals (
        name,
        address,
        phone,
        email
    ) VALUES (?, ?, ?, ?)`

	result, err := api.Deps.DB.ExecContext(ctx, stmt, req.Name, req.Address, req.Phone, req.Email)
	if err != nil {
		return 0, err
	}

	hospitalID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(hospitalID), nil
}

func (api *API) DeleteHospitalRepo(ctx context.Context, hospitalID int) error {
	stmt := `DELETE FROM hospitals WHERE id = ?`

	_, err := api.Deps.DB.ExecContext(ctx, stmt, hospitalID)
	if err != nil {
		return err
	}

	return nil
}
