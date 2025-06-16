package rest

import (
	"context"

	"github.com/bwise1/your_care_api/internal/model"
)

func (api *API) GetAllLabTestsRepo(ctx context.Context) ([]model.LabTest, error) {
	stmt := `SELECT id, name, description FROM lab_tests`
	rows, err := api.Deps.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []model.LabTest
	for rows.Next() {
		var t model.LabTest
		if err := rows.Scan(&t.ID, &t.Name, &t.Description); err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}
	return tests, nil
}

func (api *API) CreateLabTestRepo(ctx context.Context, req model.LabTest) (int, error) {
	stmt := `INSERT INTO lab_tests (hospital_id, name, description, price) VALUES (?, ?, ?, ?)`
	result, err := api.Deps.DB.ExecContext(ctx, stmt, req.HospitalID, req.Name, req.Description, req.Price)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (api *API) UpdateLabTestRepo(ctx context.Context, req model.LabTest) error {
	stmt := `UPDATE lab_tests SET hospital_id = ?, name = ?, description = ?, price = ? WHERE id = ?`
	_, err := api.Deps.DB.ExecContext(ctx, stmt, req.HospitalID, req.Name, req.Description, req.Price, req.ID)
	return err
}

func (api *API) DeleteLabTestRepo(ctx context.Context, id int) error {
	stmt := `DELETE FROM lab_tests WHERE id = ?`
	_, err := api.Deps.DB.ExecContext(ctx, stmt, id)
	return err
}
