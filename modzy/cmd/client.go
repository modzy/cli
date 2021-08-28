package cmd

import (
	modzysdk "github.com/modzy/sdk-go"
)

func defaultGetClient() modzysdk.Client {
	client := modzysdk.NewClient(rootArgs.BaseURL)
	if rootArgs.APIKey != "" {
		client = client.WithAPIKey(rootArgs.APIKey)
	} else {
		client = client.WithTeamKey(rootArgs.TeamID, rootArgs.TeamToken)
	}

	if rootArgs.VerboseHTTP {
		client = client.WithOptions(modzysdk.WithHTTPDebugging(true, true))
	}
	return client
}

// var for testing
var getClient = defaultGetClient
