package status

import (
	"context"
	"fmt"
	"strconv"

	modzysdk "github.com/modzy/sdk-go"
	modzysdkmodel "github.com/modzy/sdk-go/model"
)

type TopModel struct {
	Errors            []error
	Alerts            []modzysdk.AlertSummary
	DataProcessed     modzysdk.GetDataProcessedOutput
	Predictions       modzysdk.GetPredictionsMadeOutput
	ActiveUsers       []modzysdkmodel.ActiveUserSummary
	ActiveModels      []modzysdkmodel.ActiveModelSummary
	CPUOverallLast    float64
	CPUOverallAverage float64
	Users             []modzysdkmodel.AccountingUser
	License           modzysdkmodel.License
	EnginesProcessing int
	LatestModels      []modzysdkmodel.ModelDetails
}

// FetchTop -
func (f *StandardFetcher) FetchTop(ctx context.Context) TopModel {
	top := TopModel{}

	// alerts
	if alerts, err := f.client.Dashboard().GetAlerts(ctx, &modzysdk.GetAlertsInput{}); err != nil {
		top.Errors = append(top.Errors, fmt.Errorf("Failed to fetch alerts: %v", err))
	} else {
		top.Alerts = alerts.Alerts
	}

	// processed
	if processed, err := f.client.Dashboard().GetDataProcessed(ctx, &modzysdk.GetDataProcessedInput{}); err != nil {
		top.Errors = append(top.Errors, fmt.Errorf("Failed to fetch processed data: %v", err))
	} else {
		top.DataProcessed = *processed
	}

	// predictions-made
	if predictions, err := f.client.Dashboard().GetPredictionsMade(ctx, &modzysdk.GetPredictionsMadeInput{}); err != nil {
		top.Errors = append(top.Errors, fmt.Errorf("Failed to fetch predictions: %v", err))
	} else {
		top.Predictions = *predictions
	}

	// active-users
	if activeUsers, err := f.client.Dashboard().GetActiveUsers(ctx, &modzysdk.GetActiveUsersInput{}); err != nil {
		top.Errors = append(top.Errors, fmt.Errorf("Failed to fetch active users: %v", err))
	} else {
		top.ActiveUsers = activeUsers.Users
	}

	// active-models
	if activeModels, err := f.client.Dashboard().GetActiveModels(ctx, &modzysdk.GetActiveModelsInput{}); err != nil {
		top.Errors = append(top.Errors, fmt.Errorf("Failed to fetch active models: %v", err))
	} else {
		top.ActiveModels = activeModels.Models
	}

	// cpu-overall-usage
	if cpuOverall, err := f.client.Dashboard().GetPrometheusMetric(ctx, &modzysdk.GetPrometheusMetricInput{
		Metric: modzysdk.PrometheusMetricTypeCPUOverallUsage,
	}); err != nil {
		top.Errors = append(top.Errors, fmt.Errorf("Failed to fetch CPU usage: %v", err))
	} else {
		var num float64
		var total float64
		var last float64
		for _, promVal := range cpuOverall.Values {
			v, err := strconv.ParseFloat(promVal.Value, 64)
			if err == nil {
				num += 1.0
				total += v
				last = v
			}
		}
		top.CPUOverallLast = last
		top.CPUOverallAverage = total / num
	}

	// accounting-users
	if accountingUsers, err := f.client.Accounting().ListAccountingUsers(ctx, &modzysdk.ListAccountingUsersInput{}); err != nil {
		top.Errors = append(top.Errors, fmt.Errorf("Failed to fetch users: %v", err))
	} else {
		top.Users = accountingUsers.Users
	}

	// license
	if license, err := f.client.Accounting().GetLicense(ctx); err != nil {
		top.Errors = append(top.Errors, fmt.Errorf("Failed to fetch license: %v", err))
	} else {
		top.License = license.License
	}

	// engines
	if processingModels, err := f.client.Resources().GetProcessingModels(ctx); err != nil {
		top.Errors = append(top.Errors, fmt.Errorf("Failed to fetch engines: %v", err))
	} else {
		tot := 0
		for _, m := range processingModels.Models {
			tot += len(m.Engines)
		}
		top.EnginesProcessing = tot
	}

	// latest models
	if latestModels, err := f.client.Models().GetLatestModels(ctx); err != nil {
		top.Errors = append(top.Errors, fmt.Errorf("Failed to fetch latest models: %v", err))
	} else {
		top.LatestModels = latestModels.Models
	}

	return top
}
