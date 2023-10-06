package infra_test

import (
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/infra"
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
			name: "test datetime canonical with valid datetime",
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
			name:             "test datetime canonical with no datetime",
			time:             nil,
			datetimeExpected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			datetimeGot := infra.DatetimeCanonical(test.time)
			assert.Equal(t, test.datetimeExpected, datetimeGot)
		})
	}
}
