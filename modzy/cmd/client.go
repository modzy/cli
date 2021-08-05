package cmd

import (
	modzysdk "github.com/modzy/sdk-go"
)

func getClient() modzysdk.Client {
	client := modzysdk.NewClient(rootArgs.BaseURL)
	if rootArgs.APIKey != "" {
		client = client.WithAPIKey(rootArgs.APIKey)
	} else {
		client = client.WithTeamKey(rootArgs.TeamID, rootArgs.TeamToken)
	}
	return client
}
