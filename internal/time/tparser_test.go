package time

import (
	"fmt"
	"testing"
	"time"

	"github.com/modzy/sdk-go/model"
)

type expected struct {
	inT  string
	outT string
	err  error
}

func TestParseT(t *testing.T) {
	defer (func() {
		nowFunc = time.Now
	})()
	nowT := time.Now()
	nowFunc = func() time.Time {
		return nowT
	}

	zero := nowT.Format(model.TimeFormat)
	plus3 := nowT.AddDate(0, 0, 3).Format(model.TimeFormat)
	minus4 := nowT.AddDate(0, 0, -4).Format(model.TimeFormat)

	expectations := []expected{
		{"", zero, nil},
		{"T", zero, nil},
		{"T+3", plus3, nil},
		{"T-4", minus4, nil},
		{"T+junk", "", fmt.Errorf("Day value is not valid: T+JUNK")},
		{"T-junk", "", fmt.Errorf("Day value is not valid: T-JUNK")},
	}

	for _, exp := range expectations {
		actualT, actualErr := ParseT(exp.inT)
		if actualT != exp.outT {
			t.Errorf("result time was '%s', wanted '%s'", actualT, exp.outT)
		}
		if actualErr == nil || exp.err == nil {
			if actualErr != exp.err {
				t.Errorf("result err was '%v', wanted '%v'", actualErr, exp.err)
			}
		} else {
			if actualErr.Error() != exp.err.Error() {
				t.Errorf("result err was '%v', wanted '%v'", actualErr, exp.err)
			}
		}
	}
}
