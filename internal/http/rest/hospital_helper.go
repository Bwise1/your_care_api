package rest

import (
	"context"
	"fmt"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/bwise1/your_care_api/util/values"
)

func (api *API) GetHospitals_H() ([]model.Hospital, string, string, error) {

	var err error
	var ctx = context.TODO()

	hospitals, err := api.GetAllHospitals(ctx)
	if err != nil {
		return []model.Hospital{}, values.Error, fmt.Sprintf("%s [GeHo]", values.SystemErr), err
	}

	return hospitals, values.Success, "Fetched hospitals successfully", nil
}

func (api *API) GetLabTestsByHospital_H(hospitalID int) ([]model.LabTest, string, string, error) {

	var err error
	var ctx = context.TODO()

	tests, err := api.GetLabTestsByHospital(ctx, hospitalID)
	if err != nil {
		return []model.LabTest{}, values.Error, fmt.Sprintf("%s [GeLa]", values.SystemErr), err
	}

	return tests, values.Success, "Fetched lab tests successfully", nil
}

func (api *API) CreateHospital_H(req model.Hospital) (model.Hospital, string, string, error) {
	var err error
	var ctx = context.TODO()

	hospitalID, err := api.CreateHospitalRepo(ctx, req)
	if err != nil {
		return model.Hospital{}, values.Error, fmt.Sprintf("%s [CrHo]", values.SystemErr), err
	}

	req.ID = hospitalID
	return req, values.Created, "Hospital created successfully", nil
}

func (api *API) DeleteHospital_H(hospitalID int) (string, string, error) {
	var err error
	var ctx = context.TODO()

	err = api.DeleteHospitalRepo(ctx, hospitalID)
	if err != nil {
		return values.Error, fmt.Sprintf("%s [DlHo]", values.SystemErr), err
	}

	return values.Success, "Hospital deleted successfully", nil
}
