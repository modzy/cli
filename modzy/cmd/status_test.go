package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/modzy/cli/internal/status"

	modzysdk "github.com/modzy/sdk-go"
	modzysdkmodel "github.com/modzy/sdk-go/model"
)

func TestStatusFine(t *testing.T) {
	defer (func() { getClient = defaultGetClient })()

	// instead of a bunch of http handles, just use a fake client
	getClient = func() modzysdk.Client {
		return &modzysdk.ClientFake{
			AccountingFunc: func() modzysdk.AccountingClient {
				return &modzysdk.AccountingClientFake{
					GetLicenseFunc: func(ctx context.Context) (*modzysdk.GetLicenseOutput, error) {
						return nil, fmt.Errorf("no license")
					},
					ListAccountingUsersFunc: func(ctx context.Context, input *modzysdk.ListAccountingUsersInput) (*modzysdk.ListAccountingUsersOutput, error) {
						return nil, fmt.Errorf("no accounting users")
					},
				}
			},
			ModelsFunc: func() modzysdk.ModelsClient {
				return &modzysdk.ModelsClientFake{
					GetLatestModelsFunc: func(ctx context.Context) (*modzysdk.GetLatestModelsOutput, error) {
						return nil, fmt.Errorf("no latest models")
					},
				}
			},
			DashboardFunc: func() modzysdk.DashboardClient {
				return &modzysdk.DashboardClientFake{
					GetAlertsFunc: func(ctx context.Context, input *modzysdk.GetAlertsInput) (*modzysdk.GetAlertsOutput, error) {
						return nil, fmt.Errorf("no alerts")
					},
					GetDataProcessedFunc: func(ctx context.Context, input *modzysdk.GetDataProcessedInput) (*modzysdk.GetDataProcessedOutput, error) {
						return nil, fmt.Errorf("no data processed")
					},
					GetPredictionsMadeFunc: func(ctx context.Context, input *modzysdk.GetPredictionsMadeInput) (*modzysdk.GetPredictionsMadeOutput, error) {
						return nil, fmt.Errorf("no predictions")
					},
					GetActiveUsersFunc: func(ctx context.Context, input *modzysdk.GetActiveUsersInput) (*modzysdk.GetActiveUsersOutput, error) {
						return nil, fmt.Errorf("no active users")
					},
					GetActiveModelsFunc: func(ctx context.Context, input *modzysdk.GetActiveModelsInput) (*modzysdk.GetActiveModelsOutput, error) {
						return nil, fmt.Errorf("no active models")
					},
					GetPrometheusMetricFunc: func(ctx context.Context, input *modzysdk.GetPrometheusMetricInput) (*modzysdk.GetPrometheusMetricOutput, error) {
						return nil, fmt.Errorf("no prom")
					},
				}
			},
			ResourcesFunc: func() modzysdk.ResourcesClient {
				return &modzysdk.ResourcesClientFake{
					GetProcessingModelsFunc: func(ctx context.Context) (*modzysdk.GetProcessingModelsOutput, error) {
						return nil, fmt.Errorf("no processing models")
					},
				}
			},
		}
	}

	stdout, _ := runTestCommand(
		[]string{"status"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			// not used
		},
	)

	if !strings.Contains(stdout, "Data processed") {
		t.Fatalf("out not expected: '%s'", stdout)
	}
}

func TestStatusFineLists(t *testing.T) {
	defer (func() { getClient = defaultGetClient })()

	// instead of a bunch of http handles, just use a fake client
	getClient = func() modzysdk.Client {
		return &modzysdk.ClientFake{
			AccountingFunc: func() modzysdk.AccountingClient {
				return &modzysdk.AccountingClientFake{
					GetLicenseFunc: func(ctx context.Context) (*modzysdk.GetLicenseOutput, error) {
						return nil, fmt.Errorf("no license")
					},
					ListAccountingUsersFunc: func(ctx context.Context, input *modzysdk.ListAccountingUsersInput) (*modzysdk.ListAccountingUsersOutput, error) {
						return nil, fmt.Errorf("no accounting users")
					},
				}
			},
			ModelsFunc: func() modzysdk.ModelsClient {
				return &modzysdk.ModelsClientFake{
					GetLatestModelsFunc: func(ctx context.Context) (*modzysdk.GetLatestModelsOutput, error) {
						return nil, fmt.Errorf("no latest models")
					},
				}
			},
			DashboardFunc: func() modzysdk.DashboardClient {
				return &modzysdk.DashboardClientFake{
					GetAlertsFunc: func(ctx context.Context, input *modzysdk.GetAlertsInput) (*modzysdk.GetAlertsOutput, error) {
						return &modzysdk.GetAlertsOutput{
							Alerts: []modzysdk.AlertSummary{
								{Type: "Test"},
							},
						}, nil
					},
					GetDataProcessedFunc: func(ctx context.Context, input *modzysdk.GetDataProcessedInput) (*modzysdk.GetDataProcessedOutput, error) {
						return nil, fmt.Errorf("no data processed")
					},
					GetPredictionsMadeFunc: func(ctx context.Context, input *modzysdk.GetPredictionsMadeInput) (*modzysdk.GetPredictionsMadeOutput, error) {
						return nil, fmt.Errorf("no predictions")
					},
					GetActiveUsersFunc: func(ctx context.Context, input *modzysdk.GetActiveUsersInput) (*modzysdk.GetActiveUsersOutput, error) {
						return &modzysdk.GetActiveUsersOutput{
							Users: []modzysdkmodel.ActiveUserSummary{
								{},
							},
						}, nil
					},
					GetActiveModelsFunc: func(ctx context.Context, input *modzysdk.GetActiveModelsInput) (*modzysdk.GetActiveModelsOutput, error) {
						return &modzysdk.GetActiveModelsOutput{
							Models: []modzysdkmodel.ActiveModelSummary{
								{},
							},
						}, nil
					},
					GetPrometheusMetricFunc: func(ctx context.Context, input *modzysdk.GetPrometheusMetricInput) (*modzysdk.GetPrometheusMetricOutput, error) {
						return nil, fmt.Errorf("no prom")
					},
				}
			},
			ResourcesFunc: func() modzysdk.ResourcesClient {
				return &modzysdk.ResourcesClientFake{
					GetProcessingModelsFunc: func(ctx context.Context) (*modzysdk.GetProcessingModelsOutput, error) {
						return nil, fmt.Errorf("no processing models")
					},
				}
			},
		}
	}

	stdout, _ := runTestCommand(
		[]string{"status"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			// not used
		},
	)

	if !strings.Contains(stdout, "Data processed") {
		t.Fatalf("out not expected: '%s'", stdout)
	}
}

func TestStatusOutputerError(t *testing.T) {
	outputer := &statusOutputer{}
	err := outputer.Standard(&failWriter{}, &status.TopModel{
		Errors: []error{
			fmt.Errorf("gonna fail"),
		},
	})
	if err == nil {
		t.Fatalf("expected an error")
	}
	if err.Error() != "no" {
		t.Errorf("error was not as expected: %v", err)
	}
}
