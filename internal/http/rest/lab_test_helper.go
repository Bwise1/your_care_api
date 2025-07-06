package rest

import (
	"context"
	"log"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/bwise1/your_care_api/util/values"
)

func (api *API) GetAllLabTestsHelper() ([]model.LabTest, string, string, error) {
	tests, err := api.GetAllLabTestsRepo(context.TODO())
	if err != nil {
		log.Println(err)
		return nil, values.Error, "Failed to fetch lab tests", err
	}
	return tests, values.Success, "Fetched lab tests", nil
}

func (api *API) GetAvailableTestsForSelectionHelper() ([]model.TestForSelection, string, string, error) {
	tests, err := api.GetAvailableTestsForSelectionRepo(context.TODO())
	if err != nil {
		log.Println(err)
		return nil, values.Error, "Failed to fetch available tests", err
	}
	return tests, values.Success, "Fetched available tests for selection", nil
}

func (api *API) CreateLabTestHelper(req model.LabTest) (model.LabTest, string, string, error) {
	id, err := api.CreateLabTestRepo(context.TODO(), req)
	if err != nil {
		return model.LabTest{}, values.Error, "Failed to create lab test", err
	}
	req.ID = id
	return req, values.Created, "Lab test created", nil
}

func (api *API) UpdateLabTestHelper(req model.LabTest) (string, string, error) {
	err := api.UpdateLabTestRepo(context.TODO(), req)
	if err != nil {
		return values.Error, "Failed to update lab test", err
	}
	return values.Success, "Lab test updated", nil
}

func (api *API) DeleteLabTestHelper(id int) (string, string, error) {
	err := api.DeleteLabTestRepo(context.TODO(), id)
	if err != nil {
		return values.Error, "Failed to delete lab test", err
	}
	return values.Success, "Lab test deleted", nil
}
