package cmd

import (
	"testing"

	modzysdk "github.com/modzy/sdk-go"
)

func TestGetPagingInput(t *testing.T) {
	args := pagingArgs{
		Filter: []string{
			"a=1",
			"b=T+1",
		},
		Take: 7,
		Page: 31,
		Sort: "f,g:dESc",
	}

	pi, err := args.GetPagingInput()

	if err != nil {
		t.Fatalf("error not expected: %v", err)
	}

	if pi.Page != 31 {
		t.Errorf("page wanted 31, got %d", pi.Page)
	}

	if pi.PerPage != 7 {
		t.Errorf("per-page wanted 7, got %d", pi.PerPage)
	}

	if len(pi.SortBy) != 2 {
		t.Errorf("sort not 2 fields: %v", pi.SortBy)
	} else {
		if pi.SortBy[0] != "f" {
			t.Errorf("sort[0] wanted f, got %s", pi.SortBy[0])
		}
	}

	if pi.SortDirection != modzysdk.SortDirectionDescending {
		t.Errorf("sort was not descending")
	}
}

func TestGetPagingInputBadFilter(t *testing.T) {
	args := pagingArgs{
		Filter: []string{"a:1"},
	}

	_, err := args.GetPagingInput()

	if err == nil {
		t.Fatalf("error was expected")
	}
}
