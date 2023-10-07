package domain

import "context"

const (
	ITEMS_PER_PAGE_MIN = 10
	ITEMS_PER_PAGE_MAX = 100
)

type Pagination struct {
	Page         uint
	ItemsPerPage uint
}

func NewPagination(p, i uint) (Pagination, error) {
	if p < 1 {
		return Pagination{}, NewErrPaginationPage(p)
	}

	if i < ITEMS_PER_PAGE_MIN || i > ITEMS_PER_PAGE_MAX {
		return Pagination{}, NewErrPaginationItemsPerPage(ITEMS_PER_PAGE_MIN, ITEMS_PER_PAGE_MAX, i)
	}

	return Pagination{
		Page:         p,
		ItemsPerPage: i,
	}, nil
}

type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetAuthenticated(ctx context.Context, user User) (User, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message Message) (Message, error)
	List(ctx context.Context, pagination Pagination) ([]Message, error)
}
