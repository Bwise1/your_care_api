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
