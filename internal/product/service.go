package product

import (
	"context"

	"github.com/flukis/pagination-comparation/db/storer"
	"github.com/flukis/pagination-comparation/domain"
)

type Service interface {
	GetPaginationWithPage(ctx context.Context, page, limit int, sortColumn, sortOrder string) domain.FetchProductsResponse[int]
	GetPaginationWithCursor(ctx context.Context, cursor string, limit int) domain.FetchProductsResponse[string]
}

type service struct {
	product storer.Product
}

// GetPaginationWithCursor implements Service.
func (s *service) GetPaginationWithCursor(ctx context.Context, cursor string, limit int) domain.FetchProductsResponse[string] {
	var res domain.FetchProductsResponse[string]

	data, next, total, err := s.product.GetPaginationWithCursor(
		ctx,
		cursor,
		limit,
	)
	if err != nil {
		res.Error = err
		return res
	}

	res.Data = data
	res.Meta.Limit = limit
	res.Meta.Next = next
	res.Meta.Total = total

	return res
}

// GetPaginationWithPage implements Service.
func (s *service) GetPaginationWithPage(ctx context.Context, page, limit int, sortColumn, sortOrder string) domain.FetchProductsResponse[int] {
	var res domain.FetchProductsResponse[int]

	data, next, total, err := s.product.GetPaginationWithPage(
		ctx,
		page,
		limit,
		sortColumn,
		sortOrder,
	)
	if err != nil {
		res.Error = err
		return res
	}

	res.Data = data
	res.Meta.Limit = limit
	res.Meta.Next = next
	res.Meta.Total = total

	return res
}

func NewService(product storer.Product) Service {
	return &service{product}
}
