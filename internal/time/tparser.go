package time

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/modzy/sdk-go/model"
)

// exposed for testing
var nowFunc = time.Now

// Parse will take a T string (T+30) and return a valid time string.
// T+20 will return today + 20 days.  T will return today's date.
func ParseT(tString string) (string, error) {
	t, err := parseT(tString)
	if err != nil {
		return "", err
	}
	return t.Format(model.DateFormat), nil

}

func parseT(tString string) (time.Time, error) {
	zero := time.Time{}

	tString = strings.ToUpper(tString)
	if tString == "T" {
		return nowFunc(), nil
	}

	toAdd := 0
	if strings.HasPrefix(tString, "T+") {
		days, err := strconv.Atoi(strings.TrimPrefix(tString, "T+"))
		if err != nil {
			return zero, fmt.Errorf("Day value is not valid: %s", tString)
		}
		toAdd = days
	}

	if strings.HasPrefix(tString, "T-") {
		days, err := strconv.Atoi(strings.TrimPrefix(tString, "T-"))
		if err != nil {
			return zero, fmt.Errorf("Day value is not valid: %s", tString)
		}
		toAdd = -days
	}

	return nowFunc().AddDate(0, 0, toAdd), nil
}
