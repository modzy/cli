package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	modzysdk "github.com/modzy/sdk-go"
)

func TestGetClientAPIKey(t *testing.T) {

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "ApiKey kee" {
			t.Errorf("Authorization header not correct: %s", r.Header.Get("Authorization"))
		}
		if r.Method != "GET" {
			t.Errorf("expected method to be GET, got %s", r.Method)
		}
		if r.RequestURI != "/api/jobs/inputID" {
			t.Errorf("get url not expected: %s", r.RequestURI)
		}
		w.Write([]byte(`{"jobIdentifier": "jsonID"}`))
	}))
	defer serv.Close()

	rootArgs.BaseURL = serv.URL
	rootArgs.APIKey = "kee"
	rootArgs.TeamID = ""
	rootArgs.TeamToken = ""

	client := getClient()

	out, err := client.Jobs().GetJobDetails(context.TODO(), &modzysdk.GetJobDetailsInput{JobIdentifier: "inputID"})
	if err != nil {
		t.Errorf("err not nil: %v", err)
	}
	if out.Details.JobIdentifier != "jsonID" {
		t.Errorf("response not parsed")
	}
}

func TestGetClientTeamKey(t *testing.T) {

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("Authorization header not correct: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Modzy-Team-Id") != "team" {
			t.Errorf("Team header not correct: %s", r.Header.Get("Modzy-Team-ID"))
		}
		if r.Method != "GET" {
			t.Errorf("expected method to be GET, got %s", r.Method)
		}
		if r.RequestURI != "/api/jobs/inputID" {
			t.Errorf("get url not expected: %s", r.RequestURI)
		}
		w.Write([]byte(`{"jobIdentifier": "jsonID"}`))
	}))
	defer serv.Close()

	rootArgs.BaseURL = serv.URL
	rootArgs.APIKey = ""
	rootArgs.TeamID = "team"
	rootArgs.TeamToken = "token"
	// test verbose as well (only coverage)
	rootArgs.VerboseHTTP = true

	client := getClient()

	out, err := client.Jobs().GetJobDetails(context.TODO(), &modzysdk.GetJobDetailsInput{JobIdentifier: "inputID"})
	if err != nil {
		t.Errorf("err not nil: %v", err)
	}
	if out.Details.JobIdentifier != "jsonID" {
		t.Errorf("response not parsed")
	}
}
