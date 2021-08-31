package status

import (
	"context"
	"fmt"

	modzysdk "github.com/modzy/sdk-go"
	modzysdkmodel "github.com/modzy/sdk-go/model"
)

type ProjectModel struct {
	Errors        []error
	DataProcessed modzysdk.GetDataProcessedOutput
	Predictions   modzysdk.GetPredictionsMadeOutput
	ActiveModels  []modzysdkmodel.ActiveModelSummary
}

// FetchProject -
func (f *StandardFetcher) FetchProject(ctx context.Context, projectID string) ProjectModel {
	out := ProjectModel{}

	project, err := f.client.Accounting().GetProjectDetails(ctx, &modzysdk.GetProjectDetailsInput{
		ProjectID: projectID,
	})
	if err != nil {
		out.Errors = append(out.Errors, fmt.Errorf("Failed to reading project details: %v", err))
	}

	if project != nil {
		accessKey := project.Project.AccessKeys[0].Prefix

		// processed
		if processed, err := f.client.Dashboard().GetDataProcessed(ctx, &modzysdk.GetDataProcessedInput{
			AccessKeyPrefix: accessKey,
		}); err != nil {
			out.Errors = append(out.Errors, fmt.Errorf("Failed to fetch processed data: %v", err))
		} else {
			out.DataProcessed = *processed
		}

		// predictions-made
		if predictions, err := f.client.Dashboard().GetPredictionsMade(ctx, &modzysdk.GetPredictionsMadeInput{
			AccessKeyPrefix: accessKey,
		}); err != nil {
			out.Errors = append(out.Errors, fmt.Errorf("Failed to fetch predictions: %v", err))
		} else {
			out.Predictions = *predictions
		}

		// active-models
		if activeModels, err := f.client.Dashboard().GetActiveModels(ctx, &modzysdk.GetActiveModelsInput{
			AccessKeyPrefix: accessKey,
		}); err != nil {
			out.Errors = append(out.Errors, fmt.Errorf("Failed to fetch active models: %v", err))
		} else {
			out.ActiveModels = activeModels.Models
		}
	}

	return out
}
