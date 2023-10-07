package domain_test

import (
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	tests := []struct {
		name               string
		page               uint
		itemsPerPage       uint
		paginationExpected domain.Pagination
		errExpected        error
	}{
		{
			name:         "Test create pagination with success",
			page:         1,
			itemsPerPage: 50,
			paginationExpected: domain.Pagination{
				Page:         1,
				ItemsPerPage: 50,
			},
			errExpected: nil,
		},
		{
			name:               "Test create pagination with invalid page",
			page:               0,
			itemsPerPage:       50,
			paginationExpected: domain.Pagination{},
			errExpected:        domain.NewErrPaginationPage(0),
		},
		{
			name:               "Test create pagination with invalid items per page",
			page:               1,
			itemsPerPage:       150,
			paginationExpected: domain.Pagination{},
			errExpected: domain.NewErrPaginationItemsPerPage(
				domain.ITEMS_PER_PAGE_MIN,
				domain.ITEMS_PER_PAGE_MAX,
				150,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paginationGot, errGot := domain.NewPagination(test.page, test.itemsPerPage)

			assert.Equal(t, test.paginationExpected, paginationGot)
			assert.Equal(t, test.errExpected, errGot)
		})
	}
}
