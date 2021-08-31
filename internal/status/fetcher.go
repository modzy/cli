package status

import (
	"context"

	modzysdk "github.com/modzy/sdk-go"
)

type Fetcher interface {
	FetchTop(ctx context.Context) TopModel
	FetchProject(ctx context.Context, projectID string) ProjectModel
}

type StandardFetcher struct {
	client modzysdk.Client
}

var _ Fetcher = &StandardFetcher{}

func NewFetcher(client modzysdk.Client) Fetcher {
	return &StandardFetcher{
		client: client,
	}
}
