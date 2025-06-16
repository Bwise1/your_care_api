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
