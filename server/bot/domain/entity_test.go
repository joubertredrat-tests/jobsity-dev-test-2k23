package domain_test

import (
	"joubertredrat-tests/jobsity-dev-test-2k23/bot/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStock(t *testing.T) {
	codeExpected := "APPL.US"
	valueExpedted := "$93.42"
	stockGot := domain.NewStock(codeExpected, valueExpedted)

	assert.Equal(t, codeExpected, stockGot.Code)
	assert.Equal(t, valueExpedted, stockGot.Value)
}
