package cmd

import (
	"fmt"
	"strings"

	ttime "github.com/modzy/cli/internal/time"
	modzysdk "github.com/modzy/sdk-go"
)

type pagingArgs struct {
	Filter []string
	Take   int
	Page   int
	Sort   string
}

func (args pagingArgs) GetPagingInput() (modzysdk.PagingInput, error) {
	paging := modzysdk.PagingInput{
		Page:    args.Page,
		PerPage: args.Take,
	}
	for _, filter := range args.Filter {
		filterSplit := strings.Split(filter, "=")
		if len(filterSplit) != 2 {
			return paging, fmt.Errorf("Filter is not correctly formatted (field=value): '%s'", filter)
		}
		filterValue := filterSplit[1]
		tParsed, _ := ttime.ParseT(filterValue)
		if tParsed != "" {
			filterValue = tParsed
		}
		paging = paging.WithFilterAnd(filterSplit[0], filterValue)
	}

	sortSplits := strings.Split(args.Sort, ":")
	sortFields := strings.Split(sortSplits[0], ",")
	if len(sortFields) > 0 {
		sortDirection := modzysdk.SortDirectionAscending
		if len(sortSplits) > 1 && strings.ToLower(strings.TrimSpace(sortSplits[1])) == "desc" {
			sortDirection = modzysdk.SortDirectionDescending
		}
		paging = paging.WithSort(sortDirection, sortFields...)
	}

	return paging, nil
}
