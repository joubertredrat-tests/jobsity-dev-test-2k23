package infra_test

import (
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/infra"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongoDSN(t *testing.T) {
	dsnExpected := "mongodb://root:password@127.0.0.1:28017"
	dsnGot := infra.MongoDSN("127.0.0.1", "28017", "root", "password")

	assert.Equal(t, dsnExpected, dsnGot)
}
