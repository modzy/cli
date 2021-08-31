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

func TestProjectStatusNoProject(t *testing.T) {
	defer (func() { getClient = defaultGetClient })()

	// instead of a bunch of http handles, just use a fake client
	getClient = func() modzysdk.Client {
		return &modzysdk.ClientFake{
			AccountingFunc: func() modzysdk.AccountingClient {
				return &modzysdk.AccountingClientFake{
					GetProjectDetailsFunc: func(ctx context.Context, input *modzysdk.GetProjectDetailsInput) (*modzysdk.GetProjectDetailsOutput, error) {
						return nil, fmt.Errorf("no project details")
					},
				}
			},
		}
	}

	stdout, _ := runTestCommand(
		[]string{"projects", "status", "projectID"},
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

func TestProjectStatusNoFine(t *testing.T) {
	defer (func() { getClient = defaultGetClient })()

	// instead of a bunch of http handles, just use a fake client
	getClient = func() modzysdk.Client {
		return &modzysdk.ClientFake{
			AccountingFunc: func() modzysdk.AccountingClient {
				return &modzysdk.AccountingClientFake{
					GetProjectDetailsFunc: func(ctx context.Context, input *modzysdk.GetProjectDetailsInput) (*modzysdk.GetProjectDetailsOutput, error) {
						return &modzysdk.GetProjectDetailsOutput{
							Project: modzysdkmodel.AccountingProject{
								AccessKeys: []modzysdkmodel.AccessKey{
									{
										Prefix: "testprefix",
									},
								},
							},
						}, nil
					},
				}
			},
			DashboardFunc: func() modzysdk.DashboardClient {
				return &modzysdk.DashboardClientFake{
					GetDataProcessedFunc: func(ctx context.Context, input *modzysdk.GetDataProcessedInput) (*modzysdk.GetDataProcessedOutput, error) {
						return nil, fmt.Errorf("no data processed")
					},
					GetPredictionsMadeFunc: func(ctx context.Context, input *modzysdk.GetPredictionsMadeInput) (*modzysdk.GetPredictionsMadeOutput, error) {
						return nil, fmt.Errorf("no predictions")
					},
					GetActiveModelsFunc: func(ctx context.Context, input *modzysdk.GetActiveModelsInput) (*modzysdk.GetActiveModelsOutput, error) {
						return nil, fmt.Errorf("no active models")
					},
				}
			},
		}
	}

	stdout, _ := runTestCommand(
		[]string{"projects", "status", "projectID"},
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

func TestProjectStatusFineLists(t *testing.T) {
	defer (func() { getClient = defaultGetClient })()

	// instead of a bunch of http handles, just use a fake client
	getClient = func() modzysdk.Client {
		return &modzysdk.ClientFake{
			AccountingFunc: func() modzysdk.AccountingClient {
				return &modzysdk.AccountingClientFake{
					GetProjectDetailsFunc: func(ctx context.Context, input *modzysdk.GetProjectDetailsInput) (*modzysdk.GetProjectDetailsOutput, error) {
						return &modzysdk.GetProjectDetailsOutput{
							Project: modzysdkmodel.AccountingProject{
								AccessKeys: []modzysdkmodel.AccessKey{
									{
										Prefix: "testprefix",
									},
								},
							},
						}, nil
					},
				}
			},
			DashboardFunc: func() modzysdk.DashboardClient {
				return &modzysdk.DashboardClientFake{
					GetDataProcessedFunc: func(ctx context.Context, input *modzysdk.GetDataProcessedInput) (*modzysdk.GetDataProcessedOutput, error) {
						return &modzysdk.GetDataProcessedOutput{}, fmt.Errorf("no data processed")
					},
					GetPredictionsMadeFunc: func(ctx context.Context, input *modzysdk.GetPredictionsMadeInput) (*modzysdk.GetPredictionsMadeOutput, error) {
						return &modzysdk.GetPredictionsMadeOutput{}, nil
					},
					GetActiveModelsFunc: func(ctx context.Context, input *modzysdk.GetActiveModelsInput) (*modzysdk.GetActiveModelsOutput, error) {
						return &modzysdk.GetActiveModelsOutput{
							Models: []modzysdkmodel.ActiveModelSummary{
								{},
							},
						}, nil
					},
				}
			},
		}
	}

	stdout, _ := runTestCommand(
		[]string{"projects", "status", "projectID"},
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

func TestProjectStatusOutputerError(t *testing.T) {
	outputer := &projectStatusOutputer{}
	err := outputer.Standard(&failWriter{}, &status.ProjectModel{
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
