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

// func (api *API) GetLabTestsByHospital(ctx context.Context, hospitalID int) ([]model.LabTest, error) {

// 	stmt := `SELECT
// 		id,
// 		hospital_id,
// 		name,
// 		description,
// 		price
// 	FROM lab_tests
// 	WHERE hospital_id = ?`

// 	rows, err := api.Deps.DB.QueryContext(ctx, stmt, hospitalID)
// 	if err != nil {
// 		return []model.LabTest{}, err
// 	}
// 	defer rows.Close()

// 	var tests []model.LabTest
// 	for rows.Next() {
// 		var t model.LabTest
// 		err := rows.Scan(&t.ID, &t.HospitalID, &t.Name, &t.Description, &t.Price)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tests = append(tests, t)
// 	}

// 	return tests, nil
// }

func (api *API) GetLabTestsByHospital(ctx context.Context, hospitalID int) ([]model.HospitalLabTest, error) {
	stmt := `SELECT id,
				hospital_id,
				lab_test_id,
				name, price,
				details
			FROM hospital_lab_tests WHERE hospital_id = ?`

	rows, err := api.Deps.DB.QueryContext(ctx, stmt, hospitalID)
	if err != nil {
		return []model.HospitalLabTest{}, err
	}
	defer rows.Close()

	var tests []model.HospitalLabTest
	for rows.Next() {
		var t model.HospitalLabTest
		err := rows.Scan(&t.ID, &t.HospitalID, &t.LabTestID, &t.Name, &t.Price, &t.Details)
		if err != nil {
			return []model.HospitalLabTest{}, err
		}
		tests = append(tests, t)
	}

	// Always return an empty slice if no rows found (not nil)
	if tests == nil {
		return []model.HospitalLabTest{}, nil
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

func (api *API) CreateHospitalLabTestRepo(ctx context.Context, req model.HospitalLabTest) (int, error) {
	stmt := `INSERT INTO hospital_lab_tests (hospital_id, lab_test_id, name, price, details) VALUES (?, ?, ?, ?, ?)`
	result, err := api.Deps.DB.ExecContext(ctx, stmt, req.HospitalID, req.LabTestID, req.Name, req.Price, req.Details)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (api *API) GetAHospitalLabTestsRepo(ctx context.Context, hospitalID int) ([]model.HospitalLabTest, error) {
	stmt := `SELECT id, hospital_id, lab_test_id, name, price, details FROM hospital_lab_tests WHERE hospital_id = ?`
	rows, err := api.Deps.DB.QueryContext(ctx, stmt, hospitalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tests []model.HospitalLabTest
	for rows.Next() {
		var t model.HospitalLabTest
		if err := rows.Scan(&t.ID, &t.HospitalID, &t.LabTestID, &t.Name, &t.Price, &t.Details); err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}
	return tests, nil
}

func (api *API) UpdateHospitalLabTestRepo(ctx context.Context, req model.HospitalLabTest) error {
	stmt := `UPDATE hospital_lab_tests SET name=?, price=?, details=? WHERE id=?`
	_, err := api.Deps.DB.ExecContext(ctx, stmt, req.Name, req.Price, req.Details, req.ID)
	return err
}

func (api *API) DeleteHospitalLabTestRepo(ctx context.Context, id int) error {
	stmt := `DELETE FROM hospital_lab_tests WHERE id=?`
	_, err := api.Deps.DB.ExecContext(ctx, stmt, id)
	return err
}
