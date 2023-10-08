package pkg_test

import (
	"joubertredrat-tests/jobsity-dev-test-2k23/pkg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatetimeCanonical(t *testing.T) {
	tests := []struct {
		name             string
		time             *time.Time
		datetimeExpected *string
	}{
		{
			name: "Test datetime canonical with valid datetime",
			time: func() *time.Time {
				date, _ := time.Parse("2006-01-02 15:04:05", "2029-06-13 17:27:51")
				return &date
			}(),
			datetimeExpected: func() *string {
				str := "2029-06-13 14:27:51"
				return &str
			}(),
		},
		{
			name:             "Test datetime canonical with no datetime",
			time:             nil,
			datetimeExpected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			datetimeGot := pkg.DatetimeCanonical(test.time)
			assert.Equal(t, test.datetimeExpected, datetimeGot)
		})
	}
}

func TestTimeFromCanonical(t *testing.T) {
	tests := []struct {
		name         string
		datetime     *string
		timeExpected *time.Time
		errExpected  error
	}{
		{
			name: "Test time from canonical with valid datetime",
			datetime: func() *string {
				str := "2023-10-07 14:27:51"
				return &str
			}(),
			timeExpected: func() *time.Time {
				date, _ := time.Parse(pkg.DATETIME_FORMAT, "2023-10-07 14:27:51")
				return &date
			}(),
			errExpected: nil,
		},
		{
			name:         "Test time from canonical with no datetime",
			datetime:     nil,
			timeExpected: nil,
			errExpected:  nil,
		},
		{
			name: "Test time from canonical with invalid datetime",
			datetime: func() *string {
				str := "2023-40-07 14:27:51"
				return &str
			}(),
			timeExpected: nil,
			errExpected: &time.ParseError{
				Layout:     "2006-01-02 15:04:05",
				Value:      "2023-40-07 14:27:51",
				LayoutElem: "01",
				ValueElem:  "-07 14:27:51",
				Message:    ": month out of range",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			timeGot, errGot := pkg.TimeFromCanonical(test.datetime)
			assert.Equal(t, test.timeExpected, timeGot)
			assert.Equal(t, test.errExpected, errGot)
		})
	}
}
